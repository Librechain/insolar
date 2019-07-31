package testutils

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock"
	mm_insolar "github.com/insolar/insolar/insolar"
)

// CertificateManagerMock implements insolar.CertificateManager
type CertificateManagerMock struct {
	t minimock.Tester

	funcGetCertificate          func() (c1 mm_insolar.Certificate)
	inspectFuncGetCertificate   func()
	afterGetCertificateCounter  uint64
	beforeGetCertificateCounter uint64
	GetCertificateMock          mCertificateManagerMockGetCertificate

	funcNewUnsignedCertificate          func(pKey string, role string, nodeRef string) (c1 mm_insolar.Certificate, err error)
	inspectFuncNewUnsignedCertificate   func(pKey string, role string, nodeRef string)
	afterNewUnsignedCertificateCounter  uint64
	beforeNewUnsignedCertificateCounter uint64
	NewUnsignedCertificateMock          mCertificateManagerMockNewUnsignedCertificate

	funcVerifyAuthorizationCertificate          func(authCert mm_insolar.AuthorizationCertificate) (b1 bool, err error)
	inspectFuncVerifyAuthorizationCertificate   func(authCert mm_insolar.AuthorizationCertificate)
	afterVerifyAuthorizationCertificateCounter  uint64
	beforeVerifyAuthorizationCertificateCounter uint64
	VerifyAuthorizationCertificateMock          mCertificateManagerMockVerifyAuthorizationCertificate
}

// NewCertificateManagerMock returns a mock for insolar.CertificateManager
func NewCertificateManagerMock(t minimock.Tester) *CertificateManagerMock {
	m := &CertificateManagerMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetCertificateMock = mCertificateManagerMockGetCertificate{mock: m}

	m.NewUnsignedCertificateMock = mCertificateManagerMockNewUnsignedCertificate{mock: m}
	m.NewUnsignedCertificateMock.callArgs = []*CertificateManagerMockNewUnsignedCertificateParams{}

	m.VerifyAuthorizationCertificateMock = mCertificateManagerMockVerifyAuthorizationCertificate{mock: m}
	m.VerifyAuthorizationCertificateMock.callArgs = []*CertificateManagerMockVerifyAuthorizationCertificateParams{}

	return m
}

type mCertificateManagerMockGetCertificate struct {
	mock               *CertificateManagerMock
	defaultExpectation *CertificateManagerMockGetCertificateExpectation
	expectations       []*CertificateManagerMockGetCertificateExpectation
}

// CertificateManagerMockGetCertificateExpectation specifies expectation struct of the CertificateManager.GetCertificate
type CertificateManagerMockGetCertificateExpectation struct {
	mock *CertificateManagerMock

	results *CertificateManagerMockGetCertificateResults
	Counter uint64
}

// CertificateManagerMockGetCertificateResults contains results of the CertificateManager.GetCertificate
type CertificateManagerMockGetCertificateResults struct {
	c1 mm_insolar.Certificate
}

// Expect sets up expected params for CertificateManager.GetCertificate
func (mmGetCertificate *mCertificateManagerMockGetCertificate) Expect() *mCertificateManagerMockGetCertificate {
	if mmGetCertificate.mock.funcGetCertificate != nil {
		mmGetCertificate.mock.t.Fatalf("CertificateManagerMock.GetCertificate mock is already set by Set")
	}

	if mmGetCertificate.defaultExpectation == nil {
		mmGetCertificate.defaultExpectation = &CertificateManagerMockGetCertificateExpectation{}
	}

	return mmGetCertificate
}

// Inspect accepts an inspector function that has same arguments as the CertificateManager.GetCertificate
func (mmGetCertificate *mCertificateManagerMockGetCertificate) Inspect(f func()) *mCertificateManagerMockGetCertificate {
	if mmGetCertificate.mock.inspectFuncGetCertificate != nil {
		mmGetCertificate.mock.t.Fatalf("Inspect function is already set for CertificateManagerMock.GetCertificate")
	}

	mmGetCertificate.mock.inspectFuncGetCertificate = f

	return mmGetCertificate
}

// Return sets up results that will be returned by CertificateManager.GetCertificate
func (mmGetCertificate *mCertificateManagerMockGetCertificate) Return(c1 mm_insolar.Certificate) *CertificateManagerMock {
	if mmGetCertificate.mock.funcGetCertificate != nil {
		mmGetCertificate.mock.t.Fatalf("CertificateManagerMock.GetCertificate mock is already set by Set")
	}

	if mmGetCertificate.defaultExpectation == nil {
		mmGetCertificate.defaultExpectation = &CertificateManagerMockGetCertificateExpectation{mock: mmGetCertificate.mock}
	}
	mmGetCertificate.defaultExpectation.results = &CertificateManagerMockGetCertificateResults{c1}
	return mmGetCertificate.mock
}

//Set uses given function f to mock the CertificateManager.GetCertificate method
func (mmGetCertificate *mCertificateManagerMockGetCertificate) Set(f func() (c1 mm_insolar.Certificate)) *CertificateManagerMock {
	if mmGetCertificate.defaultExpectation != nil {
		mmGetCertificate.mock.t.Fatalf("Default expectation is already set for the CertificateManager.GetCertificate method")
	}

	if len(mmGetCertificate.expectations) > 0 {
		mmGetCertificate.mock.t.Fatalf("Some expectations are already set for the CertificateManager.GetCertificate method")
	}

	mmGetCertificate.mock.funcGetCertificate = f
	return mmGetCertificate.mock
}

// GetCertificate implements insolar.CertificateManager
func (mmGetCertificate *CertificateManagerMock) GetCertificate() (c1 mm_insolar.Certificate) {
	mm_atomic.AddUint64(&mmGetCertificate.beforeGetCertificateCounter, 1)
	defer mm_atomic.AddUint64(&mmGetCertificate.afterGetCertificateCounter, 1)

	if mmGetCertificate.inspectFuncGetCertificate != nil {
		mmGetCertificate.inspectFuncGetCertificate()
	}

	if mmGetCertificate.GetCertificateMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetCertificate.GetCertificateMock.defaultExpectation.Counter, 1)

		results := mmGetCertificate.GetCertificateMock.defaultExpectation.results
		if results == nil {
			mmGetCertificate.t.Fatal("No results are set for the CertificateManagerMock.GetCertificate")
		}
		return (*results).c1
	}
	if mmGetCertificate.funcGetCertificate != nil {
		return mmGetCertificate.funcGetCertificate()
	}
	mmGetCertificate.t.Fatalf("Unexpected call to CertificateManagerMock.GetCertificate.")
	return
}

// GetCertificateAfterCounter returns a count of finished CertificateManagerMock.GetCertificate invocations
func (mmGetCertificate *CertificateManagerMock) GetCertificateAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetCertificate.afterGetCertificateCounter)
}

// GetCertificateBeforeCounter returns a count of CertificateManagerMock.GetCertificate invocations
func (mmGetCertificate *CertificateManagerMock) GetCertificateBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetCertificate.beforeGetCertificateCounter)
}

// MinimockGetCertificateDone returns true if the count of the GetCertificate invocations corresponds
// the number of defined expectations
func (m *CertificateManagerMock) MinimockGetCertificateDone() bool {
	for _, e := range m.GetCertificateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetCertificateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetCertificateCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetCertificate != nil && mm_atomic.LoadUint64(&m.afterGetCertificateCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetCertificateInspect logs each unmet expectation
func (m *CertificateManagerMock) MinimockGetCertificateInspect() {
	for _, e := range m.GetCertificateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to CertificateManagerMock.GetCertificate")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetCertificateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetCertificateCounter) < 1 {
		m.t.Error("Expected call to CertificateManagerMock.GetCertificate")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetCertificate != nil && mm_atomic.LoadUint64(&m.afterGetCertificateCounter) < 1 {
		m.t.Error("Expected call to CertificateManagerMock.GetCertificate")
	}
}

type mCertificateManagerMockNewUnsignedCertificate struct {
	mock               *CertificateManagerMock
	defaultExpectation *CertificateManagerMockNewUnsignedCertificateExpectation
	expectations       []*CertificateManagerMockNewUnsignedCertificateExpectation

	callArgs []*CertificateManagerMockNewUnsignedCertificateParams
	mutex    sync.RWMutex
}

// CertificateManagerMockNewUnsignedCertificateExpectation specifies expectation struct of the CertificateManager.NewUnsignedCertificate
type CertificateManagerMockNewUnsignedCertificateExpectation struct {
	mock    *CertificateManagerMock
	params  *CertificateManagerMockNewUnsignedCertificateParams
	results *CertificateManagerMockNewUnsignedCertificateResults
	Counter uint64
}

// CertificateManagerMockNewUnsignedCertificateParams contains parameters of the CertificateManager.NewUnsignedCertificate
type CertificateManagerMockNewUnsignedCertificateParams struct {
	pKey    string
	role    string
	nodeRef string
}

// CertificateManagerMockNewUnsignedCertificateResults contains results of the CertificateManager.NewUnsignedCertificate
type CertificateManagerMockNewUnsignedCertificateResults struct {
	c1  mm_insolar.Certificate
	err error
}

// Expect sets up expected params for CertificateManager.NewUnsignedCertificate
func (mmNewUnsignedCertificate *mCertificateManagerMockNewUnsignedCertificate) Expect(pKey string, role string, nodeRef string) *mCertificateManagerMockNewUnsignedCertificate {
	if mmNewUnsignedCertificate.mock.funcNewUnsignedCertificate != nil {
		mmNewUnsignedCertificate.mock.t.Fatalf("CertificateManagerMock.NewUnsignedCertificate mock is already set by Set")
	}

	if mmNewUnsignedCertificate.defaultExpectation == nil {
		mmNewUnsignedCertificate.defaultExpectation = &CertificateManagerMockNewUnsignedCertificateExpectation{}
	}

	mmNewUnsignedCertificate.defaultExpectation.params = &CertificateManagerMockNewUnsignedCertificateParams{pKey, role, nodeRef}
	for _, e := range mmNewUnsignedCertificate.expectations {
		if minimock.Equal(e.params, mmNewUnsignedCertificate.defaultExpectation.params) {
			mmNewUnsignedCertificate.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmNewUnsignedCertificate.defaultExpectation.params)
		}
	}

	return mmNewUnsignedCertificate
}

// Inspect accepts an inspector function that has same arguments as the CertificateManager.NewUnsignedCertificate
func (mmNewUnsignedCertificate *mCertificateManagerMockNewUnsignedCertificate) Inspect(f func(pKey string, role string, nodeRef string)) *mCertificateManagerMockNewUnsignedCertificate {
	if mmNewUnsignedCertificate.mock.inspectFuncNewUnsignedCertificate != nil {
		mmNewUnsignedCertificate.mock.t.Fatalf("Inspect function is already set for CertificateManagerMock.NewUnsignedCertificate")
	}

	mmNewUnsignedCertificate.mock.inspectFuncNewUnsignedCertificate = f

	return mmNewUnsignedCertificate
}

// Return sets up results that will be returned by CertificateManager.NewUnsignedCertificate
func (mmNewUnsignedCertificate *mCertificateManagerMockNewUnsignedCertificate) Return(c1 mm_insolar.Certificate, err error) *CertificateManagerMock {
	if mmNewUnsignedCertificate.mock.funcNewUnsignedCertificate != nil {
		mmNewUnsignedCertificate.mock.t.Fatalf("CertificateManagerMock.NewUnsignedCertificate mock is already set by Set")
	}

	if mmNewUnsignedCertificate.defaultExpectation == nil {
		mmNewUnsignedCertificate.defaultExpectation = &CertificateManagerMockNewUnsignedCertificateExpectation{mock: mmNewUnsignedCertificate.mock}
	}
	mmNewUnsignedCertificate.defaultExpectation.results = &CertificateManagerMockNewUnsignedCertificateResults{c1, err}
	return mmNewUnsignedCertificate.mock
}

//Set uses given function f to mock the CertificateManager.NewUnsignedCertificate method
func (mmNewUnsignedCertificate *mCertificateManagerMockNewUnsignedCertificate) Set(f func(pKey string, role string, nodeRef string) (c1 mm_insolar.Certificate, err error)) *CertificateManagerMock {
	if mmNewUnsignedCertificate.defaultExpectation != nil {
		mmNewUnsignedCertificate.mock.t.Fatalf("Default expectation is already set for the CertificateManager.NewUnsignedCertificate method")
	}

	if len(mmNewUnsignedCertificate.expectations) > 0 {
		mmNewUnsignedCertificate.mock.t.Fatalf("Some expectations are already set for the CertificateManager.NewUnsignedCertificate method")
	}

	mmNewUnsignedCertificate.mock.funcNewUnsignedCertificate = f
	return mmNewUnsignedCertificate.mock
}

// When sets expectation for the CertificateManager.NewUnsignedCertificate which will trigger the result defined by the following
// Then helper
func (mmNewUnsignedCertificate *mCertificateManagerMockNewUnsignedCertificate) When(pKey string, role string, nodeRef string) *CertificateManagerMockNewUnsignedCertificateExpectation {
	if mmNewUnsignedCertificate.mock.funcNewUnsignedCertificate != nil {
		mmNewUnsignedCertificate.mock.t.Fatalf("CertificateManagerMock.NewUnsignedCertificate mock is already set by Set")
	}

	expectation := &CertificateManagerMockNewUnsignedCertificateExpectation{
		mock:   mmNewUnsignedCertificate.mock,
		params: &CertificateManagerMockNewUnsignedCertificateParams{pKey, role, nodeRef},
	}
	mmNewUnsignedCertificate.expectations = append(mmNewUnsignedCertificate.expectations, expectation)
	return expectation
}

// Then sets up CertificateManager.NewUnsignedCertificate return parameters for the expectation previously defined by the When method
func (e *CertificateManagerMockNewUnsignedCertificateExpectation) Then(c1 mm_insolar.Certificate, err error) *CertificateManagerMock {
	e.results = &CertificateManagerMockNewUnsignedCertificateResults{c1, err}
	return e.mock
}

// NewUnsignedCertificate implements insolar.CertificateManager
func (mmNewUnsignedCertificate *CertificateManagerMock) NewUnsignedCertificate(pKey string, role string, nodeRef string) (c1 mm_insolar.Certificate, err error) {
	mm_atomic.AddUint64(&mmNewUnsignedCertificate.beforeNewUnsignedCertificateCounter, 1)
	defer mm_atomic.AddUint64(&mmNewUnsignedCertificate.afterNewUnsignedCertificateCounter, 1)

	if mmNewUnsignedCertificate.inspectFuncNewUnsignedCertificate != nil {
		mmNewUnsignedCertificate.inspectFuncNewUnsignedCertificate(pKey, role, nodeRef)
	}

	params := &CertificateManagerMockNewUnsignedCertificateParams{pKey, role, nodeRef}

	// Record call args
	mmNewUnsignedCertificate.NewUnsignedCertificateMock.mutex.Lock()
	mmNewUnsignedCertificate.NewUnsignedCertificateMock.callArgs = append(mmNewUnsignedCertificate.NewUnsignedCertificateMock.callArgs, params)
	mmNewUnsignedCertificate.NewUnsignedCertificateMock.mutex.Unlock()

	for _, e := range mmNewUnsignedCertificate.NewUnsignedCertificateMock.expectations {
		if minimock.Equal(e.params, params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.c1, e.results.err
		}
	}

	if mmNewUnsignedCertificate.NewUnsignedCertificateMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmNewUnsignedCertificate.NewUnsignedCertificateMock.defaultExpectation.Counter, 1)
		want := mmNewUnsignedCertificate.NewUnsignedCertificateMock.defaultExpectation.params
		got := CertificateManagerMockNewUnsignedCertificateParams{pKey, role, nodeRef}
		if want != nil && !minimock.Equal(*want, got) {
			mmNewUnsignedCertificate.t.Errorf("CertificateManagerMock.NewUnsignedCertificate got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := mmNewUnsignedCertificate.NewUnsignedCertificateMock.defaultExpectation.results
		if results == nil {
			mmNewUnsignedCertificate.t.Fatal("No results are set for the CertificateManagerMock.NewUnsignedCertificate")
		}
		return (*results).c1, (*results).err
	}
	if mmNewUnsignedCertificate.funcNewUnsignedCertificate != nil {
		return mmNewUnsignedCertificate.funcNewUnsignedCertificate(pKey, role, nodeRef)
	}
	mmNewUnsignedCertificate.t.Fatalf("Unexpected call to CertificateManagerMock.NewUnsignedCertificate. %v %v %v", pKey, role, nodeRef)
	return
}

// NewUnsignedCertificateAfterCounter returns a count of finished CertificateManagerMock.NewUnsignedCertificate invocations
func (mmNewUnsignedCertificate *CertificateManagerMock) NewUnsignedCertificateAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmNewUnsignedCertificate.afterNewUnsignedCertificateCounter)
}

// NewUnsignedCertificateBeforeCounter returns a count of CertificateManagerMock.NewUnsignedCertificate invocations
func (mmNewUnsignedCertificate *CertificateManagerMock) NewUnsignedCertificateBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmNewUnsignedCertificate.beforeNewUnsignedCertificateCounter)
}

// Calls returns a list of arguments used in each call to CertificateManagerMock.NewUnsignedCertificate.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmNewUnsignedCertificate *mCertificateManagerMockNewUnsignedCertificate) Calls() []*CertificateManagerMockNewUnsignedCertificateParams {
	mmNewUnsignedCertificate.mutex.RLock()

	argCopy := make([]*CertificateManagerMockNewUnsignedCertificateParams, len(mmNewUnsignedCertificate.callArgs))
	copy(argCopy, mmNewUnsignedCertificate.callArgs)

	mmNewUnsignedCertificate.mutex.RUnlock()

	return argCopy
}

// MinimockNewUnsignedCertificateDone returns true if the count of the NewUnsignedCertificate invocations corresponds
// the number of defined expectations
func (m *CertificateManagerMock) MinimockNewUnsignedCertificateDone() bool {
	for _, e := range m.NewUnsignedCertificateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.NewUnsignedCertificateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterNewUnsignedCertificateCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcNewUnsignedCertificate != nil && mm_atomic.LoadUint64(&m.afterNewUnsignedCertificateCounter) < 1 {
		return false
	}
	return true
}

// MinimockNewUnsignedCertificateInspect logs each unmet expectation
func (m *CertificateManagerMock) MinimockNewUnsignedCertificateInspect() {
	for _, e := range m.NewUnsignedCertificateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CertificateManagerMock.NewUnsignedCertificate with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.NewUnsignedCertificateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterNewUnsignedCertificateCounter) < 1 {
		if m.NewUnsignedCertificateMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to CertificateManagerMock.NewUnsignedCertificate")
		} else {
			m.t.Errorf("Expected call to CertificateManagerMock.NewUnsignedCertificate with params: %#v", *m.NewUnsignedCertificateMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcNewUnsignedCertificate != nil && mm_atomic.LoadUint64(&m.afterNewUnsignedCertificateCounter) < 1 {
		m.t.Error("Expected call to CertificateManagerMock.NewUnsignedCertificate")
	}
}

type mCertificateManagerMockVerifyAuthorizationCertificate struct {
	mock               *CertificateManagerMock
	defaultExpectation *CertificateManagerMockVerifyAuthorizationCertificateExpectation
	expectations       []*CertificateManagerMockVerifyAuthorizationCertificateExpectation

	callArgs []*CertificateManagerMockVerifyAuthorizationCertificateParams
	mutex    sync.RWMutex
}

// CertificateManagerMockVerifyAuthorizationCertificateExpectation specifies expectation struct of the CertificateManager.VerifyAuthorizationCertificate
type CertificateManagerMockVerifyAuthorizationCertificateExpectation struct {
	mock    *CertificateManagerMock
	params  *CertificateManagerMockVerifyAuthorizationCertificateParams
	results *CertificateManagerMockVerifyAuthorizationCertificateResults
	Counter uint64
}

// CertificateManagerMockVerifyAuthorizationCertificateParams contains parameters of the CertificateManager.VerifyAuthorizationCertificate
type CertificateManagerMockVerifyAuthorizationCertificateParams struct {
	authCert mm_insolar.AuthorizationCertificate
}

// CertificateManagerMockVerifyAuthorizationCertificateResults contains results of the CertificateManager.VerifyAuthorizationCertificate
type CertificateManagerMockVerifyAuthorizationCertificateResults struct {
	b1  bool
	err error
}

// Expect sets up expected params for CertificateManager.VerifyAuthorizationCertificate
func (mmVerifyAuthorizationCertificate *mCertificateManagerMockVerifyAuthorizationCertificate) Expect(authCert mm_insolar.AuthorizationCertificate) *mCertificateManagerMockVerifyAuthorizationCertificate {
	if mmVerifyAuthorizationCertificate.mock.funcVerifyAuthorizationCertificate != nil {
		mmVerifyAuthorizationCertificate.mock.t.Fatalf("CertificateManagerMock.VerifyAuthorizationCertificate mock is already set by Set")
	}

	if mmVerifyAuthorizationCertificate.defaultExpectation == nil {
		mmVerifyAuthorizationCertificate.defaultExpectation = &CertificateManagerMockVerifyAuthorizationCertificateExpectation{}
	}

	mmVerifyAuthorizationCertificate.defaultExpectation.params = &CertificateManagerMockVerifyAuthorizationCertificateParams{authCert}
	for _, e := range mmVerifyAuthorizationCertificate.expectations {
		if minimock.Equal(e.params, mmVerifyAuthorizationCertificate.defaultExpectation.params) {
			mmVerifyAuthorizationCertificate.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmVerifyAuthorizationCertificate.defaultExpectation.params)
		}
	}

	return mmVerifyAuthorizationCertificate
}

// Inspect accepts an inspector function that has same arguments as the CertificateManager.VerifyAuthorizationCertificate
func (mmVerifyAuthorizationCertificate *mCertificateManagerMockVerifyAuthorizationCertificate) Inspect(f func(authCert mm_insolar.AuthorizationCertificate)) *mCertificateManagerMockVerifyAuthorizationCertificate {
	if mmVerifyAuthorizationCertificate.mock.inspectFuncVerifyAuthorizationCertificate != nil {
		mmVerifyAuthorizationCertificate.mock.t.Fatalf("Inspect function is already set for CertificateManagerMock.VerifyAuthorizationCertificate")
	}

	mmVerifyAuthorizationCertificate.mock.inspectFuncVerifyAuthorizationCertificate = f

	return mmVerifyAuthorizationCertificate
}

// Return sets up results that will be returned by CertificateManager.VerifyAuthorizationCertificate
func (mmVerifyAuthorizationCertificate *mCertificateManagerMockVerifyAuthorizationCertificate) Return(b1 bool, err error) *CertificateManagerMock {
	if mmVerifyAuthorizationCertificate.mock.funcVerifyAuthorizationCertificate != nil {
		mmVerifyAuthorizationCertificate.mock.t.Fatalf("CertificateManagerMock.VerifyAuthorizationCertificate mock is already set by Set")
	}

	if mmVerifyAuthorizationCertificate.defaultExpectation == nil {
		mmVerifyAuthorizationCertificate.defaultExpectation = &CertificateManagerMockVerifyAuthorizationCertificateExpectation{mock: mmVerifyAuthorizationCertificate.mock}
	}
	mmVerifyAuthorizationCertificate.defaultExpectation.results = &CertificateManagerMockVerifyAuthorizationCertificateResults{b1, err}
	return mmVerifyAuthorizationCertificate.mock
}

//Set uses given function f to mock the CertificateManager.VerifyAuthorizationCertificate method
func (mmVerifyAuthorizationCertificate *mCertificateManagerMockVerifyAuthorizationCertificate) Set(f func(authCert mm_insolar.AuthorizationCertificate) (b1 bool, err error)) *CertificateManagerMock {
	if mmVerifyAuthorizationCertificate.defaultExpectation != nil {
		mmVerifyAuthorizationCertificate.mock.t.Fatalf("Default expectation is already set for the CertificateManager.VerifyAuthorizationCertificate method")
	}

	if len(mmVerifyAuthorizationCertificate.expectations) > 0 {
		mmVerifyAuthorizationCertificate.mock.t.Fatalf("Some expectations are already set for the CertificateManager.VerifyAuthorizationCertificate method")
	}

	mmVerifyAuthorizationCertificate.mock.funcVerifyAuthorizationCertificate = f
	return mmVerifyAuthorizationCertificate.mock
}

// When sets expectation for the CertificateManager.VerifyAuthorizationCertificate which will trigger the result defined by the following
// Then helper
func (mmVerifyAuthorizationCertificate *mCertificateManagerMockVerifyAuthorizationCertificate) When(authCert mm_insolar.AuthorizationCertificate) *CertificateManagerMockVerifyAuthorizationCertificateExpectation {
	if mmVerifyAuthorizationCertificate.mock.funcVerifyAuthorizationCertificate != nil {
		mmVerifyAuthorizationCertificate.mock.t.Fatalf("CertificateManagerMock.VerifyAuthorizationCertificate mock is already set by Set")
	}

	expectation := &CertificateManagerMockVerifyAuthorizationCertificateExpectation{
		mock:   mmVerifyAuthorizationCertificate.mock,
		params: &CertificateManagerMockVerifyAuthorizationCertificateParams{authCert},
	}
	mmVerifyAuthorizationCertificate.expectations = append(mmVerifyAuthorizationCertificate.expectations, expectation)
	return expectation
}

// Then sets up CertificateManager.VerifyAuthorizationCertificate return parameters for the expectation previously defined by the When method
func (e *CertificateManagerMockVerifyAuthorizationCertificateExpectation) Then(b1 bool, err error) *CertificateManagerMock {
	e.results = &CertificateManagerMockVerifyAuthorizationCertificateResults{b1, err}
	return e.mock
}

// VerifyAuthorizationCertificate implements insolar.CertificateManager
func (mmVerifyAuthorizationCertificate *CertificateManagerMock) VerifyAuthorizationCertificate(authCert mm_insolar.AuthorizationCertificate) (b1 bool, err error) {
	mm_atomic.AddUint64(&mmVerifyAuthorizationCertificate.beforeVerifyAuthorizationCertificateCounter, 1)
	defer mm_atomic.AddUint64(&mmVerifyAuthorizationCertificate.afterVerifyAuthorizationCertificateCounter, 1)

	if mmVerifyAuthorizationCertificate.inspectFuncVerifyAuthorizationCertificate != nil {
		mmVerifyAuthorizationCertificate.inspectFuncVerifyAuthorizationCertificate(authCert)
	}

	params := &CertificateManagerMockVerifyAuthorizationCertificateParams{authCert}

	// Record call args
	mmVerifyAuthorizationCertificate.VerifyAuthorizationCertificateMock.mutex.Lock()
	mmVerifyAuthorizationCertificate.VerifyAuthorizationCertificateMock.callArgs = append(mmVerifyAuthorizationCertificate.VerifyAuthorizationCertificateMock.callArgs, params)
	mmVerifyAuthorizationCertificate.VerifyAuthorizationCertificateMock.mutex.Unlock()

	for _, e := range mmVerifyAuthorizationCertificate.VerifyAuthorizationCertificateMock.expectations {
		if minimock.Equal(e.params, params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.b1, e.results.err
		}
	}

	if mmVerifyAuthorizationCertificate.VerifyAuthorizationCertificateMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmVerifyAuthorizationCertificate.VerifyAuthorizationCertificateMock.defaultExpectation.Counter, 1)
		want := mmVerifyAuthorizationCertificate.VerifyAuthorizationCertificateMock.defaultExpectation.params
		got := CertificateManagerMockVerifyAuthorizationCertificateParams{authCert}
		if want != nil && !minimock.Equal(*want, got) {
			mmVerifyAuthorizationCertificate.t.Errorf("CertificateManagerMock.VerifyAuthorizationCertificate got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := mmVerifyAuthorizationCertificate.VerifyAuthorizationCertificateMock.defaultExpectation.results
		if results == nil {
			mmVerifyAuthorizationCertificate.t.Fatal("No results are set for the CertificateManagerMock.VerifyAuthorizationCertificate")
		}
		return (*results).b1, (*results).err
	}
	if mmVerifyAuthorizationCertificate.funcVerifyAuthorizationCertificate != nil {
		return mmVerifyAuthorizationCertificate.funcVerifyAuthorizationCertificate(authCert)
	}
	mmVerifyAuthorizationCertificate.t.Fatalf("Unexpected call to CertificateManagerMock.VerifyAuthorizationCertificate. %v", authCert)
	return
}

// VerifyAuthorizationCertificateAfterCounter returns a count of finished CertificateManagerMock.VerifyAuthorizationCertificate invocations
func (mmVerifyAuthorizationCertificate *CertificateManagerMock) VerifyAuthorizationCertificateAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmVerifyAuthorizationCertificate.afterVerifyAuthorizationCertificateCounter)
}

// VerifyAuthorizationCertificateBeforeCounter returns a count of CertificateManagerMock.VerifyAuthorizationCertificate invocations
func (mmVerifyAuthorizationCertificate *CertificateManagerMock) VerifyAuthorizationCertificateBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmVerifyAuthorizationCertificate.beforeVerifyAuthorizationCertificateCounter)
}

// Calls returns a list of arguments used in each call to CertificateManagerMock.VerifyAuthorizationCertificate.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmVerifyAuthorizationCertificate *mCertificateManagerMockVerifyAuthorizationCertificate) Calls() []*CertificateManagerMockVerifyAuthorizationCertificateParams {
	mmVerifyAuthorizationCertificate.mutex.RLock()

	argCopy := make([]*CertificateManagerMockVerifyAuthorizationCertificateParams, len(mmVerifyAuthorizationCertificate.callArgs))
	copy(argCopy, mmVerifyAuthorizationCertificate.callArgs)

	mmVerifyAuthorizationCertificate.mutex.RUnlock()

	return argCopy
}

// MinimockVerifyAuthorizationCertificateDone returns true if the count of the VerifyAuthorizationCertificate invocations corresponds
// the number of defined expectations
func (m *CertificateManagerMock) MinimockVerifyAuthorizationCertificateDone() bool {
	for _, e := range m.VerifyAuthorizationCertificateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.VerifyAuthorizationCertificateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterVerifyAuthorizationCertificateCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcVerifyAuthorizationCertificate != nil && mm_atomic.LoadUint64(&m.afterVerifyAuthorizationCertificateCounter) < 1 {
		return false
	}
	return true
}

// MinimockVerifyAuthorizationCertificateInspect logs each unmet expectation
func (m *CertificateManagerMock) MinimockVerifyAuthorizationCertificateInspect() {
	for _, e := range m.VerifyAuthorizationCertificateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CertificateManagerMock.VerifyAuthorizationCertificate with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.VerifyAuthorizationCertificateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterVerifyAuthorizationCertificateCounter) < 1 {
		if m.VerifyAuthorizationCertificateMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to CertificateManagerMock.VerifyAuthorizationCertificate")
		} else {
			m.t.Errorf("Expected call to CertificateManagerMock.VerifyAuthorizationCertificate with params: %#v", *m.VerifyAuthorizationCertificateMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcVerifyAuthorizationCertificate != nil && mm_atomic.LoadUint64(&m.afterVerifyAuthorizationCertificateCounter) < 1 {
		m.t.Error("Expected call to CertificateManagerMock.VerifyAuthorizationCertificate")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *CertificateManagerMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockGetCertificateInspect()

		m.MinimockNewUnsignedCertificateInspect()

		m.MinimockVerifyAuthorizationCertificateInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *CertificateManagerMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *CertificateManagerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetCertificateDone() &&
		m.MinimockNewUnsignedCertificateDone() &&
		m.MinimockVerifyAuthorizationCertificateDone()
}
