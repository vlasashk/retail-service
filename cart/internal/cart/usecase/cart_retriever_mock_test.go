// Code generated by http://github.com/gojuno/minimock (v3.3.12). DO NOT EDIT.

package usecase_test

//go:generate minimock -i route256/cart/internal/cart/usecase.cartRetriever -o cart_retriever_mock_test.go -n CartRetrieverMock -p usecase_test

import (
	"context"
	"route256/cart/internal/cart/models"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// CartRetrieverMock implements cartRetriever
type CartRetrieverMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetItemsByUserID          func(ctx context.Context, userID int64) (ia1 []models.Item, err error)
	inspectFuncGetItemsByUserID   func(ctx context.Context, userID int64)
	afterGetItemsByUserIDCounter  uint64
	beforeGetItemsByUserIDCounter uint64
	GetItemsByUserIDMock          mCartRetrieverMockGetItemsByUserID
}

// NewCartRetrieverMock returns a mock for cartRetriever
func NewCartRetrieverMock(t minimock.Tester) *CartRetrieverMock {
	m := &CartRetrieverMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetItemsByUserIDMock = mCartRetrieverMockGetItemsByUserID{mock: m}
	m.GetItemsByUserIDMock.callArgs = []*CartRetrieverMockGetItemsByUserIDParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mCartRetrieverMockGetItemsByUserID struct {
	optional           bool
	mock               *CartRetrieverMock
	defaultExpectation *CartRetrieverMockGetItemsByUserIDExpectation
	expectations       []*CartRetrieverMockGetItemsByUserIDExpectation

	callArgs []*CartRetrieverMockGetItemsByUserIDParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// CartRetrieverMockGetItemsByUserIDExpectation specifies expectation struct of the cartRetriever.GetItemsByUserID
type CartRetrieverMockGetItemsByUserIDExpectation struct {
	mock      *CartRetrieverMock
	params    *CartRetrieverMockGetItemsByUserIDParams
	paramPtrs *CartRetrieverMockGetItemsByUserIDParamPtrs
	results   *CartRetrieverMockGetItemsByUserIDResults
	Counter   uint64
}

// CartRetrieverMockGetItemsByUserIDParams contains parameters of the cartRetriever.GetItemsByUserID
type CartRetrieverMockGetItemsByUserIDParams struct {
	ctx    context.Context
	userID int64
}

// CartRetrieverMockGetItemsByUserIDParamPtrs contains pointers to parameters of the cartRetriever.GetItemsByUserID
type CartRetrieverMockGetItemsByUserIDParamPtrs struct {
	ctx    *context.Context
	userID *int64
}

// CartRetrieverMockGetItemsByUserIDResults contains results of the cartRetriever.GetItemsByUserID
type CartRetrieverMockGetItemsByUserIDResults struct {
	ia1 []models.Item
	err error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmGetItemsByUserID *mCartRetrieverMockGetItemsByUserID) Optional() *mCartRetrieverMockGetItemsByUserID {
	mmGetItemsByUserID.optional = true
	return mmGetItemsByUserID
}

// Expect sets up expected params for cartRetriever.GetItemsByUserID
func (mmGetItemsByUserID *mCartRetrieverMockGetItemsByUserID) Expect(ctx context.Context, userID int64) *mCartRetrieverMockGetItemsByUserID {
	if mmGetItemsByUserID.mock.funcGetItemsByUserID != nil {
		mmGetItemsByUserID.mock.t.Fatalf("CartRetrieverMock.GetItemsByUserID mock is already set by Set")
	}

	if mmGetItemsByUserID.defaultExpectation == nil {
		mmGetItemsByUserID.defaultExpectation = &CartRetrieverMockGetItemsByUserIDExpectation{}
	}

	if mmGetItemsByUserID.defaultExpectation.paramPtrs != nil {
		mmGetItemsByUserID.mock.t.Fatalf("CartRetrieverMock.GetItemsByUserID mock is already set by ExpectParams functions")
	}

	mmGetItemsByUserID.defaultExpectation.params = &CartRetrieverMockGetItemsByUserIDParams{ctx, userID}
	for _, e := range mmGetItemsByUserID.expectations {
		if minimock.Equal(e.params, mmGetItemsByUserID.defaultExpectation.params) {
			mmGetItemsByUserID.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetItemsByUserID.defaultExpectation.params)
		}
	}

	return mmGetItemsByUserID
}

// ExpectCtxParam1 sets up expected param ctx for cartRetriever.GetItemsByUserID
func (mmGetItemsByUserID *mCartRetrieverMockGetItemsByUserID) ExpectCtxParam1(ctx context.Context) *mCartRetrieverMockGetItemsByUserID {
	if mmGetItemsByUserID.mock.funcGetItemsByUserID != nil {
		mmGetItemsByUserID.mock.t.Fatalf("CartRetrieverMock.GetItemsByUserID mock is already set by Set")
	}

	if mmGetItemsByUserID.defaultExpectation == nil {
		mmGetItemsByUserID.defaultExpectation = &CartRetrieverMockGetItemsByUserIDExpectation{}
	}

	if mmGetItemsByUserID.defaultExpectation.params != nil {
		mmGetItemsByUserID.mock.t.Fatalf("CartRetrieverMock.GetItemsByUserID mock is already set by Expect")
	}

	if mmGetItemsByUserID.defaultExpectation.paramPtrs == nil {
		mmGetItemsByUserID.defaultExpectation.paramPtrs = &CartRetrieverMockGetItemsByUserIDParamPtrs{}
	}
	mmGetItemsByUserID.defaultExpectation.paramPtrs.ctx = &ctx

	return mmGetItemsByUserID
}

// ExpectUserIDParam2 sets up expected param userID for cartRetriever.GetItemsByUserID
func (mmGetItemsByUserID *mCartRetrieverMockGetItemsByUserID) ExpectUserIDParam2(userID int64) *mCartRetrieverMockGetItemsByUserID {
	if mmGetItemsByUserID.mock.funcGetItemsByUserID != nil {
		mmGetItemsByUserID.mock.t.Fatalf("CartRetrieverMock.GetItemsByUserID mock is already set by Set")
	}

	if mmGetItemsByUserID.defaultExpectation == nil {
		mmGetItemsByUserID.defaultExpectation = &CartRetrieverMockGetItemsByUserIDExpectation{}
	}

	if mmGetItemsByUserID.defaultExpectation.params != nil {
		mmGetItemsByUserID.mock.t.Fatalf("CartRetrieverMock.GetItemsByUserID mock is already set by Expect")
	}

	if mmGetItemsByUserID.defaultExpectation.paramPtrs == nil {
		mmGetItemsByUserID.defaultExpectation.paramPtrs = &CartRetrieverMockGetItemsByUserIDParamPtrs{}
	}
	mmGetItemsByUserID.defaultExpectation.paramPtrs.userID = &userID

	return mmGetItemsByUserID
}

// Inspect accepts an inspector function that has same arguments as the cartRetriever.GetItemsByUserID
func (mmGetItemsByUserID *mCartRetrieverMockGetItemsByUserID) Inspect(f func(ctx context.Context, userID int64)) *mCartRetrieverMockGetItemsByUserID {
	if mmGetItemsByUserID.mock.inspectFuncGetItemsByUserID != nil {
		mmGetItemsByUserID.mock.t.Fatalf("Inspect function is already set for CartRetrieverMock.GetItemsByUserID")
	}

	mmGetItemsByUserID.mock.inspectFuncGetItemsByUserID = f

	return mmGetItemsByUserID
}

// Return sets up results that will be returned by cartRetriever.GetItemsByUserID
func (mmGetItemsByUserID *mCartRetrieverMockGetItemsByUserID) Return(ia1 []models.Item, err error) *CartRetrieverMock {
	if mmGetItemsByUserID.mock.funcGetItemsByUserID != nil {
		mmGetItemsByUserID.mock.t.Fatalf("CartRetrieverMock.GetItemsByUserID mock is already set by Set")
	}

	if mmGetItemsByUserID.defaultExpectation == nil {
		mmGetItemsByUserID.defaultExpectation = &CartRetrieverMockGetItemsByUserIDExpectation{mock: mmGetItemsByUserID.mock}
	}
	mmGetItemsByUserID.defaultExpectation.results = &CartRetrieverMockGetItemsByUserIDResults{ia1, err}
	return mmGetItemsByUserID.mock
}

// Set uses given function f to mock the cartRetriever.GetItemsByUserID method
func (mmGetItemsByUserID *mCartRetrieverMockGetItemsByUserID) Set(f func(ctx context.Context, userID int64) (ia1 []models.Item, err error)) *CartRetrieverMock {
	if mmGetItemsByUserID.defaultExpectation != nil {
		mmGetItemsByUserID.mock.t.Fatalf("Default expectation is already set for the cartRetriever.GetItemsByUserID method")
	}

	if len(mmGetItemsByUserID.expectations) > 0 {
		mmGetItemsByUserID.mock.t.Fatalf("Some expectations are already set for the cartRetriever.GetItemsByUserID method")
	}

	mmGetItemsByUserID.mock.funcGetItemsByUserID = f
	return mmGetItemsByUserID.mock
}

// When sets expectation for the cartRetriever.GetItemsByUserID which will trigger the result defined by the following
// Then helper
func (mmGetItemsByUserID *mCartRetrieverMockGetItemsByUserID) When(ctx context.Context, userID int64) *CartRetrieverMockGetItemsByUserIDExpectation {
	if mmGetItemsByUserID.mock.funcGetItemsByUserID != nil {
		mmGetItemsByUserID.mock.t.Fatalf("CartRetrieverMock.GetItemsByUserID mock is already set by Set")
	}

	expectation := &CartRetrieverMockGetItemsByUserIDExpectation{
		mock:   mmGetItemsByUserID.mock,
		params: &CartRetrieverMockGetItemsByUserIDParams{ctx, userID},
	}
	mmGetItemsByUserID.expectations = append(mmGetItemsByUserID.expectations, expectation)
	return expectation
}

// Then sets up cartRetriever.GetItemsByUserID return parameters for the expectation previously defined by the When method
func (e *CartRetrieverMockGetItemsByUserIDExpectation) Then(ia1 []models.Item, err error) *CartRetrieverMock {
	e.results = &CartRetrieverMockGetItemsByUserIDResults{ia1, err}
	return e.mock
}

// Times sets number of times cartRetriever.GetItemsByUserID should be invoked
func (mmGetItemsByUserID *mCartRetrieverMockGetItemsByUserID) Times(n uint64) *mCartRetrieverMockGetItemsByUserID {
	if n == 0 {
		mmGetItemsByUserID.mock.t.Fatalf("Times of CartRetrieverMock.GetItemsByUserID mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmGetItemsByUserID.expectedInvocations, n)
	return mmGetItemsByUserID
}

func (mmGetItemsByUserID *mCartRetrieverMockGetItemsByUserID) invocationsDone() bool {
	if len(mmGetItemsByUserID.expectations) == 0 && mmGetItemsByUserID.defaultExpectation == nil && mmGetItemsByUserID.mock.funcGetItemsByUserID == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmGetItemsByUserID.mock.afterGetItemsByUserIDCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmGetItemsByUserID.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// GetItemsByUserID implements cartRetriever
func (mmGetItemsByUserID *CartRetrieverMock) GetItemsByUserID(ctx context.Context, userID int64) (ia1 []models.Item, err error) {
	mm_atomic.AddUint64(&mmGetItemsByUserID.beforeGetItemsByUserIDCounter, 1)
	defer mm_atomic.AddUint64(&mmGetItemsByUserID.afterGetItemsByUserIDCounter, 1)

	if mmGetItemsByUserID.inspectFuncGetItemsByUserID != nil {
		mmGetItemsByUserID.inspectFuncGetItemsByUserID(ctx, userID)
	}

	mm_params := CartRetrieverMockGetItemsByUserIDParams{ctx, userID}

	// Record call args
	mmGetItemsByUserID.GetItemsByUserIDMock.mutex.Lock()
	mmGetItemsByUserID.GetItemsByUserIDMock.callArgs = append(mmGetItemsByUserID.GetItemsByUserIDMock.callArgs, &mm_params)
	mmGetItemsByUserID.GetItemsByUserIDMock.mutex.Unlock()

	for _, e := range mmGetItemsByUserID.GetItemsByUserIDMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.ia1, e.results.err
		}
	}

	if mmGetItemsByUserID.GetItemsByUserIDMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetItemsByUserID.GetItemsByUserIDMock.defaultExpectation.Counter, 1)
		mm_want := mmGetItemsByUserID.GetItemsByUserIDMock.defaultExpectation.params
		mm_want_ptrs := mmGetItemsByUserID.GetItemsByUserIDMock.defaultExpectation.paramPtrs

		mm_got := CartRetrieverMockGetItemsByUserIDParams{ctx, userID}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmGetItemsByUserID.t.Errorf("CartRetrieverMock.GetItemsByUserID got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.userID != nil && !minimock.Equal(*mm_want_ptrs.userID, mm_got.userID) {
				mmGetItemsByUserID.t.Errorf("CartRetrieverMock.GetItemsByUserID got unexpected parameter userID, want: %#v, got: %#v%s\n", *mm_want_ptrs.userID, mm_got.userID, minimock.Diff(*mm_want_ptrs.userID, mm_got.userID))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetItemsByUserID.t.Errorf("CartRetrieverMock.GetItemsByUserID got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetItemsByUserID.GetItemsByUserIDMock.defaultExpectation.results
		if mm_results == nil {
			mmGetItemsByUserID.t.Fatal("No results are set for the CartRetrieverMock.GetItemsByUserID")
		}
		return (*mm_results).ia1, (*mm_results).err
	}
	if mmGetItemsByUserID.funcGetItemsByUserID != nil {
		return mmGetItemsByUserID.funcGetItemsByUserID(ctx, userID)
	}
	mmGetItemsByUserID.t.Fatalf("Unexpected call to CartRetrieverMock.GetItemsByUserID. %v %v", ctx, userID)
	return
}

// GetItemsByUserIDAfterCounter returns a count of finished CartRetrieverMock.GetItemsByUserID invocations
func (mmGetItemsByUserID *CartRetrieverMock) GetItemsByUserIDAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetItemsByUserID.afterGetItemsByUserIDCounter)
}

// GetItemsByUserIDBeforeCounter returns a count of CartRetrieverMock.GetItemsByUserID invocations
func (mmGetItemsByUserID *CartRetrieverMock) GetItemsByUserIDBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetItemsByUserID.beforeGetItemsByUserIDCounter)
}

// Calls returns a list of arguments used in each call to CartRetrieverMock.GetItemsByUserID.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetItemsByUserID *mCartRetrieverMockGetItemsByUserID) Calls() []*CartRetrieverMockGetItemsByUserIDParams {
	mmGetItemsByUserID.mutex.RLock()

	argCopy := make([]*CartRetrieverMockGetItemsByUserIDParams, len(mmGetItemsByUserID.callArgs))
	copy(argCopy, mmGetItemsByUserID.callArgs)

	mmGetItemsByUserID.mutex.RUnlock()

	return argCopy
}

// MinimockGetItemsByUserIDDone returns true if the count of the GetItemsByUserID invocations corresponds
// the number of defined expectations
func (m *CartRetrieverMock) MinimockGetItemsByUserIDDone() bool {
	if m.GetItemsByUserIDMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.GetItemsByUserIDMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.GetItemsByUserIDMock.invocationsDone()
}

// MinimockGetItemsByUserIDInspect logs each unmet expectation
func (m *CartRetrieverMock) MinimockGetItemsByUserIDInspect() {
	for _, e := range m.GetItemsByUserIDMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CartRetrieverMock.GetItemsByUserID with params: %#v", *e.params)
		}
	}

	afterGetItemsByUserIDCounter := mm_atomic.LoadUint64(&m.afterGetItemsByUserIDCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.GetItemsByUserIDMock.defaultExpectation != nil && afterGetItemsByUserIDCounter < 1 {
		if m.GetItemsByUserIDMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to CartRetrieverMock.GetItemsByUserID")
		} else {
			m.t.Errorf("Expected call to CartRetrieverMock.GetItemsByUserID with params: %#v", *m.GetItemsByUserIDMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetItemsByUserID != nil && afterGetItemsByUserIDCounter < 1 {
		m.t.Error("Expected call to CartRetrieverMock.GetItemsByUserID")
	}

	if !m.GetItemsByUserIDMock.invocationsDone() && afterGetItemsByUserIDCounter > 0 {
		m.t.Errorf("Expected %d calls to CartRetrieverMock.GetItemsByUserID but found %d calls",
			mm_atomic.LoadUint64(&m.GetItemsByUserIDMock.expectedInvocations), afterGetItemsByUserIDCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *CartRetrieverMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetItemsByUserIDInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *CartRetrieverMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *CartRetrieverMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetItemsByUserIDDone()
}
