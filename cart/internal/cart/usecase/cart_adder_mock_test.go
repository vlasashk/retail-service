// Code generated by http://github.com/gojuno/minimock (v3.3.12). DO NOT EDIT.

package usecase_test

//go:generate minimock -i route256/cart/internal/cart/usecase.cartAdder -o cart_adder_mock_test.go -n CartAdderMock -p usecase_test

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// CartAdderMock implements cartAdder
type CartAdderMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcAddItem          func(ctx context.Context, userID int64, skuID int64, count uint16) (err error)
	inspectFuncAddItem   func(ctx context.Context, userID int64, skuID int64, count uint16)
	afterAddItemCounter  uint64
	beforeAddItemCounter uint64
	AddItemMock          mCartAdderMockAddItem
}

// NewCartAdderMock returns a mock for cartAdder
func NewCartAdderMock(t minimock.Tester) *CartAdderMock {
	m := &CartAdderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddItemMock = mCartAdderMockAddItem{mock: m}
	m.AddItemMock.callArgs = []*CartAdderMockAddItemParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mCartAdderMockAddItem struct {
	optional           bool
	mock               *CartAdderMock
	defaultExpectation *CartAdderMockAddItemExpectation
	expectations       []*CartAdderMockAddItemExpectation

	callArgs []*CartAdderMockAddItemParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// CartAdderMockAddItemExpectation specifies expectation struct of the cartAdder.AddItem
type CartAdderMockAddItemExpectation struct {
	mock      *CartAdderMock
	params    *CartAdderMockAddItemParams
	paramPtrs *CartAdderMockAddItemParamPtrs
	results   *CartAdderMockAddItemResults
	Counter   uint64
}

// CartAdderMockAddItemParams contains parameters of the cartAdder.AddItem
type CartAdderMockAddItemParams struct {
	ctx    context.Context
	userID int64
	skuID  int64
	count  uint16
}

// CartAdderMockAddItemParamPtrs contains pointers to parameters of the cartAdder.AddItem
type CartAdderMockAddItemParamPtrs struct {
	ctx    *context.Context
	userID *int64
	skuID  *int64
	count  *uint16
}

// CartAdderMockAddItemResults contains results of the cartAdder.AddItem
type CartAdderMockAddItemResults struct {
	err error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmAddItem *mCartAdderMockAddItem) Optional() *mCartAdderMockAddItem {
	mmAddItem.optional = true
	return mmAddItem
}

// Expect sets up expected params for cartAdder.AddItem
func (mmAddItem *mCartAdderMockAddItem) Expect(ctx context.Context, userID int64, skuID int64, count uint16) *mCartAdderMockAddItem {
	if mmAddItem.mock.funcAddItem != nil {
		mmAddItem.mock.t.Fatalf("CartAdderMock.AddItem mock is already set by Set")
	}

	if mmAddItem.defaultExpectation == nil {
		mmAddItem.defaultExpectation = &CartAdderMockAddItemExpectation{}
	}

	if mmAddItem.defaultExpectation.paramPtrs != nil {
		mmAddItem.mock.t.Fatalf("CartAdderMock.AddItem mock is already set by ExpectParams functions")
	}

	mmAddItem.defaultExpectation.params = &CartAdderMockAddItemParams{ctx, userID, skuID, count}
	for _, e := range mmAddItem.expectations {
		if minimock.Equal(e.params, mmAddItem.defaultExpectation.params) {
			mmAddItem.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmAddItem.defaultExpectation.params)
		}
	}

	return mmAddItem
}

// ExpectCtxParam1 sets up expected param ctx for cartAdder.AddItem
func (mmAddItem *mCartAdderMockAddItem) ExpectCtxParam1(ctx context.Context) *mCartAdderMockAddItem {
	if mmAddItem.mock.funcAddItem != nil {
		mmAddItem.mock.t.Fatalf("CartAdderMock.AddItem mock is already set by Set")
	}

	if mmAddItem.defaultExpectation == nil {
		mmAddItem.defaultExpectation = &CartAdderMockAddItemExpectation{}
	}

	if mmAddItem.defaultExpectation.params != nil {
		mmAddItem.mock.t.Fatalf("CartAdderMock.AddItem mock is already set by Expect")
	}

	if mmAddItem.defaultExpectation.paramPtrs == nil {
		mmAddItem.defaultExpectation.paramPtrs = &CartAdderMockAddItemParamPtrs{}
	}
	mmAddItem.defaultExpectation.paramPtrs.ctx = &ctx

	return mmAddItem
}

// ExpectUserIDParam2 sets up expected param userID for cartAdder.AddItem
func (mmAddItem *mCartAdderMockAddItem) ExpectUserIDParam2(userID int64) *mCartAdderMockAddItem {
	if mmAddItem.mock.funcAddItem != nil {
		mmAddItem.mock.t.Fatalf("CartAdderMock.AddItem mock is already set by Set")
	}

	if mmAddItem.defaultExpectation == nil {
		mmAddItem.defaultExpectation = &CartAdderMockAddItemExpectation{}
	}

	if mmAddItem.defaultExpectation.params != nil {
		mmAddItem.mock.t.Fatalf("CartAdderMock.AddItem mock is already set by Expect")
	}

	if mmAddItem.defaultExpectation.paramPtrs == nil {
		mmAddItem.defaultExpectation.paramPtrs = &CartAdderMockAddItemParamPtrs{}
	}
	mmAddItem.defaultExpectation.paramPtrs.userID = &userID

	return mmAddItem
}

// ExpectSkuIDParam3 sets up expected param skuID for cartAdder.AddItem
func (mmAddItem *mCartAdderMockAddItem) ExpectSkuIDParam3(skuID int64) *mCartAdderMockAddItem {
	if mmAddItem.mock.funcAddItem != nil {
		mmAddItem.mock.t.Fatalf("CartAdderMock.AddItem mock is already set by Set")
	}

	if mmAddItem.defaultExpectation == nil {
		mmAddItem.defaultExpectation = &CartAdderMockAddItemExpectation{}
	}

	if mmAddItem.defaultExpectation.params != nil {
		mmAddItem.mock.t.Fatalf("CartAdderMock.AddItem mock is already set by Expect")
	}

	if mmAddItem.defaultExpectation.paramPtrs == nil {
		mmAddItem.defaultExpectation.paramPtrs = &CartAdderMockAddItemParamPtrs{}
	}
	mmAddItem.defaultExpectation.paramPtrs.skuID = &skuID

	return mmAddItem
}

// ExpectCountParam4 sets up expected param count for cartAdder.AddItem
func (mmAddItem *mCartAdderMockAddItem) ExpectCountParam4(count uint16) *mCartAdderMockAddItem {
	if mmAddItem.mock.funcAddItem != nil {
		mmAddItem.mock.t.Fatalf("CartAdderMock.AddItem mock is already set by Set")
	}

	if mmAddItem.defaultExpectation == nil {
		mmAddItem.defaultExpectation = &CartAdderMockAddItemExpectation{}
	}

	if mmAddItem.defaultExpectation.params != nil {
		mmAddItem.mock.t.Fatalf("CartAdderMock.AddItem mock is already set by Expect")
	}

	if mmAddItem.defaultExpectation.paramPtrs == nil {
		mmAddItem.defaultExpectation.paramPtrs = &CartAdderMockAddItemParamPtrs{}
	}
	mmAddItem.defaultExpectation.paramPtrs.count = &count

	return mmAddItem
}

// Inspect accepts an inspector function that has same arguments as the cartAdder.AddItem
func (mmAddItem *mCartAdderMockAddItem) Inspect(f func(ctx context.Context, userID int64, skuID int64, count uint16)) *mCartAdderMockAddItem {
	if mmAddItem.mock.inspectFuncAddItem != nil {
		mmAddItem.mock.t.Fatalf("Inspect function is already set for CartAdderMock.AddItem")
	}

	mmAddItem.mock.inspectFuncAddItem = f

	return mmAddItem
}

// Return sets up results that will be returned by cartAdder.AddItem
func (mmAddItem *mCartAdderMockAddItem) Return(err error) *CartAdderMock {
	if mmAddItem.mock.funcAddItem != nil {
		mmAddItem.mock.t.Fatalf("CartAdderMock.AddItem mock is already set by Set")
	}

	if mmAddItem.defaultExpectation == nil {
		mmAddItem.defaultExpectation = &CartAdderMockAddItemExpectation{mock: mmAddItem.mock}
	}
	mmAddItem.defaultExpectation.results = &CartAdderMockAddItemResults{err}
	return mmAddItem.mock
}

// Set uses given function f to mock the cartAdder.AddItem method
func (mmAddItem *mCartAdderMockAddItem) Set(f func(ctx context.Context, userID int64, skuID int64, count uint16) (err error)) *CartAdderMock {
	if mmAddItem.defaultExpectation != nil {
		mmAddItem.mock.t.Fatalf("Default expectation is already set for the cartAdder.AddItem method")
	}

	if len(mmAddItem.expectations) > 0 {
		mmAddItem.mock.t.Fatalf("Some expectations are already set for the cartAdder.AddItem method")
	}

	mmAddItem.mock.funcAddItem = f
	return mmAddItem.mock
}

// When sets expectation for the cartAdder.AddItem which will trigger the result defined by the following
// Then helper
func (mmAddItem *mCartAdderMockAddItem) When(ctx context.Context, userID int64, skuID int64, count uint16) *CartAdderMockAddItemExpectation {
	if mmAddItem.mock.funcAddItem != nil {
		mmAddItem.mock.t.Fatalf("CartAdderMock.AddItem mock is already set by Set")
	}

	expectation := &CartAdderMockAddItemExpectation{
		mock:   mmAddItem.mock,
		params: &CartAdderMockAddItemParams{ctx, userID, skuID, count},
	}
	mmAddItem.expectations = append(mmAddItem.expectations, expectation)
	return expectation
}

// Then sets up cartAdder.AddItem return parameters for the expectation previously defined by the When method
func (e *CartAdderMockAddItemExpectation) Then(err error) *CartAdderMock {
	e.results = &CartAdderMockAddItemResults{err}
	return e.mock
}

// Times sets number of times cartAdder.AddItem should be invoked
func (mmAddItem *mCartAdderMockAddItem) Times(n uint64) *mCartAdderMockAddItem {
	if n == 0 {
		mmAddItem.mock.t.Fatalf("Times of CartAdderMock.AddItem mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmAddItem.expectedInvocations, n)
	return mmAddItem
}

func (mmAddItem *mCartAdderMockAddItem) invocationsDone() bool {
	if len(mmAddItem.expectations) == 0 && mmAddItem.defaultExpectation == nil && mmAddItem.mock.funcAddItem == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmAddItem.mock.afterAddItemCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmAddItem.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// AddItem implements cartAdder
func (mmAddItem *CartAdderMock) AddItem(ctx context.Context, userID int64, skuID int64, count uint16) (err error) {
	mm_atomic.AddUint64(&mmAddItem.beforeAddItemCounter, 1)
	defer mm_atomic.AddUint64(&mmAddItem.afterAddItemCounter, 1)

	if mmAddItem.inspectFuncAddItem != nil {
		mmAddItem.inspectFuncAddItem(ctx, userID, skuID, count)
	}

	mm_params := CartAdderMockAddItemParams{ctx, userID, skuID, count}

	// Record call args
	mmAddItem.AddItemMock.mutex.Lock()
	mmAddItem.AddItemMock.callArgs = append(mmAddItem.AddItemMock.callArgs, &mm_params)
	mmAddItem.AddItemMock.mutex.Unlock()

	for _, e := range mmAddItem.AddItemMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmAddItem.AddItemMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmAddItem.AddItemMock.defaultExpectation.Counter, 1)
		mm_want := mmAddItem.AddItemMock.defaultExpectation.params
		mm_want_ptrs := mmAddItem.AddItemMock.defaultExpectation.paramPtrs

		mm_got := CartAdderMockAddItemParams{ctx, userID, skuID, count}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmAddItem.t.Errorf("CartAdderMock.AddItem got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.userID != nil && !minimock.Equal(*mm_want_ptrs.userID, mm_got.userID) {
				mmAddItem.t.Errorf("CartAdderMock.AddItem got unexpected parameter userID, want: %#v, got: %#v%s\n", *mm_want_ptrs.userID, mm_got.userID, minimock.Diff(*mm_want_ptrs.userID, mm_got.userID))
			}

			if mm_want_ptrs.skuID != nil && !minimock.Equal(*mm_want_ptrs.skuID, mm_got.skuID) {
				mmAddItem.t.Errorf("CartAdderMock.AddItem got unexpected parameter skuID, want: %#v, got: %#v%s\n", *mm_want_ptrs.skuID, mm_got.skuID, minimock.Diff(*mm_want_ptrs.skuID, mm_got.skuID))
			}

			if mm_want_ptrs.count != nil && !minimock.Equal(*mm_want_ptrs.count, mm_got.count) {
				mmAddItem.t.Errorf("CartAdderMock.AddItem got unexpected parameter count, want: %#v, got: %#v%s\n", *mm_want_ptrs.count, mm_got.count, minimock.Diff(*mm_want_ptrs.count, mm_got.count))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmAddItem.t.Errorf("CartAdderMock.AddItem got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmAddItem.AddItemMock.defaultExpectation.results
		if mm_results == nil {
			mmAddItem.t.Fatal("No results are set for the CartAdderMock.AddItem")
		}
		return (*mm_results).err
	}
	if mmAddItem.funcAddItem != nil {
		return mmAddItem.funcAddItem(ctx, userID, skuID, count)
	}
	mmAddItem.t.Fatalf("Unexpected call to CartAdderMock.AddItem. %v %v %v %v", ctx, userID, skuID, count)
	return
}

// AddItemAfterCounter returns a count of finished CartAdderMock.AddItem invocations
func (mmAddItem *CartAdderMock) AddItemAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmAddItem.afterAddItemCounter)
}

// AddItemBeforeCounter returns a count of CartAdderMock.AddItem invocations
func (mmAddItem *CartAdderMock) AddItemBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmAddItem.beforeAddItemCounter)
}

// Calls returns a list of arguments used in each call to CartAdderMock.AddItem.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmAddItem *mCartAdderMockAddItem) Calls() []*CartAdderMockAddItemParams {
	mmAddItem.mutex.RLock()

	argCopy := make([]*CartAdderMockAddItemParams, len(mmAddItem.callArgs))
	copy(argCopy, mmAddItem.callArgs)

	mmAddItem.mutex.RUnlock()

	return argCopy
}

// MinimockAddItemDone returns true if the count of the AddItem invocations corresponds
// the number of defined expectations
func (m *CartAdderMock) MinimockAddItemDone() bool {
	if m.AddItemMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.AddItemMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.AddItemMock.invocationsDone()
}

// MinimockAddItemInspect logs each unmet expectation
func (m *CartAdderMock) MinimockAddItemInspect() {
	for _, e := range m.AddItemMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CartAdderMock.AddItem with params: %#v", *e.params)
		}
	}

	afterAddItemCounter := mm_atomic.LoadUint64(&m.afterAddItemCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.AddItemMock.defaultExpectation != nil && afterAddItemCounter < 1 {
		if m.AddItemMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to CartAdderMock.AddItem")
		} else {
			m.t.Errorf("Expected call to CartAdderMock.AddItem with params: %#v", *m.AddItemMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAddItem != nil && afterAddItemCounter < 1 {
		m.t.Error("Expected call to CartAdderMock.AddItem")
	}

	if !m.AddItemMock.invocationsDone() && afterAddItemCounter > 0 {
		m.t.Errorf("Expected %d calls to CartAdderMock.AddItem but found %d calls",
			mm_atomic.LoadUint64(&m.AddItemMock.expectedInvocations), afterAddItemCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *CartAdderMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockAddItemInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *CartAdderMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *CartAdderMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockAddItemDone()
}
