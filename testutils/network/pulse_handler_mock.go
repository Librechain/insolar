package network

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	mm_network "github.com/insolar/insolar/network"
)

// PulseHandlerMock implements network.PulseHandler
type PulseHandlerMock struct {
	t minimock.Tester

	funcHandlePulse          func(ctx context.Context, pulse insolar.Pulse, originalPacket mm_network.ReceivedPacket)
	inspectFuncHandlePulse   func(ctx context.Context, pulse insolar.Pulse, originalPacket mm_network.ReceivedPacket)
	afterHandlePulseCounter  uint64
	beforeHandlePulseCounter uint64
	HandlePulseMock          mPulseHandlerMockHandlePulse
}

// NewPulseHandlerMock returns a mock for network.PulseHandler
func NewPulseHandlerMock(t minimock.Tester) *PulseHandlerMock {
	m := &PulseHandlerMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.HandlePulseMock = mPulseHandlerMockHandlePulse{mock: m}
	m.HandlePulseMock.callArgs = []*PulseHandlerMockHandlePulseParams{}

	return m
}

type mPulseHandlerMockHandlePulse struct {
	mock               *PulseHandlerMock
	defaultExpectation *PulseHandlerMockHandlePulseExpectation
	expectations       []*PulseHandlerMockHandlePulseExpectation

	callArgs []*PulseHandlerMockHandlePulseParams
	mutex    sync.RWMutex
}

// PulseHandlerMockHandlePulseExpectation specifies expectation struct of the PulseHandler.HandlePulse
type PulseHandlerMockHandlePulseExpectation struct {
	mock   *PulseHandlerMock
	params *PulseHandlerMockHandlePulseParams

	Counter uint64
}

// PulseHandlerMockHandlePulseParams contains parameters of the PulseHandler.HandlePulse
type PulseHandlerMockHandlePulseParams struct {
	ctx            context.Context
	pulse          insolar.Pulse
	originalPacket mm_network.ReceivedPacket
}

// Expect sets up expected params for PulseHandler.HandlePulse
func (mmHandlePulse *mPulseHandlerMockHandlePulse) Expect(ctx context.Context, pulse insolar.Pulse, originalPacket mm_network.ReceivedPacket) *mPulseHandlerMockHandlePulse {
	if mmHandlePulse.mock.funcHandlePulse != nil {
		mmHandlePulse.mock.t.Fatalf("PulseHandlerMock.HandlePulse mock is already set by Set")
	}

	if mmHandlePulse.defaultExpectation == nil {
		mmHandlePulse.defaultExpectation = &PulseHandlerMockHandlePulseExpectation{}
	}

	mmHandlePulse.defaultExpectation.params = &PulseHandlerMockHandlePulseParams{ctx, pulse, originalPacket}
	for _, e := range mmHandlePulse.expectations {
		if minimock.Equal(e.params, mmHandlePulse.defaultExpectation.params) {
			mmHandlePulse.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmHandlePulse.defaultExpectation.params)
		}
	}

	return mmHandlePulse
}

// Inspect accepts an inspector function that has same arguments as the PulseHandler.HandlePulse
func (mmHandlePulse *mPulseHandlerMockHandlePulse) Inspect(f func(ctx context.Context, pulse insolar.Pulse, originalPacket mm_network.ReceivedPacket)) *mPulseHandlerMockHandlePulse {
	if mmHandlePulse.mock.inspectFuncHandlePulse != nil {
		mmHandlePulse.mock.t.Fatalf("Inspect function is already set for PulseHandlerMock.HandlePulse")
	}

	mmHandlePulse.mock.inspectFuncHandlePulse = f

	return mmHandlePulse
}

// Return sets up results that will be returned by PulseHandler.HandlePulse
func (mmHandlePulse *mPulseHandlerMockHandlePulse) Return() *PulseHandlerMock {
	if mmHandlePulse.mock.funcHandlePulse != nil {
		mmHandlePulse.mock.t.Fatalf("PulseHandlerMock.HandlePulse mock is already set by Set")
	}

	if mmHandlePulse.defaultExpectation == nil {
		mmHandlePulse.defaultExpectation = &PulseHandlerMockHandlePulseExpectation{mock: mmHandlePulse.mock}
	}

	return mmHandlePulse.mock
}

//Set uses given function f to mock the PulseHandler.HandlePulse method
func (mmHandlePulse *mPulseHandlerMockHandlePulse) Set(f func(ctx context.Context, pulse insolar.Pulse, originalPacket mm_network.ReceivedPacket)) *PulseHandlerMock {
	if mmHandlePulse.defaultExpectation != nil {
		mmHandlePulse.mock.t.Fatalf("Default expectation is already set for the PulseHandler.HandlePulse method")
	}

	if len(mmHandlePulse.expectations) > 0 {
		mmHandlePulse.mock.t.Fatalf("Some expectations are already set for the PulseHandler.HandlePulse method")
	}

	mmHandlePulse.mock.funcHandlePulse = f
	return mmHandlePulse.mock
}

// HandlePulse implements network.PulseHandler
func (mmHandlePulse *PulseHandlerMock) HandlePulse(ctx context.Context, pulse insolar.Pulse, originalPacket mm_network.ReceivedPacket) {
	mm_atomic.AddUint64(&mmHandlePulse.beforeHandlePulseCounter, 1)
	defer mm_atomic.AddUint64(&mmHandlePulse.afterHandlePulseCounter, 1)

	if mmHandlePulse.inspectFuncHandlePulse != nil {
		mmHandlePulse.inspectFuncHandlePulse(ctx, pulse, originalPacket)
	}

	params := &PulseHandlerMockHandlePulseParams{ctx, pulse, originalPacket}

	// Record call args
	mmHandlePulse.HandlePulseMock.mutex.Lock()
	mmHandlePulse.HandlePulseMock.callArgs = append(mmHandlePulse.HandlePulseMock.callArgs, params)
	mmHandlePulse.HandlePulseMock.mutex.Unlock()

	for _, e := range mmHandlePulse.HandlePulseMock.expectations {
		if minimock.Equal(e.params, params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return
		}
	}

	if mmHandlePulse.HandlePulseMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmHandlePulse.HandlePulseMock.defaultExpectation.Counter, 1)
		want := mmHandlePulse.HandlePulseMock.defaultExpectation.params
		got := PulseHandlerMockHandlePulseParams{ctx, pulse, originalPacket}
		if want != nil && !minimock.Equal(*want, got) {
			mmHandlePulse.t.Errorf("PulseHandlerMock.HandlePulse got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		return

	}
	if mmHandlePulse.funcHandlePulse != nil {
		mmHandlePulse.funcHandlePulse(ctx, pulse, originalPacket)
		return
	}
	mmHandlePulse.t.Fatalf("Unexpected call to PulseHandlerMock.HandlePulse. %v %v %v", ctx, pulse, originalPacket)

}

// HandlePulseAfterCounter returns a count of finished PulseHandlerMock.HandlePulse invocations
func (mmHandlePulse *PulseHandlerMock) HandlePulseAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmHandlePulse.afterHandlePulseCounter)
}

// HandlePulseBeforeCounter returns a count of PulseHandlerMock.HandlePulse invocations
func (mmHandlePulse *PulseHandlerMock) HandlePulseBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmHandlePulse.beforeHandlePulseCounter)
}

// Calls returns a list of arguments used in each call to PulseHandlerMock.HandlePulse.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmHandlePulse *mPulseHandlerMockHandlePulse) Calls() []*PulseHandlerMockHandlePulseParams {
	mmHandlePulse.mutex.RLock()

	argCopy := make([]*PulseHandlerMockHandlePulseParams, len(mmHandlePulse.callArgs))
	copy(argCopy, mmHandlePulse.callArgs)

	mmHandlePulse.mutex.RUnlock()

	return argCopy
}

// MinimockHandlePulseDone returns true if the count of the HandlePulse invocations corresponds
// the number of defined expectations
func (m *PulseHandlerMock) MinimockHandlePulseDone() bool {
	for _, e := range m.HandlePulseMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.HandlePulseMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterHandlePulseCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcHandlePulse != nil && mm_atomic.LoadUint64(&m.afterHandlePulseCounter) < 1 {
		return false
	}
	return true
}

// MinimockHandlePulseInspect logs each unmet expectation
func (m *PulseHandlerMock) MinimockHandlePulseInspect() {
	for _, e := range m.HandlePulseMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to PulseHandlerMock.HandlePulse with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.HandlePulseMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterHandlePulseCounter) < 1 {
		if m.HandlePulseMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to PulseHandlerMock.HandlePulse")
		} else {
			m.t.Errorf("Expected call to PulseHandlerMock.HandlePulse with params: %#v", *m.HandlePulseMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcHandlePulse != nil && mm_atomic.LoadUint64(&m.afterHandlePulseCounter) < 1 {
		m.t.Error("Expected call to PulseHandlerMock.HandlePulse")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *PulseHandlerMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockHandlePulseInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *PulseHandlerMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *PulseHandlerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockHandlePulseDone()
}
