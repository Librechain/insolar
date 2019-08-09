package logicrunner

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock"

	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
)

// LogicExecutorMock implements LogicExecutor
type LogicExecutorMock struct {
	t minimock.Tester

	funcExecute          func(ctx context.Context, transcript *common.Transcript) (r1 artifacts.RequestResult, err error)
	inspectFuncExecute   func(ctx context.Context, transcript *common.Transcript)
	afterExecuteCounter  uint64
	beforeExecuteCounter uint64
	ExecuteMock          mLogicExecutorMockExecute

	funcExecuteConstructor          func(ctx context.Context, transcript *common.Transcript) (r1 artifacts.RequestResult, err error)
	inspectFuncExecuteConstructor   func(ctx context.Context, transcript *common.Transcript)
	afterExecuteConstructorCounter  uint64
	beforeExecuteConstructorCounter uint64
	ExecuteConstructorMock          mLogicExecutorMockExecuteConstructor

	funcExecuteMethod          func(ctx context.Context, transcript *common.Transcript) (r1 artifacts.RequestResult, err error)
	inspectFuncExecuteMethod   func(ctx context.Context, transcript *common.Transcript)
	afterExecuteMethodCounter  uint64
	beforeExecuteMethodCounter uint64
	ExecuteMethodMock          mLogicExecutorMockExecuteMethod
}

// NewLogicExecutorMock returns a mock for LogicExecutor
func NewLogicExecutorMock(t minimock.Tester) *LogicExecutorMock {
	m := &LogicExecutorMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ExecuteMock = mLogicExecutorMockExecute{mock: m}
	m.ExecuteMock.callArgs = []*LogicExecutorMockExecuteParams{}

	m.ExecuteConstructorMock = mLogicExecutorMockExecuteConstructor{mock: m}
	m.ExecuteConstructorMock.callArgs = []*LogicExecutorMockExecuteConstructorParams{}

	m.ExecuteMethodMock = mLogicExecutorMockExecuteMethod{mock: m}
	m.ExecuteMethodMock.callArgs = []*LogicExecutorMockExecuteMethodParams{}

	return m
}

type mLogicExecutorMockExecute struct {
	mock               *LogicExecutorMock
	defaultExpectation *LogicExecutorMockExecuteExpectation
	expectations       []*LogicExecutorMockExecuteExpectation

	callArgs []*LogicExecutorMockExecuteParams
	mutex    sync.RWMutex
}

// LogicExecutorMockExecuteExpectation specifies expectation struct of the LogicExecutor.Execute
type LogicExecutorMockExecuteExpectation struct {
	mock    *LogicExecutorMock
	params  *LogicExecutorMockExecuteParams
	results *LogicExecutorMockExecuteResults
	Counter uint64
}

// LogicExecutorMockExecuteParams contains parameters of the LogicExecutor.Execute
type LogicExecutorMockExecuteParams struct {
	ctx        context.Context
	transcript *common.Transcript
}

// LogicExecutorMockExecuteResults contains results of the LogicExecutor.Execute
type LogicExecutorMockExecuteResults struct {
	r1  artifacts.RequestResult
	err error
}

// Expect sets up expected params for LogicExecutor.Execute
func (mmExecute *mLogicExecutorMockExecute) Expect(ctx context.Context, transcript *common.Transcript) *mLogicExecutorMockExecute {
	if mmExecute.mock.funcExecute != nil {
		mmExecute.mock.t.Fatalf("LogicExecutorMock.Execute mock is already set by Set")
	}

	if mmExecute.defaultExpectation == nil {
		mmExecute.defaultExpectation = &LogicExecutorMockExecuteExpectation{}
	}

	mmExecute.defaultExpectation.params = &LogicExecutorMockExecuteParams{ctx, transcript}
	for _, e := range mmExecute.expectations {
		if minimock.Equal(e.params, mmExecute.defaultExpectation.params) {
			mmExecute.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmExecute.defaultExpectation.params)
		}
	}

	return mmExecute
}

// Inspect accepts an inspector function that has same arguments as the LogicExecutor.Execute
func (mmExecute *mLogicExecutorMockExecute) Inspect(f func(ctx context.Context, transcript *common.Transcript)) *mLogicExecutorMockExecute {
	if mmExecute.mock.inspectFuncExecute != nil {
		mmExecute.mock.t.Fatalf("Inspect function is already set for LogicExecutorMock.Execute")
	}

	mmExecute.mock.inspectFuncExecute = f

	return mmExecute
}

// Return sets up results that will be returned by LogicExecutor.Execute
func (mmExecute *mLogicExecutorMockExecute) Return(r1 artifacts.RequestResult, err error) *LogicExecutorMock {
	if mmExecute.mock.funcExecute != nil {
		mmExecute.mock.t.Fatalf("LogicExecutorMock.Execute mock is already set by Set")
	}

	if mmExecute.defaultExpectation == nil {
		mmExecute.defaultExpectation = &LogicExecutorMockExecuteExpectation{mock: mmExecute.mock}
	}
	mmExecute.defaultExpectation.results = &LogicExecutorMockExecuteResults{r1, err}
	return mmExecute.mock
}

//Set uses given function f to mock the LogicExecutor.Execute method
func (mmExecute *mLogicExecutorMockExecute) Set(f func(ctx context.Context, transcript *common.Transcript) (r1 artifacts.RequestResult, err error)) *LogicExecutorMock {
	if mmExecute.defaultExpectation != nil {
		mmExecute.mock.t.Fatalf("Default expectation is already set for the LogicExecutor.Execute method")
	}

	if len(mmExecute.expectations) > 0 {
		mmExecute.mock.t.Fatalf("Some expectations are already set for the LogicExecutor.Execute method")
	}

	mmExecute.mock.funcExecute = f
	return mmExecute.mock
}

// When sets expectation for the LogicExecutor.Execute which will trigger the result defined by the following
// Then helper
func (mmExecute *mLogicExecutorMockExecute) When(ctx context.Context, transcript *common.Transcript) *LogicExecutorMockExecuteExpectation {
	if mmExecute.mock.funcExecute != nil {
		mmExecute.mock.t.Fatalf("LogicExecutorMock.Execute mock is already set by Set")
	}

	expectation := &LogicExecutorMockExecuteExpectation{
		mock:   mmExecute.mock,
		params: &LogicExecutorMockExecuteParams{ctx, transcript},
	}
	mmExecute.expectations = append(mmExecute.expectations, expectation)
	return expectation
}

// Then sets up LogicExecutor.Execute return parameters for the expectation previously defined by the When method
func (e *LogicExecutorMockExecuteExpectation) Then(r1 artifacts.RequestResult, err error) *LogicExecutorMock {
	e.results = &LogicExecutorMockExecuteResults{r1, err}
	return e.mock
}

// Execute implements LogicExecutor
func (mmExecute *LogicExecutorMock) Execute(ctx context.Context, transcript *common.Transcript) (r1 artifacts.RequestResult, err error) {
	mm_atomic.AddUint64(&mmExecute.beforeExecuteCounter, 1)
	defer mm_atomic.AddUint64(&mmExecute.afterExecuteCounter, 1)

	if mmExecute.inspectFuncExecute != nil {
		mmExecute.inspectFuncExecute(ctx, transcript)
	}

	params := &LogicExecutorMockExecuteParams{ctx, transcript}

	// Record call args
	mmExecute.ExecuteMock.mutex.Lock()
	mmExecute.ExecuteMock.callArgs = append(mmExecute.ExecuteMock.callArgs, params)
	mmExecute.ExecuteMock.mutex.Unlock()

	for _, e := range mmExecute.ExecuteMock.expectations {
		if minimock.Equal(e.params, params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.r1, e.results.err
		}
	}

	if mmExecute.ExecuteMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmExecute.ExecuteMock.defaultExpectation.Counter, 1)
		want := mmExecute.ExecuteMock.defaultExpectation.params
		got := LogicExecutorMockExecuteParams{ctx, transcript}
		if want != nil && !minimock.Equal(*want, got) {
			mmExecute.t.Errorf("LogicExecutorMock.Execute got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := mmExecute.ExecuteMock.defaultExpectation.results
		if results == nil {
			mmExecute.t.Fatal("No results are set for the LogicExecutorMock.Execute")
		}
		return (*results).r1, (*results).err
	}
	if mmExecute.funcExecute != nil {
		return mmExecute.funcExecute(ctx, transcript)
	}
	mmExecute.t.Fatalf("Unexpected call to LogicExecutorMock.Execute. %v %v", ctx, transcript)
	return
}

// ExecuteAfterCounter returns a count of finished LogicExecutorMock.Execute invocations
func (mmExecute *LogicExecutorMock) ExecuteAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmExecute.afterExecuteCounter)
}

// ExecuteBeforeCounter returns a count of LogicExecutorMock.Execute invocations
func (mmExecute *LogicExecutorMock) ExecuteBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmExecute.beforeExecuteCounter)
}

// Calls returns a list of arguments used in each call to LogicExecutorMock.Execute.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmExecute *mLogicExecutorMockExecute) Calls() []*LogicExecutorMockExecuteParams {
	mmExecute.mutex.RLock()

	argCopy := make([]*LogicExecutorMockExecuteParams, len(mmExecute.callArgs))
	copy(argCopy, mmExecute.callArgs)

	mmExecute.mutex.RUnlock()

	return argCopy
}

// MinimockExecuteDone returns true if the count of the Execute invocations corresponds
// the number of defined expectations
func (m *LogicExecutorMock) MinimockExecuteDone() bool {
	for _, e := range m.ExecuteMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ExecuteMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterExecuteCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcExecute != nil && mm_atomic.LoadUint64(&m.afterExecuteCounter) < 1 {
		return false
	}
	return true
}

// MinimockExecuteInspect logs each unmet expectation
func (m *LogicExecutorMock) MinimockExecuteInspect() {
	for _, e := range m.ExecuteMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to LogicExecutorMock.Execute with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ExecuteMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterExecuteCounter) < 1 {
		if m.ExecuteMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to LogicExecutorMock.Execute")
		} else {
			m.t.Errorf("Expected call to LogicExecutorMock.Execute with params: %#v", *m.ExecuteMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcExecute != nil && mm_atomic.LoadUint64(&m.afterExecuteCounter) < 1 {
		m.t.Error("Expected call to LogicExecutorMock.Execute")
	}
}

type mLogicExecutorMockExecuteConstructor struct {
	mock               *LogicExecutorMock
	defaultExpectation *LogicExecutorMockExecuteConstructorExpectation
	expectations       []*LogicExecutorMockExecuteConstructorExpectation

	callArgs []*LogicExecutorMockExecuteConstructorParams
	mutex    sync.RWMutex
}

// LogicExecutorMockExecuteConstructorExpectation specifies expectation struct of the LogicExecutor.ExecuteConstructor
type LogicExecutorMockExecuteConstructorExpectation struct {
	mock    *LogicExecutorMock
	params  *LogicExecutorMockExecuteConstructorParams
	results *LogicExecutorMockExecuteConstructorResults
	Counter uint64
}

// LogicExecutorMockExecuteConstructorParams contains parameters of the LogicExecutor.ExecuteConstructor
type LogicExecutorMockExecuteConstructorParams struct {
	ctx        context.Context
	transcript *common.Transcript
}

// LogicExecutorMockExecuteConstructorResults contains results of the LogicExecutor.ExecuteConstructor
type LogicExecutorMockExecuteConstructorResults struct {
	r1  artifacts.RequestResult
	err error
}

// Expect sets up expected params for LogicExecutor.ExecuteConstructor
func (mmExecuteConstructor *mLogicExecutorMockExecuteConstructor) Expect(ctx context.Context, transcript *common.Transcript) *mLogicExecutorMockExecuteConstructor {
	if mmExecuteConstructor.mock.funcExecuteConstructor != nil {
		mmExecuteConstructor.mock.t.Fatalf("LogicExecutorMock.ExecuteConstructor mock is already set by Set")
	}

	if mmExecuteConstructor.defaultExpectation == nil {
		mmExecuteConstructor.defaultExpectation = &LogicExecutorMockExecuteConstructorExpectation{}
	}

	mmExecuteConstructor.defaultExpectation.params = &LogicExecutorMockExecuteConstructorParams{ctx, transcript}
	for _, e := range mmExecuteConstructor.expectations {
		if minimock.Equal(e.params, mmExecuteConstructor.defaultExpectation.params) {
			mmExecuteConstructor.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmExecuteConstructor.defaultExpectation.params)
		}
	}

	return mmExecuteConstructor
}

// Inspect accepts an inspector function that has same arguments as the LogicExecutor.ExecuteConstructor
func (mmExecuteConstructor *mLogicExecutorMockExecuteConstructor) Inspect(f func(ctx context.Context, transcript *common.Transcript)) *mLogicExecutorMockExecuteConstructor {
	if mmExecuteConstructor.mock.inspectFuncExecuteConstructor != nil {
		mmExecuteConstructor.mock.t.Fatalf("Inspect function is already set for LogicExecutorMock.ExecuteConstructor")
	}

	mmExecuteConstructor.mock.inspectFuncExecuteConstructor = f

	return mmExecuteConstructor
}

// Return sets up results that will be returned by LogicExecutor.ExecuteConstructor
func (mmExecuteConstructor *mLogicExecutorMockExecuteConstructor) Return(r1 artifacts.RequestResult, err error) *LogicExecutorMock {
	if mmExecuteConstructor.mock.funcExecuteConstructor != nil {
		mmExecuteConstructor.mock.t.Fatalf("LogicExecutorMock.ExecuteConstructor mock is already set by Set")
	}

	if mmExecuteConstructor.defaultExpectation == nil {
		mmExecuteConstructor.defaultExpectation = &LogicExecutorMockExecuteConstructorExpectation{mock: mmExecuteConstructor.mock}
	}
	mmExecuteConstructor.defaultExpectation.results = &LogicExecutorMockExecuteConstructorResults{r1, err}
	return mmExecuteConstructor.mock
}

//Set uses given function f to mock the LogicExecutor.ExecuteConstructor method
func (mmExecuteConstructor *mLogicExecutorMockExecuteConstructor) Set(f func(ctx context.Context, transcript *common.Transcript) (r1 artifacts.RequestResult, err error)) *LogicExecutorMock {
	if mmExecuteConstructor.defaultExpectation != nil {
		mmExecuteConstructor.mock.t.Fatalf("Default expectation is already set for the LogicExecutor.ExecuteConstructor method")
	}

	if len(mmExecuteConstructor.expectations) > 0 {
		mmExecuteConstructor.mock.t.Fatalf("Some expectations are already set for the LogicExecutor.ExecuteConstructor method")
	}

	mmExecuteConstructor.mock.funcExecuteConstructor = f
	return mmExecuteConstructor.mock
}

// When sets expectation for the LogicExecutor.ExecuteConstructor which will trigger the result defined by the following
// Then helper
func (mmExecuteConstructor *mLogicExecutorMockExecuteConstructor) When(ctx context.Context, transcript *common.Transcript) *LogicExecutorMockExecuteConstructorExpectation {
	if mmExecuteConstructor.mock.funcExecuteConstructor != nil {
		mmExecuteConstructor.mock.t.Fatalf("LogicExecutorMock.ExecuteConstructor mock is already set by Set")
	}

	expectation := &LogicExecutorMockExecuteConstructorExpectation{
		mock:   mmExecuteConstructor.mock,
		params: &LogicExecutorMockExecuteConstructorParams{ctx, transcript},
	}
	mmExecuteConstructor.expectations = append(mmExecuteConstructor.expectations, expectation)
	return expectation
}

// Then sets up LogicExecutor.ExecuteConstructor return parameters for the expectation previously defined by the When method
func (e *LogicExecutorMockExecuteConstructorExpectation) Then(r1 artifacts.RequestResult, err error) *LogicExecutorMock {
	e.results = &LogicExecutorMockExecuteConstructorResults{r1, err}
	return e.mock
}

// ExecuteConstructor implements LogicExecutor
func (mmExecuteConstructor *LogicExecutorMock) ExecuteConstructor(ctx context.Context, transcript *common.Transcript) (r1 artifacts.RequestResult, err error) {
	mm_atomic.AddUint64(&mmExecuteConstructor.beforeExecuteConstructorCounter, 1)
	defer mm_atomic.AddUint64(&mmExecuteConstructor.afterExecuteConstructorCounter, 1)

	if mmExecuteConstructor.inspectFuncExecuteConstructor != nil {
		mmExecuteConstructor.inspectFuncExecuteConstructor(ctx, transcript)
	}

	params := &LogicExecutorMockExecuteConstructorParams{ctx, transcript}

	// Record call args
	mmExecuteConstructor.ExecuteConstructorMock.mutex.Lock()
	mmExecuteConstructor.ExecuteConstructorMock.callArgs = append(mmExecuteConstructor.ExecuteConstructorMock.callArgs, params)
	mmExecuteConstructor.ExecuteConstructorMock.mutex.Unlock()

	for _, e := range mmExecuteConstructor.ExecuteConstructorMock.expectations {
		if minimock.Equal(e.params, params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.r1, e.results.err
		}
	}

	if mmExecuteConstructor.ExecuteConstructorMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmExecuteConstructor.ExecuteConstructorMock.defaultExpectation.Counter, 1)
		want := mmExecuteConstructor.ExecuteConstructorMock.defaultExpectation.params
		got := LogicExecutorMockExecuteConstructorParams{ctx, transcript}
		if want != nil && !minimock.Equal(*want, got) {
			mmExecuteConstructor.t.Errorf("LogicExecutorMock.ExecuteConstructor got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := mmExecuteConstructor.ExecuteConstructorMock.defaultExpectation.results
		if results == nil {
			mmExecuteConstructor.t.Fatal("No results are set for the LogicExecutorMock.ExecuteConstructor")
		}
		return (*results).r1, (*results).err
	}
	if mmExecuteConstructor.funcExecuteConstructor != nil {
		return mmExecuteConstructor.funcExecuteConstructor(ctx, transcript)
	}
	mmExecuteConstructor.t.Fatalf("Unexpected call to LogicExecutorMock.ExecuteConstructor. %v %v", ctx, transcript)
	return
}

// ExecuteConstructorAfterCounter returns a count of finished LogicExecutorMock.ExecuteConstructor invocations
func (mmExecuteConstructor *LogicExecutorMock) ExecuteConstructorAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmExecuteConstructor.afterExecuteConstructorCounter)
}

// ExecuteConstructorBeforeCounter returns a count of LogicExecutorMock.ExecuteConstructor invocations
func (mmExecuteConstructor *LogicExecutorMock) ExecuteConstructorBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmExecuteConstructor.beforeExecuteConstructorCounter)
}

// Calls returns a list of arguments used in each call to LogicExecutorMock.ExecuteConstructor.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmExecuteConstructor *mLogicExecutorMockExecuteConstructor) Calls() []*LogicExecutorMockExecuteConstructorParams {
	mmExecuteConstructor.mutex.RLock()

	argCopy := make([]*LogicExecutorMockExecuteConstructorParams, len(mmExecuteConstructor.callArgs))
	copy(argCopy, mmExecuteConstructor.callArgs)

	mmExecuteConstructor.mutex.RUnlock()

	return argCopy
}

// MinimockExecuteConstructorDone returns true if the count of the ExecuteConstructor invocations corresponds
// the number of defined expectations
func (m *LogicExecutorMock) MinimockExecuteConstructorDone() bool {
	for _, e := range m.ExecuteConstructorMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ExecuteConstructorMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterExecuteConstructorCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcExecuteConstructor != nil && mm_atomic.LoadUint64(&m.afterExecuteConstructorCounter) < 1 {
		return false
	}
	return true
}

// MinimockExecuteConstructorInspect logs each unmet expectation
func (m *LogicExecutorMock) MinimockExecuteConstructorInspect() {
	for _, e := range m.ExecuteConstructorMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to LogicExecutorMock.ExecuteConstructor with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ExecuteConstructorMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterExecuteConstructorCounter) < 1 {
		if m.ExecuteConstructorMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to LogicExecutorMock.ExecuteConstructor")
		} else {
			m.t.Errorf("Expected call to LogicExecutorMock.ExecuteConstructor with params: %#v", *m.ExecuteConstructorMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcExecuteConstructor != nil && mm_atomic.LoadUint64(&m.afterExecuteConstructorCounter) < 1 {
		m.t.Error("Expected call to LogicExecutorMock.ExecuteConstructor")
	}
}

type mLogicExecutorMockExecuteMethod struct {
	mock               *LogicExecutorMock
	defaultExpectation *LogicExecutorMockExecuteMethodExpectation
	expectations       []*LogicExecutorMockExecuteMethodExpectation

	callArgs []*LogicExecutorMockExecuteMethodParams
	mutex    sync.RWMutex
}

// LogicExecutorMockExecuteMethodExpectation specifies expectation struct of the LogicExecutor.ExecuteMethod
type LogicExecutorMockExecuteMethodExpectation struct {
	mock    *LogicExecutorMock
	params  *LogicExecutorMockExecuteMethodParams
	results *LogicExecutorMockExecuteMethodResults
	Counter uint64
}

// LogicExecutorMockExecuteMethodParams contains parameters of the LogicExecutor.ExecuteMethod
type LogicExecutorMockExecuteMethodParams struct {
	ctx        context.Context
	transcript *common.Transcript
}

// LogicExecutorMockExecuteMethodResults contains results of the LogicExecutor.ExecuteMethod
type LogicExecutorMockExecuteMethodResults struct {
	r1  artifacts.RequestResult
	err error
}

// Expect sets up expected params for LogicExecutor.ExecuteMethod
func (mmExecuteMethod *mLogicExecutorMockExecuteMethod) Expect(ctx context.Context, transcript *common.Transcript) *mLogicExecutorMockExecuteMethod {
	if mmExecuteMethod.mock.funcExecuteMethod != nil {
		mmExecuteMethod.mock.t.Fatalf("LogicExecutorMock.ExecuteMethod mock is already set by Set")
	}

	if mmExecuteMethod.defaultExpectation == nil {
		mmExecuteMethod.defaultExpectation = &LogicExecutorMockExecuteMethodExpectation{}
	}

	mmExecuteMethod.defaultExpectation.params = &LogicExecutorMockExecuteMethodParams{ctx, transcript}
	for _, e := range mmExecuteMethod.expectations {
		if minimock.Equal(e.params, mmExecuteMethod.defaultExpectation.params) {
			mmExecuteMethod.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmExecuteMethod.defaultExpectation.params)
		}
	}

	return mmExecuteMethod
}

// Inspect accepts an inspector function that has same arguments as the LogicExecutor.ExecuteMethod
func (mmExecuteMethod *mLogicExecutorMockExecuteMethod) Inspect(f func(ctx context.Context, transcript *common.Transcript)) *mLogicExecutorMockExecuteMethod {
	if mmExecuteMethod.mock.inspectFuncExecuteMethod != nil {
		mmExecuteMethod.mock.t.Fatalf("Inspect function is already set for LogicExecutorMock.ExecuteMethod")
	}

	mmExecuteMethod.mock.inspectFuncExecuteMethod = f

	return mmExecuteMethod
}

// Return sets up results that will be returned by LogicExecutor.ExecuteMethod
func (mmExecuteMethod *mLogicExecutorMockExecuteMethod) Return(r1 artifacts.RequestResult, err error) *LogicExecutorMock {
	if mmExecuteMethod.mock.funcExecuteMethod != nil {
		mmExecuteMethod.mock.t.Fatalf("LogicExecutorMock.ExecuteMethod mock is already set by Set")
	}

	if mmExecuteMethod.defaultExpectation == nil {
		mmExecuteMethod.defaultExpectation = &LogicExecutorMockExecuteMethodExpectation{mock: mmExecuteMethod.mock}
	}
	mmExecuteMethod.defaultExpectation.results = &LogicExecutorMockExecuteMethodResults{r1, err}
	return mmExecuteMethod.mock
}

//Set uses given function f to mock the LogicExecutor.ExecuteMethod method
func (mmExecuteMethod *mLogicExecutorMockExecuteMethod) Set(f func(ctx context.Context, transcript *common.Transcript) (r1 artifacts.RequestResult, err error)) *LogicExecutorMock {
	if mmExecuteMethod.defaultExpectation != nil {
		mmExecuteMethod.mock.t.Fatalf("Default expectation is already set for the LogicExecutor.ExecuteMethod method")
	}

	if len(mmExecuteMethod.expectations) > 0 {
		mmExecuteMethod.mock.t.Fatalf("Some expectations are already set for the LogicExecutor.ExecuteMethod method")
	}

	mmExecuteMethod.mock.funcExecuteMethod = f
	return mmExecuteMethod.mock
}

// When sets expectation for the LogicExecutor.ExecuteMethod which will trigger the result defined by the following
// Then helper
func (mmExecuteMethod *mLogicExecutorMockExecuteMethod) When(ctx context.Context, transcript *common.Transcript) *LogicExecutorMockExecuteMethodExpectation {
	if mmExecuteMethod.mock.funcExecuteMethod != nil {
		mmExecuteMethod.mock.t.Fatalf("LogicExecutorMock.ExecuteMethod mock is already set by Set")
	}

	expectation := &LogicExecutorMockExecuteMethodExpectation{
		mock:   mmExecuteMethod.mock,
		params: &LogicExecutorMockExecuteMethodParams{ctx, transcript},
	}
	mmExecuteMethod.expectations = append(mmExecuteMethod.expectations, expectation)
	return expectation
}

// Then sets up LogicExecutor.ExecuteMethod return parameters for the expectation previously defined by the When method
func (e *LogicExecutorMockExecuteMethodExpectation) Then(r1 artifacts.RequestResult, err error) *LogicExecutorMock {
	e.results = &LogicExecutorMockExecuteMethodResults{r1, err}
	return e.mock
}

// ExecuteMethod implements LogicExecutor
func (mmExecuteMethod *LogicExecutorMock) ExecuteMethod(ctx context.Context, transcript *common.Transcript) (r1 artifacts.RequestResult, err error) {
	mm_atomic.AddUint64(&mmExecuteMethod.beforeExecuteMethodCounter, 1)
	defer mm_atomic.AddUint64(&mmExecuteMethod.afterExecuteMethodCounter, 1)

	if mmExecuteMethod.inspectFuncExecuteMethod != nil {
		mmExecuteMethod.inspectFuncExecuteMethod(ctx, transcript)
	}

	params := &LogicExecutorMockExecuteMethodParams{ctx, transcript}

	// Record call args
	mmExecuteMethod.ExecuteMethodMock.mutex.Lock()
	mmExecuteMethod.ExecuteMethodMock.callArgs = append(mmExecuteMethod.ExecuteMethodMock.callArgs, params)
	mmExecuteMethod.ExecuteMethodMock.mutex.Unlock()

	for _, e := range mmExecuteMethod.ExecuteMethodMock.expectations {
		if minimock.Equal(e.params, params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.r1, e.results.err
		}
	}

	if mmExecuteMethod.ExecuteMethodMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmExecuteMethod.ExecuteMethodMock.defaultExpectation.Counter, 1)
		want := mmExecuteMethod.ExecuteMethodMock.defaultExpectation.params
		got := LogicExecutorMockExecuteMethodParams{ctx, transcript}
		if want != nil && !minimock.Equal(*want, got) {
			mmExecuteMethod.t.Errorf("LogicExecutorMock.ExecuteMethod got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := mmExecuteMethod.ExecuteMethodMock.defaultExpectation.results
		if results == nil {
			mmExecuteMethod.t.Fatal("No results are set for the LogicExecutorMock.ExecuteMethod")
		}
		return (*results).r1, (*results).err
	}
	if mmExecuteMethod.funcExecuteMethod != nil {
		return mmExecuteMethod.funcExecuteMethod(ctx, transcript)
	}
	mmExecuteMethod.t.Fatalf("Unexpected call to LogicExecutorMock.ExecuteMethod. %v %v", ctx, transcript)
	return
}

// ExecuteMethodAfterCounter returns a count of finished LogicExecutorMock.ExecuteMethod invocations
func (mmExecuteMethod *LogicExecutorMock) ExecuteMethodAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmExecuteMethod.afterExecuteMethodCounter)
}

// ExecuteMethodBeforeCounter returns a count of LogicExecutorMock.ExecuteMethod invocations
func (mmExecuteMethod *LogicExecutorMock) ExecuteMethodBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmExecuteMethod.beforeExecuteMethodCounter)
}

// Calls returns a list of arguments used in each call to LogicExecutorMock.ExecuteMethod.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmExecuteMethod *mLogicExecutorMockExecuteMethod) Calls() []*LogicExecutorMockExecuteMethodParams {
	mmExecuteMethod.mutex.RLock()

	argCopy := make([]*LogicExecutorMockExecuteMethodParams, len(mmExecuteMethod.callArgs))
	copy(argCopy, mmExecuteMethod.callArgs)

	mmExecuteMethod.mutex.RUnlock()

	return argCopy
}

// MinimockExecuteMethodDone returns true if the count of the ExecuteMethod invocations corresponds
// the number of defined expectations
func (m *LogicExecutorMock) MinimockExecuteMethodDone() bool {
	for _, e := range m.ExecuteMethodMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ExecuteMethodMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterExecuteMethodCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcExecuteMethod != nil && mm_atomic.LoadUint64(&m.afterExecuteMethodCounter) < 1 {
		return false
	}
	return true
}

// MinimockExecuteMethodInspect logs each unmet expectation
func (m *LogicExecutorMock) MinimockExecuteMethodInspect() {
	for _, e := range m.ExecuteMethodMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to LogicExecutorMock.ExecuteMethod with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ExecuteMethodMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterExecuteMethodCounter) < 1 {
		if m.ExecuteMethodMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to LogicExecutorMock.ExecuteMethod")
		} else {
			m.t.Errorf("Expected call to LogicExecutorMock.ExecuteMethod with params: %#v", *m.ExecuteMethodMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcExecuteMethod != nil && mm_atomic.LoadUint64(&m.afterExecuteMethodCounter) < 1 {
		m.t.Error("Expected call to LogicExecutorMock.ExecuteMethod")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *LogicExecutorMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockExecuteInspect()

		m.MinimockExecuteConstructorInspect()

		m.MinimockExecuteMethodInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *LogicExecutorMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *LogicExecutorMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockExecuteDone() &&
		m.MinimockExecuteConstructorDone() &&
		m.MinimockExecuteMethodDone()
}
