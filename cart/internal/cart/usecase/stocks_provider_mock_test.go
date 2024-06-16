// Code generated by http://github.com/gojuno/minimock (v3.3.12). DO NOT EDIT.

package usecase_test

//go:generate minimock -i route256/cart/internal/cart/usecase.stocksProvider -o stocks_provider_mock_test.go -n StocksProviderMock -p usecase_test

import (
	"context"
	"route256/cart/internal/cart/models"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// StocksProviderMock implements stocksProvider
type StocksProviderMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcOrderCreate          func(ctx context.Context, order models.Order) (i1 int64, err error)
	inspectFuncOrderCreate   func(ctx context.Context, order models.Order)
	afterOrderCreateCounter  uint64
	beforeOrderCreateCounter uint64
	OrderCreateMock          mStocksProviderMockOrderCreate

	funcStocksInfo          func(ctx context.Context, skuID int64) (u1 uint64, err error)
	inspectFuncStocksInfo   func(ctx context.Context, skuID int64)
	afterStocksInfoCounter  uint64
	beforeStocksInfoCounter uint64
	StocksInfoMock          mStocksProviderMockStocksInfo
}

// NewStocksProviderMock returns a mock for stocksProvider
func NewStocksProviderMock(t minimock.Tester) *StocksProviderMock {
	m := &StocksProviderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.OrderCreateMock = mStocksProviderMockOrderCreate{mock: m}
	m.OrderCreateMock.callArgs = []*StocksProviderMockOrderCreateParams{}

	m.StocksInfoMock = mStocksProviderMockStocksInfo{mock: m}
	m.StocksInfoMock.callArgs = []*StocksProviderMockStocksInfoParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mStocksProviderMockOrderCreate struct {
	optional           bool
	mock               *StocksProviderMock
	defaultExpectation *StocksProviderMockOrderCreateExpectation
	expectations       []*StocksProviderMockOrderCreateExpectation

	callArgs []*StocksProviderMockOrderCreateParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// StocksProviderMockOrderCreateExpectation specifies expectation struct of the stocksProvider.OrderCreate
type StocksProviderMockOrderCreateExpectation struct {
	mock      *StocksProviderMock
	params    *StocksProviderMockOrderCreateParams
	paramPtrs *StocksProviderMockOrderCreateParamPtrs
	results   *StocksProviderMockOrderCreateResults
	Counter   uint64
}

// StocksProviderMockOrderCreateParams contains parameters of the stocksProvider.OrderCreate
type StocksProviderMockOrderCreateParams struct {
	ctx   context.Context
	order models.Order
}

// StocksProviderMockOrderCreateParamPtrs contains pointers to parameters of the stocksProvider.OrderCreate
type StocksProviderMockOrderCreateParamPtrs struct {
	ctx   *context.Context
	order *models.Order
}

// StocksProviderMockOrderCreateResults contains results of the stocksProvider.OrderCreate
type StocksProviderMockOrderCreateResults struct {
	i1  int64
	err error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmOrderCreate *mStocksProviderMockOrderCreate) Optional() *mStocksProviderMockOrderCreate {
	mmOrderCreate.optional = true
	return mmOrderCreate
}

// Expect sets up expected params for stocksProvider.OrderCreate
func (mmOrderCreate *mStocksProviderMockOrderCreate) Expect(ctx context.Context, order models.Order) *mStocksProviderMockOrderCreate {
	if mmOrderCreate.mock.funcOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("StocksProviderMock.OrderCreate mock is already set by Set")
	}

	if mmOrderCreate.defaultExpectation == nil {
		mmOrderCreate.defaultExpectation = &StocksProviderMockOrderCreateExpectation{}
	}

	if mmOrderCreate.defaultExpectation.paramPtrs != nil {
		mmOrderCreate.mock.t.Fatalf("StocksProviderMock.OrderCreate mock is already set by ExpectParams functions")
	}

	mmOrderCreate.defaultExpectation.params = &StocksProviderMockOrderCreateParams{ctx, order}
	for _, e := range mmOrderCreate.expectations {
		if minimock.Equal(e.params, mmOrderCreate.defaultExpectation.params) {
			mmOrderCreate.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmOrderCreate.defaultExpectation.params)
		}
	}

	return mmOrderCreate
}

// ExpectCtxParam1 sets up expected param ctx for stocksProvider.OrderCreate
func (mmOrderCreate *mStocksProviderMockOrderCreate) ExpectCtxParam1(ctx context.Context) *mStocksProviderMockOrderCreate {
	if mmOrderCreate.mock.funcOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("StocksProviderMock.OrderCreate mock is already set by Set")
	}

	if mmOrderCreate.defaultExpectation == nil {
		mmOrderCreate.defaultExpectation = &StocksProviderMockOrderCreateExpectation{}
	}

	if mmOrderCreate.defaultExpectation.params != nil {
		mmOrderCreate.mock.t.Fatalf("StocksProviderMock.OrderCreate mock is already set by Expect")
	}

	if mmOrderCreate.defaultExpectation.paramPtrs == nil {
		mmOrderCreate.defaultExpectation.paramPtrs = &StocksProviderMockOrderCreateParamPtrs{}
	}
	mmOrderCreate.defaultExpectation.paramPtrs.ctx = &ctx

	return mmOrderCreate
}

// ExpectOrderParam2 sets up expected param order for stocksProvider.OrderCreate
func (mmOrderCreate *mStocksProviderMockOrderCreate) ExpectOrderParam2(order models.Order) *mStocksProviderMockOrderCreate {
	if mmOrderCreate.mock.funcOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("StocksProviderMock.OrderCreate mock is already set by Set")
	}

	if mmOrderCreate.defaultExpectation == nil {
		mmOrderCreate.defaultExpectation = &StocksProviderMockOrderCreateExpectation{}
	}

	if mmOrderCreate.defaultExpectation.params != nil {
		mmOrderCreate.mock.t.Fatalf("StocksProviderMock.OrderCreate mock is already set by Expect")
	}

	if mmOrderCreate.defaultExpectation.paramPtrs == nil {
		mmOrderCreate.defaultExpectation.paramPtrs = &StocksProviderMockOrderCreateParamPtrs{}
	}
	mmOrderCreate.defaultExpectation.paramPtrs.order = &order

	return mmOrderCreate
}

// Inspect accepts an inspector function that has same arguments as the stocksProvider.OrderCreate
func (mmOrderCreate *mStocksProviderMockOrderCreate) Inspect(f func(ctx context.Context, order models.Order)) *mStocksProviderMockOrderCreate {
	if mmOrderCreate.mock.inspectFuncOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("Inspect function is already set for StocksProviderMock.OrderCreate")
	}

	mmOrderCreate.mock.inspectFuncOrderCreate = f

	return mmOrderCreate
}

// Return sets up results that will be returned by stocksProvider.OrderCreate
func (mmOrderCreate *mStocksProviderMockOrderCreate) Return(i1 int64, err error) *StocksProviderMock {
	if mmOrderCreate.mock.funcOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("StocksProviderMock.OrderCreate mock is already set by Set")
	}

	if mmOrderCreate.defaultExpectation == nil {
		mmOrderCreate.defaultExpectation = &StocksProviderMockOrderCreateExpectation{mock: mmOrderCreate.mock}
	}
	mmOrderCreate.defaultExpectation.results = &StocksProviderMockOrderCreateResults{i1, err}
	return mmOrderCreate.mock
}

// Set uses given function f to mock the stocksProvider.OrderCreate method
func (mmOrderCreate *mStocksProviderMockOrderCreate) Set(f func(ctx context.Context, order models.Order) (i1 int64, err error)) *StocksProviderMock {
	if mmOrderCreate.defaultExpectation != nil {
		mmOrderCreate.mock.t.Fatalf("Default expectation is already set for the stocksProvider.OrderCreate method")
	}

	if len(mmOrderCreate.expectations) > 0 {
		mmOrderCreate.mock.t.Fatalf("Some expectations are already set for the stocksProvider.OrderCreate method")
	}

	mmOrderCreate.mock.funcOrderCreate = f
	return mmOrderCreate.mock
}

// When sets expectation for the stocksProvider.OrderCreate which will trigger the result defined by the following
// Then helper
func (mmOrderCreate *mStocksProviderMockOrderCreate) When(ctx context.Context, order models.Order) *StocksProviderMockOrderCreateExpectation {
	if mmOrderCreate.mock.funcOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("StocksProviderMock.OrderCreate mock is already set by Set")
	}

	expectation := &StocksProviderMockOrderCreateExpectation{
		mock:   mmOrderCreate.mock,
		params: &StocksProviderMockOrderCreateParams{ctx, order},
	}
	mmOrderCreate.expectations = append(mmOrderCreate.expectations, expectation)
	return expectation
}

// Then sets up stocksProvider.OrderCreate return parameters for the expectation previously defined by the When method
func (e *StocksProviderMockOrderCreateExpectation) Then(i1 int64, err error) *StocksProviderMock {
	e.results = &StocksProviderMockOrderCreateResults{i1, err}
	return e.mock
}

// Times sets number of times stocksProvider.OrderCreate should be invoked
func (mmOrderCreate *mStocksProviderMockOrderCreate) Times(n uint64) *mStocksProviderMockOrderCreate {
	if n == 0 {
		mmOrderCreate.mock.t.Fatalf("Times of StocksProviderMock.OrderCreate mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmOrderCreate.expectedInvocations, n)
	return mmOrderCreate
}

func (mmOrderCreate *mStocksProviderMockOrderCreate) invocationsDone() bool {
	if len(mmOrderCreate.expectations) == 0 && mmOrderCreate.defaultExpectation == nil && mmOrderCreate.mock.funcOrderCreate == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmOrderCreate.mock.afterOrderCreateCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmOrderCreate.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// OrderCreate implements stocksProvider
func (mmOrderCreate *StocksProviderMock) OrderCreate(ctx context.Context, order models.Order) (i1 int64, err error) {
	mm_atomic.AddUint64(&mmOrderCreate.beforeOrderCreateCounter, 1)
	defer mm_atomic.AddUint64(&mmOrderCreate.afterOrderCreateCounter, 1)

	if mmOrderCreate.inspectFuncOrderCreate != nil {
		mmOrderCreate.inspectFuncOrderCreate(ctx, order)
	}

	mm_params := StocksProviderMockOrderCreateParams{ctx, order}

	// Record call args
	mmOrderCreate.OrderCreateMock.mutex.Lock()
	mmOrderCreate.OrderCreateMock.callArgs = append(mmOrderCreate.OrderCreateMock.callArgs, &mm_params)
	mmOrderCreate.OrderCreateMock.mutex.Unlock()

	for _, e := range mmOrderCreate.OrderCreateMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.i1, e.results.err
		}
	}

	if mmOrderCreate.OrderCreateMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmOrderCreate.OrderCreateMock.defaultExpectation.Counter, 1)
		mm_want := mmOrderCreate.OrderCreateMock.defaultExpectation.params
		mm_want_ptrs := mmOrderCreate.OrderCreateMock.defaultExpectation.paramPtrs

		mm_got := StocksProviderMockOrderCreateParams{ctx, order}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmOrderCreate.t.Errorf("StocksProviderMock.OrderCreate got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.order != nil && !minimock.Equal(*mm_want_ptrs.order, mm_got.order) {
				mmOrderCreate.t.Errorf("StocksProviderMock.OrderCreate got unexpected parameter order, want: %#v, got: %#v%s\n", *mm_want_ptrs.order, mm_got.order, minimock.Diff(*mm_want_ptrs.order, mm_got.order))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmOrderCreate.t.Errorf("StocksProviderMock.OrderCreate got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmOrderCreate.OrderCreateMock.defaultExpectation.results
		if mm_results == nil {
			mmOrderCreate.t.Fatal("No results are set for the StocksProviderMock.OrderCreate")
		}
		return (*mm_results).i1, (*mm_results).err
	}
	if mmOrderCreate.funcOrderCreate != nil {
		return mmOrderCreate.funcOrderCreate(ctx, order)
	}
	mmOrderCreate.t.Fatalf("Unexpected call to StocksProviderMock.OrderCreate. %v %v", ctx, order)
	return
}

// OrderCreateAfterCounter returns a count of finished StocksProviderMock.OrderCreate invocations
func (mmOrderCreate *StocksProviderMock) OrderCreateAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmOrderCreate.afterOrderCreateCounter)
}

// OrderCreateBeforeCounter returns a count of StocksProviderMock.OrderCreate invocations
func (mmOrderCreate *StocksProviderMock) OrderCreateBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmOrderCreate.beforeOrderCreateCounter)
}

// Calls returns a list of arguments used in each call to StocksProviderMock.OrderCreate.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmOrderCreate *mStocksProviderMockOrderCreate) Calls() []*StocksProviderMockOrderCreateParams {
	mmOrderCreate.mutex.RLock()

	argCopy := make([]*StocksProviderMockOrderCreateParams, len(mmOrderCreate.callArgs))
	copy(argCopy, mmOrderCreate.callArgs)

	mmOrderCreate.mutex.RUnlock()

	return argCopy
}

// MinimockOrderCreateDone returns true if the count of the OrderCreate invocations corresponds
// the number of defined expectations
func (m *StocksProviderMock) MinimockOrderCreateDone() bool {
	if m.OrderCreateMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.OrderCreateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.OrderCreateMock.invocationsDone()
}

// MinimockOrderCreateInspect logs each unmet expectation
func (m *StocksProviderMock) MinimockOrderCreateInspect() {
	for _, e := range m.OrderCreateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to StocksProviderMock.OrderCreate with params: %#v", *e.params)
		}
	}

	afterOrderCreateCounter := mm_atomic.LoadUint64(&m.afterOrderCreateCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.OrderCreateMock.defaultExpectation != nil && afterOrderCreateCounter < 1 {
		if m.OrderCreateMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to StocksProviderMock.OrderCreate")
		} else {
			m.t.Errorf("Expected call to StocksProviderMock.OrderCreate with params: %#v", *m.OrderCreateMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcOrderCreate != nil && afterOrderCreateCounter < 1 {
		m.t.Error("Expected call to StocksProviderMock.OrderCreate")
	}

	if !m.OrderCreateMock.invocationsDone() && afterOrderCreateCounter > 0 {
		m.t.Errorf("Expected %d calls to StocksProviderMock.OrderCreate but found %d calls",
			mm_atomic.LoadUint64(&m.OrderCreateMock.expectedInvocations), afterOrderCreateCounter)
	}
}

type mStocksProviderMockStocksInfo struct {
	optional           bool
	mock               *StocksProviderMock
	defaultExpectation *StocksProviderMockStocksInfoExpectation
	expectations       []*StocksProviderMockStocksInfoExpectation

	callArgs []*StocksProviderMockStocksInfoParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// StocksProviderMockStocksInfoExpectation specifies expectation struct of the stocksProvider.StocksInfo
type StocksProviderMockStocksInfoExpectation struct {
	mock      *StocksProviderMock
	params    *StocksProviderMockStocksInfoParams
	paramPtrs *StocksProviderMockStocksInfoParamPtrs
	results   *StocksProviderMockStocksInfoResults
	Counter   uint64
}

// StocksProviderMockStocksInfoParams contains parameters of the stocksProvider.StocksInfo
type StocksProviderMockStocksInfoParams struct {
	ctx   context.Context
	skuID int64
}

// StocksProviderMockStocksInfoParamPtrs contains pointers to parameters of the stocksProvider.StocksInfo
type StocksProviderMockStocksInfoParamPtrs struct {
	ctx   *context.Context
	skuID *int64
}

// StocksProviderMockStocksInfoResults contains results of the stocksProvider.StocksInfo
type StocksProviderMockStocksInfoResults struct {
	u1  uint64
	err error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmStocksInfo *mStocksProviderMockStocksInfo) Optional() *mStocksProviderMockStocksInfo {
	mmStocksInfo.optional = true
	return mmStocksInfo
}

// Expect sets up expected params for stocksProvider.StocksInfo
func (mmStocksInfo *mStocksProviderMockStocksInfo) Expect(ctx context.Context, skuID int64) *mStocksProviderMockStocksInfo {
	if mmStocksInfo.mock.funcStocksInfo != nil {
		mmStocksInfo.mock.t.Fatalf("StocksProviderMock.StocksInfo mock is already set by Set")
	}

	if mmStocksInfo.defaultExpectation == nil {
		mmStocksInfo.defaultExpectation = &StocksProviderMockStocksInfoExpectation{}
	}

	if mmStocksInfo.defaultExpectation.paramPtrs != nil {
		mmStocksInfo.mock.t.Fatalf("StocksProviderMock.StocksInfo mock is already set by ExpectParams functions")
	}

	mmStocksInfo.defaultExpectation.params = &StocksProviderMockStocksInfoParams{ctx, skuID}
	for _, e := range mmStocksInfo.expectations {
		if minimock.Equal(e.params, mmStocksInfo.defaultExpectation.params) {
			mmStocksInfo.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmStocksInfo.defaultExpectation.params)
		}
	}

	return mmStocksInfo
}

// ExpectCtxParam1 sets up expected param ctx for stocksProvider.StocksInfo
func (mmStocksInfo *mStocksProviderMockStocksInfo) ExpectCtxParam1(ctx context.Context) *mStocksProviderMockStocksInfo {
	if mmStocksInfo.mock.funcStocksInfo != nil {
		mmStocksInfo.mock.t.Fatalf("StocksProviderMock.StocksInfo mock is already set by Set")
	}

	if mmStocksInfo.defaultExpectation == nil {
		mmStocksInfo.defaultExpectation = &StocksProviderMockStocksInfoExpectation{}
	}

	if mmStocksInfo.defaultExpectation.params != nil {
		mmStocksInfo.mock.t.Fatalf("StocksProviderMock.StocksInfo mock is already set by Expect")
	}

	if mmStocksInfo.defaultExpectation.paramPtrs == nil {
		mmStocksInfo.defaultExpectation.paramPtrs = &StocksProviderMockStocksInfoParamPtrs{}
	}
	mmStocksInfo.defaultExpectation.paramPtrs.ctx = &ctx

	return mmStocksInfo
}

// ExpectSkuIDParam2 sets up expected param skuID for stocksProvider.StocksInfo
func (mmStocksInfo *mStocksProviderMockStocksInfo) ExpectSkuIDParam2(skuID int64) *mStocksProviderMockStocksInfo {
	if mmStocksInfo.mock.funcStocksInfo != nil {
		mmStocksInfo.mock.t.Fatalf("StocksProviderMock.StocksInfo mock is already set by Set")
	}

	if mmStocksInfo.defaultExpectation == nil {
		mmStocksInfo.defaultExpectation = &StocksProviderMockStocksInfoExpectation{}
	}

	if mmStocksInfo.defaultExpectation.params != nil {
		mmStocksInfo.mock.t.Fatalf("StocksProviderMock.StocksInfo mock is already set by Expect")
	}

	if mmStocksInfo.defaultExpectation.paramPtrs == nil {
		mmStocksInfo.defaultExpectation.paramPtrs = &StocksProviderMockStocksInfoParamPtrs{}
	}
	mmStocksInfo.defaultExpectation.paramPtrs.skuID = &skuID

	return mmStocksInfo
}

// Inspect accepts an inspector function that has same arguments as the stocksProvider.StocksInfo
func (mmStocksInfo *mStocksProviderMockStocksInfo) Inspect(f func(ctx context.Context, skuID int64)) *mStocksProviderMockStocksInfo {
	if mmStocksInfo.mock.inspectFuncStocksInfo != nil {
		mmStocksInfo.mock.t.Fatalf("Inspect function is already set for StocksProviderMock.StocksInfo")
	}

	mmStocksInfo.mock.inspectFuncStocksInfo = f

	return mmStocksInfo
}

// Return sets up results that will be returned by stocksProvider.StocksInfo
func (mmStocksInfo *mStocksProviderMockStocksInfo) Return(u1 uint64, err error) *StocksProviderMock {
	if mmStocksInfo.mock.funcStocksInfo != nil {
		mmStocksInfo.mock.t.Fatalf("StocksProviderMock.StocksInfo mock is already set by Set")
	}

	if mmStocksInfo.defaultExpectation == nil {
		mmStocksInfo.defaultExpectation = &StocksProviderMockStocksInfoExpectation{mock: mmStocksInfo.mock}
	}
	mmStocksInfo.defaultExpectation.results = &StocksProviderMockStocksInfoResults{u1, err}
	return mmStocksInfo.mock
}

// Set uses given function f to mock the stocksProvider.StocksInfo method
func (mmStocksInfo *mStocksProviderMockStocksInfo) Set(f func(ctx context.Context, skuID int64) (u1 uint64, err error)) *StocksProviderMock {
	if mmStocksInfo.defaultExpectation != nil {
		mmStocksInfo.mock.t.Fatalf("Default expectation is already set for the stocksProvider.StocksInfo method")
	}

	if len(mmStocksInfo.expectations) > 0 {
		mmStocksInfo.mock.t.Fatalf("Some expectations are already set for the stocksProvider.StocksInfo method")
	}

	mmStocksInfo.mock.funcStocksInfo = f
	return mmStocksInfo.mock
}

// When sets expectation for the stocksProvider.StocksInfo which will trigger the result defined by the following
// Then helper
func (mmStocksInfo *mStocksProviderMockStocksInfo) When(ctx context.Context, skuID int64) *StocksProviderMockStocksInfoExpectation {
	if mmStocksInfo.mock.funcStocksInfo != nil {
		mmStocksInfo.mock.t.Fatalf("StocksProviderMock.StocksInfo mock is already set by Set")
	}

	expectation := &StocksProviderMockStocksInfoExpectation{
		mock:   mmStocksInfo.mock,
		params: &StocksProviderMockStocksInfoParams{ctx, skuID},
	}
	mmStocksInfo.expectations = append(mmStocksInfo.expectations, expectation)
	return expectation
}

// Then sets up stocksProvider.StocksInfo return parameters for the expectation previously defined by the When method
func (e *StocksProviderMockStocksInfoExpectation) Then(u1 uint64, err error) *StocksProviderMock {
	e.results = &StocksProviderMockStocksInfoResults{u1, err}
	return e.mock
}

// Times sets number of times stocksProvider.StocksInfo should be invoked
func (mmStocksInfo *mStocksProviderMockStocksInfo) Times(n uint64) *mStocksProviderMockStocksInfo {
	if n == 0 {
		mmStocksInfo.mock.t.Fatalf("Times of StocksProviderMock.StocksInfo mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmStocksInfo.expectedInvocations, n)
	return mmStocksInfo
}

func (mmStocksInfo *mStocksProviderMockStocksInfo) invocationsDone() bool {
	if len(mmStocksInfo.expectations) == 0 && mmStocksInfo.defaultExpectation == nil && mmStocksInfo.mock.funcStocksInfo == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmStocksInfo.mock.afterStocksInfoCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmStocksInfo.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// StocksInfo implements stocksProvider
func (mmStocksInfo *StocksProviderMock) StocksInfo(ctx context.Context, skuID int64) (u1 uint64, err error) {
	mm_atomic.AddUint64(&mmStocksInfo.beforeStocksInfoCounter, 1)
	defer mm_atomic.AddUint64(&mmStocksInfo.afterStocksInfoCounter, 1)

	if mmStocksInfo.inspectFuncStocksInfo != nil {
		mmStocksInfo.inspectFuncStocksInfo(ctx, skuID)
	}

	mm_params := StocksProviderMockStocksInfoParams{ctx, skuID}

	// Record call args
	mmStocksInfo.StocksInfoMock.mutex.Lock()
	mmStocksInfo.StocksInfoMock.callArgs = append(mmStocksInfo.StocksInfoMock.callArgs, &mm_params)
	mmStocksInfo.StocksInfoMock.mutex.Unlock()

	for _, e := range mmStocksInfo.StocksInfoMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.u1, e.results.err
		}
	}

	if mmStocksInfo.StocksInfoMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmStocksInfo.StocksInfoMock.defaultExpectation.Counter, 1)
		mm_want := mmStocksInfo.StocksInfoMock.defaultExpectation.params
		mm_want_ptrs := mmStocksInfo.StocksInfoMock.defaultExpectation.paramPtrs

		mm_got := StocksProviderMockStocksInfoParams{ctx, skuID}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmStocksInfo.t.Errorf("StocksProviderMock.StocksInfo got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.skuID != nil && !minimock.Equal(*mm_want_ptrs.skuID, mm_got.skuID) {
				mmStocksInfo.t.Errorf("StocksProviderMock.StocksInfo got unexpected parameter skuID, want: %#v, got: %#v%s\n", *mm_want_ptrs.skuID, mm_got.skuID, minimock.Diff(*mm_want_ptrs.skuID, mm_got.skuID))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmStocksInfo.t.Errorf("StocksProviderMock.StocksInfo got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmStocksInfo.StocksInfoMock.defaultExpectation.results
		if mm_results == nil {
			mmStocksInfo.t.Fatal("No results are set for the StocksProviderMock.StocksInfo")
		}
		return (*mm_results).u1, (*mm_results).err
	}
	if mmStocksInfo.funcStocksInfo != nil {
		return mmStocksInfo.funcStocksInfo(ctx, skuID)
	}
	mmStocksInfo.t.Fatalf("Unexpected call to StocksProviderMock.StocksInfo. %v %v", ctx, skuID)
	return
}

// StocksInfoAfterCounter returns a count of finished StocksProviderMock.StocksInfo invocations
func (mmStocksInfo *StocksProviderMock) StocksInfoAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmStocksInfo.afterStocksInfoCounter)
}

// StocksInfoBeforeCounter returns a count of StocksProviderMock.StocksInfo invocations
func (mmStocksInfo *StocksProviderMock) StocksInfoBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmStocksInfo.beforeStocksInfoCounter)
}

// Calls returns a list of arguments used in each call to StocksProviderMock.StocksInfo.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmStocksInfo *mStocksProviderMockStocksInfo) Calls() []*StocksProviderMockStocksInfoParams {
	mmStocksInfo.mutex.RLock()

	argCopy := make([]*StocksProviderMockStocksInfoParams, len(mmStocksInfo.callArgs))
	copy(argCopy, mmStocksInfo.callArgs)

	mmStocksInfo.mutex.RUnlock()

	return argCopy
}

// MinimockStocksInfoDone returns true if the count of the StocksInfo invocations corresponds
// the number of defined expectations
func (m *StocksProviderMock) MinimockStocksInfoDone() bool {
	if m.StocksInfoMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.StocksInfoMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.StocksInfoMock.invocationsDone()
}

// MinimockStocksInfoInspect logs each unmet expectation
func (m *StocksProviderMock) MinimockStocksInfoInspect() {
	for _, e := range m.StocksInfoMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to StocksProviderMock.StocksInfo with params: %#v", *e.params)
		}
	}

	afterStocksInfoCounter := mm_atomic.LoadUint64(&m.afterStocksInfoCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.StocksInfoMock.defaultExpectation != nil && afterStocksInfoCounter < 1 {
		if m.StocksInfoMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to StocksProviderMock.StocksInfo")
		} else {
			m.t.Errorf("Expected call to StocksProviderMock.StocksInfo with params: %#v", *m.StocksInfoMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcStocksInfo != nil && afterStocksInfoCounter < 1 {
		m.t.Error("Expected call to StocksProviderMock.StocksInfo")
	}

	if !m.StocksInfoMock.invocationsDone() && afterStocksInfoCounter > 0 {
		m.t.Errorf("Expected %d calls to StocksProviderMock.StocksInfo but found %d calls",
			mm_atomic.LoadUint64(&m.StocksInfoMock.expectedInvocations), afterStocksInfoCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *StocksProviderMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockOrderCreateInspect()

			m.MinimockStocksInfoInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *StocksProviderMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *StocksProviderMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockOrderCreateDone() &&
		m.MinimockStocksInfoDone()
}
