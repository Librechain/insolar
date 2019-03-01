/*
 *    Copyright 2019 Insolar
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

package db

import (
	"sync"

	"github.com/insolar/insolar/core"
)

type JetIndex struct {
	lock    sync.Mutex
	storage map[core.JetID]recordSet
}

func NewJetIndex() *JetIndex {
	return &JetIndex{storage: map[core.JetID]recordSet{}}
}

type recordSet map[core.RecordID]struct{}

func (i *JetIndex) Add(id core.RecordID, jetID core.JetID) {
	i.lock.Lock()
	defer i.lock.Unlock()

	jet, ok := i.storage[jetID]
	if !ok {
		jet = recordSet{}
		i.storage[jetID] = jet
	}
	jet[id] = struct{}{}
}

func (i *JetIndex) Delete(id core.RecordID, jetID core.JetID) {
	i.lock.Lock()
	defer i.lock.Unlock()

	jet, ok := i.storage[jetID]
	if !ok {
		return
	}

	delete(jet, id)
	if len(jet) == 0 {
		delete(i.storage, jetID)
	}
}
