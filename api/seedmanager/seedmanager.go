//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package seedmanager

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
)

// Expiration represents time of expiration
type Expiration = int64

// DefaultTTL is default time period for deleting expired seeds
const DefaultTTL = 5 * time.Second

// DefaultCleanPeriod default time period for launching cleaning goroutine
const DefaultCleanPeriod = 5 * time.Second

type storedSeed struct {
	expiration Expiration
	pulse      insolar.PulseNumber
}

// SeedManager manages working with seed pool
// It's thread safe
type SeedManager struct {
	mutex    sync.RWMutex
	seedPool map[Seed]storedSeed
	ttl      time.Duration
	stopped  chan struct{}
	id       uint32
}

// New creates new seed manager with default params
func New() *SeedManager {
	return NewSpecified(DefaultTTL, DefaultCleanPeriod)
}

var globalID uint32

// NewSpecified creates new seed manager with custom params
func NewSpecified(ttl time.Duration, cleanPeriod time.Duration) *SeedManager {
	atomic.AddUint32(&globalID, 1)
	sm := SeedManager{
		seedPool: make(map[Seed]storedSeed),
		ttl:      ttl,
		stopped:  make(chan struct{}),
		id:       atomic.LoadUint32(&globalID),
	}

	log.Info("Creating NewSpecified: ", ", ID: ", sm.id)

	ticker := time.NewTicker(cleanPeriod)

	go func() {
		var stop = false
		for !stop {
			select {
			case <-ticker.C:
				sm.deleteExpired()
			case <-sm.stopped:
				stop = true
			}
		}
		sm.stopped <- struct{}{}
	}()

	return &sm
}

func (sm *SeedManager) Stop() {
	sm.stopped <- struct{}{}
	<-sm.stopped
}

// Add adds seed to pool
func (sm *SeedManager) Add(seed Seed, pulse insolar.PulseNumber) {
	expTime := time.Now().Add(sm.ttl).UnixNano()

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.seedPool[seed] = storedSeed{
		expiration: expTime,
		pulse:      pulse,
	}

	log.Info("SeedManager.Add. ID: ", sm.id, ". Now: ", time.Now(),
		", ttl: ", sm.ttl.String(),
		", expTime:", expTime,
		"content: ", sm.seedPool,
		", seed: ", seed)

}

func (sm *SeedManager) isExpired(expTime Expiration) bool {
	log.Info("isExpired. expTime: ", expTime, " , now: ", time.Now().UnixNano())
	return expTime < time.Now().UnixNano()
}

// Exists checks whether seed in the pool
func (sm *SeedManager) Pop(seed Seed) (insolar.PulseNumber, bool) {
	sm.mutex.RLock()
	stored, ok := sm.seedPool[seed]
	log.Info("Pop: ok: ", ok, ", ID: ", sm.id, ". content: ", sm.seedPool, ", seed: ", seed)
	sm.mutex.RUnlock()

	if ok && !sm.isExpired(stored.expiration) {
		sm.mutex.Lock()
		defer sm.mutex.Unlock()

		log.Info("Pop seed: delete: ", seed, ", ID: ", sm.id)
		delete(sm.seedPool, seed)
		return stored.pulse, true
	}

	log.Info("Pop seed: NO: ", seed, ", ID: ", sm.id)

	return 0, false
}

func (sm *SeedManager) deleteExpired() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for seed, stored := range sm.seedPool {
		if sm.isExpired(stored.expiration) {
			log.Info("deleteExpired: ", seed, ", ID: ", sm.id)
			delete(sm.seedPool, seed)
		}
	}
}

// SeedFromBytes converts slice of bytes to Seed. Returns nil if slice's size is not equal to SeedSize
func SeedFromBytes(slice []byte) *Seed {
	if len(slice) != int(SeedSize) {
		return nil
	}
	var result Seed
	copy(result[:], slice[:SeedSize])
	return &result
}
