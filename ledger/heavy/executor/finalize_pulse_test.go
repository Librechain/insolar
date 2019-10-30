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

package executor_test

import (
	"context"
	"sync"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/testutils/network"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type TestBadgerGCRunner struct {
	lock   sync.RWMutex
	count  uint
	called chan struct{}
}

func (t *TestBadgerGCRunner) RunValueGC(ctx context.Context) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.count++
	t.called <- struct{}{}
}

func (t *TestBadgerGCRunner) getCount() uint {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.count
}

func TestBadgerGCRunInfo(t *testing.T) {

	ctx := inslogger.TestContext(t)

	t.Run("call every time if frequency equal 1", func(t *testing.T) {
		t.Parallel()
		runner := &TestBadgerGCRunner{
			called: make(chan struct{}),
		}
		info := executor.NewBadgerGCRunInfo(runner, 1)
		for i := 1; i < 5; i++ {
			info.RunGCIfNeeded(ctx)
			<-runner.called
			require.Equal(t, uint(i), runner.getCount())
		}
	})

	t.Run("no call if frequency equal 0", func(t *testing.T) {
		t.Parallel()
		runner := &TestBadgerGCRunner{}
		info := executor.NewBadgerGCRunInfo(runner, 0)
		for i := 1; i < 5; i++ {
			info.RunGCIfNeeded(ctx)
			require.Equal(t, uint(0), runner.getCount())
		}
	})
	t.Run("call every second time if frequency equal 2", func(t *testing.T) {
		t.Parallel()
		runner := &TestBadgerGCRunner{
			called: make(chan struct{}, 1),
		}
		frequency := uint(2)
		info := executor.NewBadgerGCRunInfo(runner, frequency)
		info.RunGCIfNeeded(ctx)
		require.Equal(t, uint(0), runner.getCount())
		info.RunGCIfNeeded(ctx)
		<-runner.called
		require.Equal(t, uint(1), runner.getCount())
		info.RunGCIfNeeded(ctx)
		require.Equal(t, uint(1), runner.getCount())
		info.RunGCIfNeeded(ctx)
		<-runner.called
		require.Equal(t, uint(2), runner.getCount())
	})
}

func TestFinalizePulse_HappyPath(t *testing.T) {
	ctx := inslogger.TestContext(t)

	testPulse := insolar.PulseNumber(pulse.MinTimePulse)
	targetPulse := testPulse + 1

	pc := network.NewPulseCalculatorMock(t)
	pc.ForwardsMock.Return(insolar.Pulse{PulseNumber: targetPulse}, nil)

	bkp := executor.NewBackupMakerMock(t)
	bkp.MakeBackupMock.Return(nil)

	jk := executor.NewJetKeeperMock(t)
	var hasConfirmCount uint32
	hasConfirm := func(ctx context.Context, pulse insolar.PulseNumber) bool {
		var p bool
		switch hasConfirmCount {
		case 0:
			p = true
		case 1:
			p = false
		}
		hasConfirmCount++
		return p
	}

	jk.HasAllJetConfirmsMock.Set(hasConfirm)

	js := jet.NewStorageMock(t)
	js.AllMock.Return(nil)
	jk.StorageMock.Return(js)

	var topSyncCount uint32
	topSync := func() insolar.PulseNumber {
		var p insolar.PulseNumber
		switch topSyncCount {
		case 0:
			p = testPulse
		case 1:
			p = targetPulse
		}
		topSyncCount++
		return p
	}

	jk.TopSyncPulseMock.Set(topSync)
	jk.AddBackupConfirmationMock.Return(nil)

	indexes := object.NewIndexModifierMock(t)
	indexes.UpdateLastKnownPulseMock.Return(nil)

	executor.FinalizePulse(ctx, pc, bkp, jk, indexes, targetPulse, testBadgerGCInfo())
}

func testBadgerGCInfo() *executor.BadgerGCRunInfo {
	return executor.NewBadgerGCRunInfo(&TestBadgerGCRunner{}, 1)
}

func TestFinalizePulse_JetIsNotConfirmed(t *testing.T) {
	ctx := inslogger.TestContext(t)

	testPulse := insolar.PulseNumber(pulse.MinTimePulse)

	jk := executor.NewJetKeeperMock(t)
	jk.HasAllJetConfirmsMock.Return(false)

	executor.FinalizePulse(ctx, nil, nil, jk, nil, testPulse, testBadgerGCInfo())
}

func TestFinalizePulse_CantGteNextPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)

	testPulse := insolar.PulseNumber(pulse.MinTimePulse)

	jk := executor.NewJetKeeperMock(t)
	jk.HasAllJetConfirmsMock.Return(true)
	jk.TopSyncPulseMock.Return(testPulse)

	pc := network.NewPulseCalculatorMock(t)
	pc.ForwardsMock.Return(insolar.Pulse{}, errors.New("Test"))

	executor.FinalizePulse(ctx, pc, nil, jk, nil, testPulse, testBadgerGCInfo())
}

func TestFinalizePulse_BackupError(t *testing.T) {
	ctx := inslogger.TestContext(t)

	testPulse := insolar.PulseNumber(pulse.MinTimePulse)
	targetPulse := testPulse + 1

	jk := executor.NewJetKeeperMock(t)
	jk.HasAllJetConfirmsMock.Return(true)
	jk.TopSyncPulseMock.Return(targetPulse)

	js := jet.NewStorageMock(t)
	js.AllMock.Return(nil)
	jk.StorageMock.Return(js)

	pc := network.NewPulseCalculatorMock(t)
	pc.ForwardsMock.Return(insolar.Pulse{PulseNumber: targetPulse}, nil)

	bkp := executor.NewBackupMakerMock(t)
	bkp.MakeBackupMock.Return(executor.ErrAlreadyDone)

	executor.FinalizePulse(ctx, pc, bkp, jk, nil, targetPulse, testBadgerGCInfo())
}

func TestFinalizePulse_NotNextPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)

	testPulse := insolar.PulseNumber(pulse.MinTimePulse)

	jk := executor.NewJetKeeperMock(t)
	jk.HasAllJetConfirmsMock.Return(true)
	jk.TopSyncPulseMock.Return(testPulse)

	pc := network.NewPulseCalculatorMock(t)
	pc.ForwardsMock.Return(insolar.Pulse{PulseNumber: testPulse}, nil)

	executor.FinalizePulse(ctx, pc, nil, jk, nil, testPulse+10, testBadgerGCInfo())
}
