// Code generated by http://github.com/gojuno/minimock (v3.3.11). DO NOT EDIT.

package usecase_test

//go:generate minimock -i route256/cart/internal/cart/usecase.itemRemover -o item_remover_mock_test.go -n ItemRemoverMock -p usecase_test

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// ItemRemoverMock implements itemRemover
type ItemRemoverMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcDeleteItem          func(ctx context.Context, userID int64, skuID int64) (err error)
	inspectFuncDeleteItem   func(ctx context.Context, userID int64, skuID int64)
	afterDeleteItemCounter  uint64
	beforeDeleteItemCounter uint64
	DeleteItemMock          mItemRemoverMockDeleteItem
}

// NewItemRemoverMock returns a mock for itemRemover
func NewItemRemoverMock(t minimock.Tester) *ItemRemoverMock {
	m := &ItemRemoverMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeleteItemMock = mItemRemoverMockDeleteItem{mock: m}
	m.DeleteItemMock.callArgs = []*ItemRemoverMockDeleteItemParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mItemRemoverMockDeleteItem struct {
	optional           bool
	mock               *ItemRemoverMock
	defaultExpectation *ItemRemoverMockDeleteItemExpectation
	expectations       []*ItemRemoverMockDeleteItemExpectation

	callArgs []*ItemRemoverMockDeleteItemParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// ItemRemoverMockDeleteItemExpectation specifies expectation struct of the itemRemover.DeleteItem
type ItemRemoverMockDeleteItemExpectation struct {
	mock      *ItemRemoverMock
	params    *ItemRemoverMockDeleteItemParams
	paramPtrs *ItemRemoverMockDeleteItemParamPtrs
	results   *ItemRemoverMockDeleteItemResults
	Counter   uint64
}

// ItemRemoverMockDeleteItemParams contains parameters of the itemRemover.DeleteItem
type ItemRemoverMockDeleteItemParams struct {
	ctx    context.Context
	userID int64
	skuID  int64
}

// ItemRemoverMockDeleteItemParamPtrs contains pointers to parameters of the itemRemover.DeleteItem
type ItemRemoverMockDeleteItemParamPtrs struct {
	ctx    *context.Context
	userID *int64
	skuID  *int64
}

// ItemRemoverMockDeleteItemResults contains results of the itemRemover.DeleteItem
type ItemRemoverMockDeleteItemResults struct {
	err error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option by default unless you really need it, as it helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmDeleteItem *mItemRemoverMockDeleteItem) Optional() *mItemRemoverMockDeleteItem {
	mmDeleteItem.optional = true
	return mmDeleteItem
}

// Expect sets up expected params for itemRemover.DeleteItem
func (mmDeleteItem *mItemRemoverMockDeleteItem) Expect(ctx context.Context, userID int64, skuID int64) *mItemRemoverMockDeleteItem {
	if mmDeleteItem.mock.funcDeleteItem != nil {
		mmDeleteItem.mock.t.Fatalf("ItemRemoverMock.DeleteItem mock is already set by Set")
	}

	if mmDeleteItem.defaultExpectation == nil {
		mmDeleteItem.defaultExpectation = &ItemRemoverMockDeleteItemExpectation{}
	}

	if mmDeleteItem.defaultExpectation.paramPtrs != nil {
		mmDeleteItem.mock.t.Fatalf("ItemRemoverMock.DeleteItem mock is already set by ExpectParams functions")
	}

	mmDeleteItem.defaultExpectation.params = &ItemRemoverMockDeleteItemParams{ctx, userID, skuID}
	for _, e := range mmDeleteItem.expectations {
		if minimock.Equal(e.params, mmDeleteItem.defaultExpectation.params) {
			mmDeleteItem.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmDeleteItem.defaultExpectation.params)
		}
	}

	return mmDeleteItem
}

// ExpectCtxParam1 sets up expected param ctx for itemRemover.DeleteItem
func (mmDeleteItem *mItemRemoverMockDeleteItem) ExpectCtxParam1(ctx context.Context) *mItemRemoverMockDeleteItem {
	if mmDeleteItem.mock.funcDeleteItem != nil {
		mmDeleteItem.mock.t.Fatalf("ItemRemoverMock.DeleteItem mock is already set by Set")
	}

	if mmDeleteItem.defaultExpectation == nil {
		mmDeleteItem.defaultExpectation = &ItemRemoverMockDeleteItemExpectation{}
	}

	if mmDeleteItem.defaultExpectation.params != nil {
		mmDeleteItem.mock.t.Fatalf("ItemRemoverMock.DeleteItem mock is already set by Expect")
	}

	if mmDeleteItem.defaultExpectation.paramPtrs == nil {
		mmDeleteItem.defaultExpectation.paramPtrs = &ItemRemoverMockDeleteItemParamPtrs{}
	}
	mmDeleteItem.defaultExpectation.paramPtrs.ctx = &ctx

	return mmDeleteItem
}

// ExpectUserIDParam2 sets up expected param userID for itemRemover.DeleteItem
func (mmDeleteItem *mItemRemoverMockDeleteItem) ExpectUserIDParam2(userID int64) *mItemRemoverMockDeleteItem {
	if mmDeleteItem.mock.funcDeleteItem != nil {
		mmDeleteItem.mock.t.Fatalf("ItemRemoverMock.DeleteItem mock is already set by Set")
	}

	if mmDeleteItem.defaultExpectation == nil {
		mmDeleteItem.defaultExpectation = &ItemRemoverMockDeleteItemExpectation{}
	}

	if mmDeleteItem.defaultExpectation.params != nil {
		mmDeleteItem.mock.t.Fatalf("ItemRemoverMock.DeleteItem mock is already set by Expect")
	}

	if mmDeleteItem.defaultExpectation.paramPtrs == nil {
		mmDeleteItem.defaultExpectation.paramPtrs = &ItemRemoverMockDeleteItemParamPtrs{}
	}
	mmDeleteItem.defaultExpectation.paramPtrs.userID = &userID

	return mmDeleteItem
}

// ExpectSkuIDParam3 sets up expected param skuID for itemRemover.DeleteItem
func (mmDeleteItem *mItemRemoverMockDeleteItem) ExpectSkuIDParam3(skuID int64) *mItemRemoverMockDeleteItem {
	if mmDeleteItem.mock.funcDeleteItem != nil {
		mmDeleteItem.mock.t.Fatalf("ItemRemoverMock.DeleteItem mock is already set by Set")
	}

	if mmDeleteItem.defaultExpectation == nil {
		mmDeleteItem.defaultExpectation = &ItemRemoverMockDeleteItemExpectation{}
	}

	if mmDeleteItem.defaultExpectation.params != nil {
		mmDeleteItem.mock.t.Fatalf("ItemRemoverMock.DeleteItem mock is already set by Expect")
	}

	if mmDeleteItem.defaultExpectation.paramPtrs == nil {
		mmDeleteItem.defaultExpectation.paramPtrs = &ItemRemoverMockDeleteItemParamPtrs{}
	}
	mmDeleteItem.defaultExpectation.paramPtrs.skuID = &skuID

	return mmDeleteItem
}

// Inspect accepts an inspector function that has same arguments as the itemRemover.DeleteItem
func (mmDeleteItem *mItemRemoverMockDeleteItem) Inspect(f func(ctx context.Context, userID int64, skuID int64)) *mItemRemoverMockDeleteItem {
	if mmDeleteItem.mock.inspectFuncDeleteItem != nil {
		mmDeleteItem.mock.t.Fatalf("Inspect function is already set for ItemRemoverMock.DeleteItem")
	}

	mmDeleteItem.mock.inspectFuncDeleteItem = f

	return mmDeleteItem
}

// Return sets up results that will be returned by itemRemover.DeleteItem
func (mmDeleteItem *mItemRemoverMockDeleteItem) Return(err error) *ItemRemoverMock {
	if mmDeleteItem.mock.funcDeleteItem != nil {
		mmDeleteItem.mock.t.Fatalf("ItemRemoverMock.DeleteItem mock is already set by Set")
	}

	if mmDeleteItem.defaultExpectation == nil {
		mmDeleteItem.defaultExpectation = &ItemRemoverMockDeleteItemExpectation{mock: mmDeleteItem.mock}
	}
	mmDeleteItem.defaultExpectation.results = &ItemRemoverMockDeleteItemResults{err}
	return mmDeleteItem.mock
}

// Set uses given function f to mock the itemRemover.DeleteItem method
func (mmDeleteItem *mItemRemoverMockDeleteItem) Set(f func(ctx context.Context, userID int64, skuID int64) (err error)) *ItemRemoverMock {
	if mmDeleteItem.defaultExpectation != nil {
		mmDeleteItem.mock.t.Fatalf("Default expectation is already set for the itemRemover.DeleteItem method")
	}

	if len(mmDeleteItem.expectations) > 0 {
		mmDeleteItem.mock.t.Fatalf("Some expectations are already set for the itemRemover.DeleteItem method")
	}

	mmDeleteItem.mock.funcDeleteItem = f
	return mmDeleteItem.mock
}

// When sets expectation for the itemRemover.DeleteItem which will trigger the result defined by the following
// Then helper
func (mmDeleteItem *mItemRemoverMockDeleteItem) When(ctx context.Context, userID int64, skuID int64) *ItemRemoverMockDeleteItemExpectation {
	if mmDeleteItem.mock.funcDeleteItem != nil {
		mmDeleteItem.mock.t.Fatalf("ItemRemoverMock.DeleteItem mock is already set by Set")
	}

	expectation := &ItemRemoverMockDeleteItemExpectation{
		mock:   mmDeleteItem.mock,
		params: &ItemRemoverMockDeleteItemParams{ctx, userID, skuID},
	}
	mmDeleteItem.expectations = append(mmDeleteItem.expectations, expectation)
	return expectation
}

// Then sets up itemRemover.DeleteItem return parameters for the expectation previously defined by the When method
func (e *ItemRemoverMockDeleteItemExpectation) Then(err error) *ItemRemoverMock {
	e.results = &ItemRemoverMockDeleteItemResults{err}
	return e.mock
}

// Times sets number of times itemRemover.DeleteItem should be invoked
func (mmDeleteItem *mItemRemoverMockDeleteItem) Times(n uint64) *mItemRemoverMockDeleteItem {
	if n == 0 {
		mmDeleteItem.mock.t.Fatalf("Times of ItemRemoverMock.DeleteItem mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmDeleteItem.expectedInvocations, n)
	return mmDeleteItem
}

func (mmDeleteItem *mItemRemoverMockDeleteItem) invocationsDone() bool {
	if len(mmDeleteItem.expectations) == 0 && mmDeleteItem.defaultExpectation == nil && mmDeleteItem.mock.funcDeleteItem == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmDeleteItem.mock.afterDeleteItemCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmDeleteItem.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// DeleteItem implements itemRemover
func (mmDeleteItem *ItemRemoverMock) DeleteItem(ctx context.Context, userID int64, skuID int64) (err error) {
	mm_atomic.AddUint64(&mmDeleteItem.beforeDeleteItemCounter, 1)
	defer mm_atomic.AddUint64(&mmDeleteItem.afterDeleteItemCounter, 1)

	if mmDeleteItem.inspectFuncDeleteItem != nil {
		mmDeleteItem.inspectFuncDeleteItem(ctx, userID, skuID)
	}

	mm_params := ItemRemoverMockDeleteItemParams{ctx, userID, skuID}

	// Record call args
	mmDeleteItem.DeleteItemMock.mutex.Lock()
	mmDeleteItem.DeleteItemMock.callArgs = append(mmDeleteItem.DeleteItemMock.callArgs, &mm_params)
	mmDeleteItem.DeleteItemMock.mutex.Unlock()

	for _, e := range mmDeleteItem.DeleteItemMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmDeleteItem.DeleteItemMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmDeleteItem.DeleteItemMock.defaultExpectation.Counter, 1)
		mm_want := mmDeleteItem.DeleteItemMock.defaultExpectation.params
		mm_want_ptrs := mmDeleteItem.DeleteItemMock.defaultExpectation.paramPtrs

		mm_got := ItemRemoverMockDeleteItemParams{ctx, userID, skuID}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmDeleteItem.t.Errorf("ItemRemoverMock.DeleteItem got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.userID != nil && !minimock.Equal(*mm_want_ptrs.userID, mm_got.userID) {
				mmDeleteItem.t.Errorf("ItemRemoverMock.DeleteItem got unexpected parameter userID, want: %#v, got: %#v%s\n", *mm_want_ptrs.userID, mm_got.userID, minimock.Diff(*mm_want_ptrs.userID, mm_got.userID))
			}

			if mm_want_ptrs.skuID != nil && !minimock.Equal(*mm_want_ptrs.skuID, mm_got.skuID) {
				mmDeleteItem.t.Errorf("ItemRemoverMock.DeleteItem got unexpected parameter skuID, want: %#v, got: %#v%s\n", *mm_want_ptrs.skuID, mm_got.skuID, minimock.Diff(*mm_want_ptrs.skuID, mm_got.skuID))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmDeleteItem.t.Errorf("ItemRemoverMock.DeleteItem got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmDeleteItem.DeleteItemMock.defaultExpectation.results
		if mm_results == nil {
			mmDeleteItem.t.Fatal("No results are set for the ItemRemoverMock.DeleteItem")
		}
		return (*mm_results).err
	}
	if mmDeleteItem.funcDeleteItem != nil {
		return mmDeleteItem.funcDeleteItem(ctx, userID, skuID)
	}
	mmDeleteItem.t.Fatalf("Unexpected call to ItemRemoverMock.DeleteItem. %v %v %v", ctx, userID, skuID)
	return
}

// DeleteItemAfterCounter returns a count of finished ItemRemoverMock.DeleteItem invocations
func (mmDeleteItem *ItemRemoverMock) DeleteItemAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDeleteItem.afterDeleteItemCounter)
}

// DeleteItemBeforeCounter returns a count of ItemRemoverMock.DeleteItem invocations
func (mmDeleteItem *ItemRemoverMock) DeleteItemBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDeleteItem.beforeDeleteItemCounter)
}

// Calls returns a list of arguments used in each call to ItemRemoverMock.DeleteItem.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmDeleteItem *mItemRemoverMockDeleteItem) Calls() []*ItemRemoverMockDeleteItemParams {
	mmDeleteItem.mutex.RLock()

	argCopy := make([]*ItemRemoverMockDeleteItemParams, len(mmDeleteItem.callArgs))
	copy(argCopy, mmDeleteItem.callArgs)

	mmDeleteItem.mutex.RUnlock()

	return argCopy
}

// MinimockDeleteItemDone returns true if the count of the DeleteItem invocations corresponds
// the number of defined expectations
func (m *ItemRemoverMock) MinimockDeleteItemDone() bool {
	if m.DeleteItemMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.DeleteItemMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.DeleteItemMock.invocationsDone()
}

// MinimockDeleteItemInspect logs each unmet expectation
func (m *ItemRemoverMock) MinimockDeleteItemInspect() {
	for _, e := range m.DeleteItemMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ItemRemoverMock.DeleteItem with params: %#v", *e.params)
		}
	}

	afterDeleteItemCounter := mm_atomic.LoadUint64(&m.afterDeleteItemCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.DeleteItemMock.defaultExpectation != nil && afterDeleteItemCounter < 1 {
		if m.DeleteItemMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ItemRemoverMock.DeleteItem")
		} else {
			m.t.Errorf("Expected call to ItemRemoverMock.DeleteItem with params: %#v", *m.DeleteItemMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDeleteItem != nil && afterDeleteItemCounter < 1 {
		m.t.Error("Expected call to ItemRemoverMock.DeleteItem")
	}

	if !m.DeleteItemMock.invocationsDone() && afterDeleteItemCounter > 0 {
		m.t.Errorf("Expected %d calls to ItemRemoverMock.DeleteItem but found %d calls",
			mm_atomic.LoadUint64(&m.DeleteItemMock.expectedInvocations), afterDeleteItemCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ItemRemoverMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockDeleteItemInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ItemRemoverMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *ItemRemoverMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockDeleteItemDone()
}