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

package executor

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/stretchr/testify/require"
)

type splitCase struct {
	name             string
	jetID            insolar.JetID
	recordsThreshold int
	overflowCount    int

	// represents expected values on every pulse (elements count should match!)
	recordsPerPulse []int
	dropExceed      []int
	hasSplit        []bool
}

var cases = []splitCase{
	{
		name:             "no_split",
		jetID:            jet.NewIDFromString("10"),
		recordsThreshold: 5,
		overflowCount:    0,
		recordsPerPulse:  []int{5, 3},
		dropExceed:       []int{0, 0},
		hasSplit:         []bool{false, false},
	},
	{
		name:             "split",
		jetID:            jet.NewIDFromString("10"),
		recordsThreshold: 4,
		overflowCount:    0,
		recordsPerPulse:  []int{5, 3},
		dropExceed:       []int{1, 0},
		hasSplit:         []bool{true, false},
	},
	{
		name:             "split_with_overflow",
		jetID:            jet.NewIDFromString("10"),
		recordsThreshold: 4,
		overflowCount:    1,
		recordsPerPulse:  []int{5, 5, 5},
		dropExceed:       []int{1, 2, 0},
		hasSplit:         []bool{false, true, false},
	},
	{
		name:             "no_split_with_overflow",
		jetID:            jet.NewIDFromString("10"),
		recordsThreshold: 4,
		overflowCount:    1,
		recordsPerPulse:  []int{5, 4, 5},
		dropExceed:       []int{1, 0, 1},
		hasSplit:         []bool{false, false, false},
	},
}

func TestJetSplitter(t *testing.T) {
	ctx := inslogger.TestContext(t)

	// prepare initial pulses
	pn := gen.PulseNumber()
	// just avoid special pulses
	if pn < 60000 {
		pn += 60000
	}
	previousPulse, currentPulse, newPulse := pn, pn+1, pn+2

	initialJets := []insolar.JetID{
		jet.NewIDFromString("0"),
		jet.NewIDFromString("10"),
		jet.NewIDFromString("11"),
	}

	checkCase := func(t *testing.T, sc splitCase) {
		// real compomnents
		jetStore := jet.NewStore()
		err := jetStore.Update(ctx, currentPulse, true, initialJets...)
		require.NoError(t, err, "jet store updated with initial jets")
		db := drop.NewStorageMemory()
		dropAccessor := db
		dropModifier := db
		cfg := configuration.JetSplit{
			ThresholdRecordsCount:  sc.recordsThreshold,
			ThresholdOverflowCount: sc.overflowCount,
		}
		// mocks
		jetCalc := NewJetCalculatorMock(t)
		collectionAccessor := object.NewRecordCollectionAccessorMock(t)
		pulseCalc := pulse.NewCalculatorMock(t)

		// create splitter
		splitter := NewJetSplitter(
			cfg, jetCalc, jetStore, jetStore,
			dropAccessor, dropModifier,
			pulseCalc, collectionAccessor,
		)

		// no filter for ID
		jetCalc.MineForPulseFunc = func(ctx context.Context, pn insolar.PulseNumber) []insolar.JetID {
			return jetStore.All(ctx, pn)
		}

		for i, recordsCount := range sc.recordsPerPulse {
			delta := insolar.PulseNumber(i)
			current, newpulse := currentPulse+delta, newPulse+delta
			pulseCalc.BackwardsMock.Return(insolar.Pulse{PulseNumber: previousPulse + delta}, nil)

			pulseStartedWithJets := jetStore.All(ctx, current)

			collectionAccessor.ForPulseFunc = func(_ context.Context, jetID insolar.JetID, pn insolar.PulseNumber) []record.Material {
				if pn == current && jetID == sc.jetID {
					return make([]record.Material, recordsCount)
				}
				return nil
			}

			gotJets, err := splitter.Do(ctx, current, newpulse)
			require.NoError(t, err, "splitter.Do performed")

			dropThreshold := splitter.getDropThreshold(ctx, sc.jetID, current)
			require.NoErrorf(t, err, "get jet's %v drop for pulse %v", sc.jetID.DebugString(), current)

			require.Equalf(t, sc.dropExceed[i], dropThreshold,
				"check drop.SplitThresholdExceeded for pulse with offset: %v", i)

			var expectJets []insolar.JetID
			for _, jetID := range pulseStartedWithJets {
				if sc.hasSplit[i] && (jetID == sc.jetID) {
					left, right := jet.Siblings(sc.jetID)
					expectJets = append(expectJets, left, right)
					continue
				}
				expectJets = append(expectJets, jetID)
			}
			expectMsg := fmt.Sprintf("jet %v should%v split on +%v pulse",
				sc.jetID.DebugString(),
				map[bool]string{false: " not"}[sc.hasSplit[i]],
				i,
			)
			require.Equal(t, jsort(expectJets), jsort(gotJets), expectMsg)
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { checkCase(t, c) })
	}
}

func jsort(jets []insolar.JetID) []string {
	sort.Slice(jets, func(i, j int) bool {
		return bytes.Compare(jets[i][:], jets[j][:]) == -1
	})
	result := make([]string, 0, len(jets))
	for _, j := range jets {
		result = append(result, j.DebugString())
	}
	return result
}
