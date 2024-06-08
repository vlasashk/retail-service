// Code generated by http://github.com/gojuno/minimock (v3.3.11). DO NOT EDIT.

package usecase_test

//go:generate minimock -i route256/cart/internal/cart/usecase.productProvider -o product_provider_mock_test.go -n ProductProviderMock -p usecase_test

import (
	"context"
	"route256/cart/internal/cart/models"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// ProductProviderMock implements productProvider
type ProductProviderMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetProduct          func(ctx context.Context, sku int64) (i1 models.ItemDescription, err error)
	inspectFuncGetProduct   func(ctx context.Context, sku int64)
	afterGetProductCounter  uint64
	beforeGetProductCounter uint64
	GetProductMock          mProductProviderMockGetProduct
}

// NewProductProviderMock returns a mock for productProvider
func NewProductProviderMock(t minimock.Tester) *ProductProviderMock {
	m := &ProductProviderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetProductMock = mProductProviderMockGetProduct{mock: m}
	m.GetProductMock.callArgs = []*ProductProviderMockGetProductParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mProductProviderMockGetProduct struct {
	optional           bool
	mock               *ProductProviderMock
	defaultExpectation *ProductProviderMockGetProductExpectation
	expectations       []*ProductProviderMockGetProductExpectation

	callArgs []*ProductProviderMockGetProductParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// ProductProviderMockGetProductExpectation specifies expectation struct of the productProvider.GetProduct
type ProductProviderMockGetProductExpectation struct {
	mock      *ProductProviderMock
	params    *ProductProviderMockGetProductParams
	paramPtrs *ProductProviderMockGetProductParamPtrs
	results   *ProductProviderMockGetProductResults
	Counter   uint64
}

// ProductProviderMockGetProductParams contains parameters of the productProvider.GetProduct
type ProductProviderMockGetProductParams struct {
	ctx context.Context
	sku int64
}

// ProductProviderMockGetProductParamPtrs contains pointers to parameters of the productProvider.GetProduct
type ProductProviderMockGetProductParamPtrs struct {
	ctx *context.Context
	sku *int64
}

// ProductProviderMockGetProductResults contains results of the productProvider.GetProduct
type ProductProviderMockGetProductResults struct {
	i1  models.ItemDescription
	err error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option by default unless you really need it, as it helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmGetProduct *mProductProviderMockGetProduct) Optional() *mProductProviderMockGetProduct {
	mmGetProduct.optional = true
	return mmGetProduct
}

// Expect sets up expected params for productProvider.GetProduct
func (mmGetProduct *mProductProviderMockGetProduct) Expect(ctx context.Context, sku int64) *mProductProviderMockGetProduct {
	if mmGetProduct.mock.funcGetProduct != nil {
		mmGetProduct.mock.t.Fatalf("ProductProviderMock.GetProduct mock is already set by Set")
	}

	if mmGetProduct.defaultExpectation == nil {
		mmGetProduct.defaultExpectation = &ProductProviderMockGetProductExpectation{}
	}

	if mmGetProduct.defaultExpectation.paramPtrs != nil {
		mmGetProduct.mock.t.Fatalf("ProductProviderMock.GetProduct mock is already set by ExpectParams functions")
	}

	mmGetProduct.defaultExpectation.params = &ProductProviderMockGetProductParams{ctx, sku}
	for _, e := range mmGetProduct.expectations {
		if minimock.Equal(e.params, mmGetProduct.defaultExpectation.params) {
			mmGetProduct.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetProduct.defaultExpectation.params)
		}
	}

	return mmGetProduct
}

// ExpectCtxParam1 sets up expected param ctx for productProvider.GetProduct
func (mmGetProduct *mProductProviderMockGetProduct) ExpectCtxParam1(ctx context.Context) *mProductProviderMockGetProduct {
	if mmGetProduct.mock.funcGetProduct != nil {
		mmGetProduct.mock.t.Fatalf("ProductProviderMock.GetProduct mock is already set by Set")
	}

	if mmGetProduct.defaultExpectation == nil {
		mmGetProduct.defaultExpectation = &ProductProviderMockGetProductExpectation{}
	}

	if mmGetProduct.defaultExpectation.params != nil {
		mmGetProduct.mock.t.Fatalf("ProductProviderMock.GetProduct mock is already set by Expect")
	}

	if mmGetProduct.defaultExpectation.paramPtrs == nil {
		mmGetProduct.defaultExpectation.paramPtrs = &ProductProviderMockGetProductParamPtrs{}
	}
	mmGetProduct.defaultExpectation.paramPtrs.ctx = &ctx

	return mmGetProduct
}

// ExpectSkuParam2 sets up expected param sku for productProvider.GetProduct
func (mmGetProduct *mProductProviderMockGetProduct) ExpectSkuParam2(sku int64) *mProductProviderMockGetProduct {
	if mmGetProduct.mock.funcGetProduct != nil {
		mmGetProduct.mock.t.Fatalf("ProductProviderMock.GetProduct mock is already set by Set")
	}

	if mmGetProduct.defaultExpectation == nil {
		mmGetProduct.defaultExpectation = &ProductProviderMockGetProductExpectation{}
	}

	if mmGetProduct.defaultExpectation.params != nil {
		mmGetProduct.mock.t.Fatalf("ProductProviderMock.GetProduct mock is already set by Expect")
	}

	if mmGetProduct.defaultExpectation.paramPtrs == nil {
		mmGetProduct.defaultExpectation.paramPtrs = &ProductProviderMockGetProductParamPtrs{}
	}
	mmGetProduct.defaultExpectation.paramPtrs.sku = &sku

	return mmGetProduct
}

// Inspect accepts an inspector function that has same arguments as the productProvider.GetProduct
func (mmGetProduct *mProductProviderMockGetProduct) Inspect(f func(ctx context.Context, sku int64)) *mProductProviderMockGetProduct {
	if mmGetProduct.mock.inspectFuncGetProduct != nil {
		mmGetProduct.mock.t.Fatalf("Inspect function is already set for ProductProviderMock.GetProduct")
	}

	mmGetProduct.mock.inspectFuncGetProduct = f

	return mmGetProduct
}

// Return sets up results that will be returned by productProvider.GetProduct
func (mmGetProduct *mProductProviderMockGetProduct) Return(i1 models.ItemDescription, err error) *ProductProviderMock {
	if mmGetProduct.mock.funcGetProduct != nil {
		mmGetProduct.mock.t.Fatalf("ProductProviderMock.GetProduct mock is already set by Set")
	}

	if mmGetProduct.defaultExpectation == nil {
		mmGetProduct.defaultExpectation = &ProductProviderMockGetProductExpectation{mock: mmGetProduct.mock}
	}
	mmGetProduct.defaultExpectation.results = &ProductProviderMockGetProductResults{i1, err}
	return mmGetProduct.mock
}

// Set uses given function f to mock the productProvider.GetProduct method
func (mmGetProduct *mProductProviderMockGetProduct) Set(f func(ctx context.Context, sku int64) (i1 models.ItemDescription, err error)) *ProductProviderMock {
	if mmGetProduct.defaultExpectation != nil {
		mmGetProduct.mock.t.Fatalf("Default expectation is already set for the productProvider.GetProduct method")
	}

	if len(mmGetProduct.expectations) > 0 {
		mmGetProduct.mock.t.Fatalf("Some expectations are already set for the productProvider.GetProduct method")
	}

	mmGetProduct.mock.funcGetProduct = f
	return mmGetProduct.mock
}

// When sets expectation for the productProvider.GetProduct which will trigger the result defined by the following
// Then helper
func (mmGetProduct *mProductProviderMockGetProduct) When(ctx context.Context, sku int64) *ProductProviderMockGetProductExpectation {
	if mmGetProduct.mock.funcGetProduct != nil {
		mmGetProduct.mock.t.Fatalf("ProductProviderMock.GetProduct mock is already set by Set")
	}

	expectation := &ProductProviderMockGetProductExpectation{
		mock:   mmGetProduct.mock,
		params: &ProductProviderMockGetProductParams{ctx, sku},
	}
	mmGetProduct.expectations = append(mmGetProduct.expectations, expectation)
	return expectation
}

// Then sets up productProvider.GetProduct return parameters for the expectation previously defined by the When method
func (e *ProductProviderMockGetProductExpectation) Then(i1 models.ItemDescription, err error) *ProductProviderMock {
	e.results = &ProductProviderMockGetProductResults{i1, err}
	return e.mock
}

// Times sets number of times productProvider.GetProduct should be invoked
func (mmGetProduct *mProductProviderMockGetProduct) Times(n uint64) *mProductProviderMockGetProduct {
	if n == 0 {
		mmGetProduct.mock.t.Fatalf("Times of ProductProviderMock.GetProduct mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmGetProduct.expectedInvocations, n)
	return mmGetProduct
}

func (mmGetProduct *mProductProviderMockGetProduct) invocationsDone() bool {
	if len(mmGetProduct.expectations) == 0 && mmGetProduct.defaultExpectation == nil && mmGetProduct.mock.funcGetProduct == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmGetProduct.mock.afterGetProductCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmGetProduct.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// GetProduct implements productProvider
func (mmGetProduct *ProductProviderMock) GetProduct(ctx context.Context, sku int64) (i1 models.ItemDescription, err error) {
	mm_atomic.AddUint64(&mmGetProduct.beforeGetProductCounter, 1)
	defer mm_atomic.AddUint64(&mmGetProduct.afterGetProductCounter, 1)

	if mmGetProduct.inspectFuncGetProduct != nil {
		mmGetProduct.inspectFuncGetProduct(ctx, sku)
	}

	mm_params := ProductProviderMockGetProductParams{ctx, sku}

	// Record call args
	mmGetProduct.GetProductMock.mutex.Lock()
	mmGetProduct.GetProductMock.callArgs = append(mmGetProduct.GetProductMock.callArgs, &mm_params)
	mmGetProduct.GetProductMock.mutex.Unlock()

	for _, e := range mmGetProduct.GetProductMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.i1, e.results.err
		}
	}

	if mmGetProduct.GetProductMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetProduct.GetProductMock.defaultExpectation.Counter, 1)
		mm_want := mmGetProduct.GetProductMock.defaultExpectation.params
		mm_want_ptrs := mmGetProduct.GetProductMock.defaultExpectation.paramPtrs

		mm_got := ProductProviderMockGetProductParams{ctx, sku}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmGetProduct.t.Errorf("ProductProviderMock.GetProduct got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.sku != nil && !minimock.Equal(*mm_want_ptrs.sku, mm_got.sku) {
				mmGetProduct.t.Errorf("ProductProviderMock.GetProduct got unexpected parameter sku, want: %#v, got: %#v%s\n", *mm_want_ptrs.sku, mm_got.sku, minimock.Diff(*mm_want_ptrs.sku, mm_got.sku))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetProduct.t.Errorf("ProductProviderMock.GetProduct got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetProduct.GetProductMock.defaultExpectation.results
		if mm_results == nil {
			mmGetProduct.t.Fatal("No results are set for the ProductProviderMock.GetProduct")
		}
		return (*mm_results).i1, (*mm_results).err
	}
	if mmGetProduct.funcGetProduct != nil {
		return mmGetProduct.funcGetProduct(ctx, sku)
	}
	mmGetProduct.t.Fatalf("Unexpected call to ProductProviderMock.GetProduct. %v %v", ctx, sku)
	return
}

// GetProductAfterCounter returns a count of finished ProductProviderMock.GetProduct invocations
func (mmGetProduct *ProductProviderMock) GetProductAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetProduct.afterGetProductCounter)
}

// GetProductBeforeCounter returns a count of ProductProviderMock.GetProduct invocations
func (mmGetProduct *ProductProviderMock) GetProductBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetProduct.beforeGetProductCounter)
}

// Calls returns a list of arguments used in each call to ProductProviderMock.GetProduct.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetProduct *mProductProviderMockGetProduct) Calls() []*ProductProviderMockGetProductParams {
	mmGetProduct.mutex.RLock()

	argCopy := make([]*ProductProviderMockGetProductParams, len(mmGetProduct.callArgs))
	copy(argCopy, mmGetProduct.callArgs)

	mmGetProduct.mutex.RUnlock()

	return argCopy
}

// MinimockGetProductDone returns true if the count of the GetProduct invocations corresponds
// the number of defined expectations
func (m *ProductProviderMock) MinimockGetProductDone() bool {
	if m.GetProductMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.GetProductMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.GetProductMock.invocationsDone()
}

// MinimockGetProductInspect logs each unmet expectation
func (m *ProductProviderMock) MinimockGetProductInspect() {
	for _, e := range m.GetProductMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ProductProviderMock.GetProduct with params: %#v", *e.params)
		}
	}

	afterGetProductCounter := mm_atomic.LoadUint64(&m.afterGetProductCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.GetProductMock.defaultExpectation != nil && afterGetProductCounter < 1 {
		if m.GetProductMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ProductProviderMock.GetProduct")
		} else {
			m.t.Errorf("Expected call to ProductProviderMock.GetProduct with params: %#v", *m.GetProductMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetProduct != nil && afterGetProductCounter < 1 {
		m.t.Error("Expected call to ProductProviderMock.GetProduct")
	}

	if !m.GetProductMock.invocationsDone() && afterGetProductCounter > 0 {
		m.t.Errorf("Expected %d calls to ProductProviderMock.GetProduct but found %d calls",
			mm_atomic.LoadUint64(&m.GetProductMock.expectedInvocations), afterGetProductCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ProductProviderMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetProductInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ProductProviderMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *ProductProviderMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetProductDone()
}
