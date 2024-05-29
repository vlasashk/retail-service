// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "route256/cart/internal/cart/models"
)

// CartRetriever is an autogenerated mock type for the CartRetriever type
type CartRetriever struct {
	mock.Mock
}

// GetItemsByUserID provides a mock function with given fields: ctx, userID
func (_m *CartRetriever) GetItemsByUserID(ctx context.Context, userID int64) ([]models.Item, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetItemsByUserID")
	}

	var r0 []models.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]models.Item, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []models.Item); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCartRetriever creates a new instance of CartRetriever. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCartRetriever(t interface {
	mock.TestingT
	Cleanup(func())
}) *CartRetriever {
	mock := &CartRetriever{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}