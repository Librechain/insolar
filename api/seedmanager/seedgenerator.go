/*
 *    Copyright 2018 Insolar
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

package seedmanager

import (
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

// Seed is a type of seed
type Seed = [32]byte

// SeedGenerator holds logic with seed generation
type SeedGenerator struct {
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Next returns next random seed
func (sg *SeedGenerator) Next() (*Seed, error) {
	seed := Seed{}
	_, err := rand.Read(seed[:])
	if err != nil {
		return nil, errors.Wrap(err, "[ SeedGenerator::Next ]")
	}

	return &seed, nil
}
