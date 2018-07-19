/*
 *    Copyright 2018 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package record

import (
	"crypto/sha256"
)

const (
	HashSize = sha256.Size

	RecTypeCodeBase = RecordType(iota + 1)
	RecTypeCodeAmendment
	RecTypeObjectInstance
	RecTypeObjectAmendment
	RecTypeObjectData
)

type RecordHash [HashSize]byte

type RecordType uint

type ProjectionType uint

type Memory []byte

type Record interface {
	Hash() RecordHash
	TimeSlot() uint64
	Type() RecordType
}

// TODO: Should implement normal RecordReference type (not interface)
// TODO: Globally unique record identifier must be found
type RecordReference interface {
	Record
}

type AppDataRecord struct {
	timeSlotNo     uint64
	recType        RecordType
	collisionNonce uint // TODO: combine nonce with type?
}

func (r *AppDataRecord) Hash() RecordHash {
	panic("implement me")
}

func (r *AppDataRecord) TimeSlot() uint64 {
	return r.timeSlotNo
}

func (r *AppDataRecord) Type() RecordType {
	return r.recType
}
