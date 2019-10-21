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

package conveyor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/conveyor/sworker"
	"github.com/insolar/insolar/conveyor/tools"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/pulse"
)

func TestConveyor(t *testing.T) {
	machineConfig := smachine.SlotMachineConfig{
		SlotPageSize:    1000,
		PollingPeriod:   500 * time.Millisecond,
		PollingTruncate: 1 * time.Millisecond,
		ScanCountLimit:  100000,
	}

	factoryFn := func(pulse.Number, InputEvent) smachine.CreateFunc {
		return nil
	}
	conveyor := NewPulseConveyor(context.Background(), machineConfig, factoryFn, machineConfig, nil)

	pd := pulse.NewFirstPulsarData(10, longbits.Bits256{})
	signal := tools.NewVersionedSignal()

	go worker(conveyor, signal)

	require.NoError(t, conveyor.CommitPulseChange(pd))

	for i := 0; i < 100; i++ {
		pd = pd.CreateNextPulsarPulse(2, func() longbits.Bits256 {
			return longbits.Bits256{}
		})
		require.NoError(t, conveyor.CommitPulseChange(pd))
		time.Sleep(time.Second)
	}
	fmt.Println("======================")
	time.Sleep(time.Hour)
}

func worker(conveyor *PulseConveyor, signal tools.VersionedSignal) {
	workerFactory := sworker.NewAttachableSimpleSlotWorker()
	sm := conveyor.slotMachine
	for {
		var (
			repeatNow    bool
			nextPollTime time.Time
		)
		workerFactory.AttachTo(sm, signal.Mark(), 100, func(worker smachine.AttachedSlotWorker) {
			repeatNow, nextPollTime = sm.ScanOnce(0, worker)
		})

		if repeatNow {
			continue
		}
		if !nextPollTime.IsZero() {
			time.Sleep(time.Until(nextPollTime))
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}
