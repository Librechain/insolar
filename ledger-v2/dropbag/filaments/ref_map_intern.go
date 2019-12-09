//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package filaments

import (
	"math"
	"math/bits"
	"runtime"

	"github.com/insolar/insolar/ledger-v2/unsafekit"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/reference"
)

const MinInterningPage = 16
const MaxInterningPage = 65536

func NewInterningRefMap(pageSize int, hashSeed uint32) *RefSemiMap {
	switch {
	case pageSize < MinInterningPage:
		panic("illegal value")
	case pageSize > MaxInterningPage:
		panic("illegal value")
	}

	pageBits := uint8(bits.Len(uint(pageSize)))

	buckets := [][]refMapBucket{make([]refMapBucket, 1, 1<<pageBits)}
	buckets[0][0] = refMapBucket{localRef: emptyLocalRef}

	return &RefSemiMap{
		map0:     map[longbits.ByteString]uint32{emptyLocalRefKey: 0},
		buckets:  buckets,
		hashSeed: hashSeed,
		pageBits: pageBits,
	}
}

var emptyLocalRef = reference.EmptyLocal()
var emptyLocalRefKey = unsafekit.WrapLocalRef(emptyLocalRef)

type RefSemiMap struct {
	hashSeed uint32

	map0    map[longbits.ByteString]uint32
	buckets [][]refMapBucket

	pageBits uint8
}

type BucketState uint32
type ValueSelector struct {
	BucketId, ValueId uint32
	State             BucketState
}

type RefMapHash = uint32

type refMapBucket struct {
	localRef *reference.Local
	refHash  RefMapHash
	state    BucketState
}

func (b refMapBucket) IsEmpty() bool {
	return b.state == 0
}

func (m *RefSemiMap) InternHolder(ref reference.Holder) reference.Holder {
	switch {
	case ref == nil:
		return nil
	case ref.IsEmpty():
		return reference.Empty()
	}

	p0, p1 := ref.GetLocal(), ref.GetBase()
	p0i, p1i := m.Intern(p0), p1
	if p0 == p1 {
		p1i = p0i
	} else {
		p1i = m.Intern(p1)
	}
	if p0 != p0i || p1 != p1i {
		return reference.NewNoCopy(p0, p1)
	}
	return ref
}

const bucketIndexMask = math.MaxInt32

func (m *RefSemiMap) InternedKeyCount() int {
	return len(m.map0)
}

func (m *RefSemiMap) Intern(ref *reference.Local) *reference.Local {
	if ref == nil {
		return nil
	}
	_, r, _ := m.intern(ref)
	return r
}

func (m *RefSemiMap) getBucket(bucketIndex uint32) *refMapBucket {
	return &m.buckets[bucketIndex>>m.pageBits][bucketIndex&(1<<m.pageBits-1)]
}

func (m *RefSemiMap) GetInterned(bucketIndex uint32) *reference.Local {
	bucket := m.getBucket(bucketIndex)
	return bucket.localRef
}

func (m *RefSemiMap) intern(ref *reference.Local) (uint32, *reference.Local, longbits.ByteString) {
	switch {
	case ref == nil:
		panic("illegal value")
	case ref.IsEmpty():
		return 0, reference.EmptyLocal(), ""
	}

	key := unsafekit.WrapLocalRef(ref) // MUST keep ref as long as key is in use
	if bucketIndex, ok := m.getIndexWithKey(ref, key); ok {
		return bucketIndex, m.getBucket(bucketIndex).localRef, key
	}
	bucketIndex := uint32(len(m.buckets))
	if bucketIndex > bucketIndexMask {
		panic("illegal state - overflow")
	}

	pageN := len(m.buckets)
	page := &m.buckets[pageN-1]
	if len(*page) == cap(*page) {
		m.buckets = append(m.buckets, make([]refMapBucket, 0, 1<<m.pageBits))
		page = &m.buckets[pageN]
	}

	*page = append(*page, refMapBucket{localRef: ref})

	m.map0[key] = bucketIndex
	return bucketIndex, ref, key // (ref) stays, no need for KeepAlive()
}

func (m *RefSemiMap) getIndexWithKey(ref *reference.Local, refKey longbits.ByteString) (uint32, bool) {
	bucketIndex, ok := m.map0[refKey]
	runtime.KeepAlive(ref) // make sure that (ref) stays while (key) is in use by mapaccess()
	return bucketIndex & bucketIndexMask, ok
}

func (m *RefSemiMap) getIndex(ref *reference.Local) (uint32, bool) {
	return m.getIndexWithKey(ref, unsafekit.WrapLocalRef(ref))
}

func (m *RefSemiMap) Find(key reference.Holder) (ValueSelector, bool) {
	switch {
	case key.IsEmpty():
		return ValueSelector{}, false
	}

	var bucket *refMapBucket
	bucketIndex, ok := m.getIndex(key.GetLocal())
	if !ok {
		return ValueSelector{}, false
	}
	bucket = m.getBucket(bucketIndex)
	if bucket.IsEmpty() {
		return ValueSelector{}, false
	}

	if baseIndex, ok := m.getIndex(key.GetBase()); !ok {
		return ValueSelector{}, false
	} else {
		return ValueSelector{bucketIndex, baseIndex, bucket.state}, true
	}
}

func (m *RefSemiMap) TryPut(key reference.Holder,
	valueFn func(internedKey reference.Holder, selector ValueSelector) BucketState,
) bool {
	switch {
	case valueFn == nil:
		panic("illegal value")
	case key.IsEmpty():
		return false
	}

	p0, p1 := key.GetLocal(), key.GetBase()
	bucket0, p0i, p0k := m.intern(p0)
	bucket1, p1i, _ := m.intern(p1)

	if p0 != p0i || p1 != p1i {
		key = reference.NewNoCopy(p0i, p1i)
	}

	bucket := m.getBucket(bucket0)
	prevState := bucket.state
	switch newState := valueFn(key, ValueSelector{bucket0, bucket1, prevState}); {
	case prevState == newState:
		return false
	case newState == 0:
		panic("illegal value")
	case prevState == 0:
		bucket.refHash = Hash32(p0k, m.hashSeed)
		runtime.KeepAlive(p0) // ensures that (p0k) is ok
		fallthrough
	default:
		bucket.state = newState
		return true
	}
}
