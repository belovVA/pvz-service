// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"
	model "pvz-service/internal/model"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// ProductRepository is an autogenerated mock type for the ProductRepository type
type ProductRepository struct {
	mock.Mock
}

// CreateProduct provides a mock function with given fields: ctx, typeProduct, recepID
func (_m *ProductRepository) CreateProduct(ctx context.Context, typeProduct string, recepID uuid.UUID) (uuid.UUID, error) {
	ret := _m.Called(ctx, typeProduct, recepID)

	if len(ret) == 0 {
		panic("no return value specified for CreateProduct")
	}

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uuid.UUID) (uuid.UUID, error)); ok {
		return rf(ctx, typeProduct, recepID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, uuid.UUID) uuid.UUID); ok {
		r0 = rf(ctx, typeProduct, recepID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, uuid.UUID) error); ok {
		r1 = rf(ctx, typeProduct, recepID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteProductByID provides a mock function with given fields: ctx, id
func (_m *ProductRepository) DeleteProductByID(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteProductByID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetLastProduct provides a mock function with given fields: ctx, receptionID
func (_m *ProductRepository) GetLastProduct(ctx context.Context, receptionID uuid.UUID) (*model.Product, error) {
	ret := _m.Called(ctx, receptionID)

	if len(ret) == 0 {
		panic("no return value specified for GetLastProduct")
	}

	var r0 *model.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*model.Product, error)); ok {
		return rf(ctx, receptionID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *model.Product); ok {
		r0 = rf(ctx, receptionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, receptionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProductByID provides a mock function with given fields: ctx, id
func (_m *ProductRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetProductByID")
	}

	var r0 *model.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*model.Product, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *model.Product); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProductSliceByReceptionID provides a mock function with given fields: ctx, receptionID
func (_m *ProductRepository) GetProductSliceByReceptionID(ctx context.Context, receptionID uuid.UUID) ([]model.Product, error) {
	ret := _m.Called(ctx, receptionID)

	if len(ret) == 0 {
		panic("no return value specified for GetProductSliceByReceptionID")
	}

	var r0 []model.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]model.Product, error)); ok {
		return rf(ctx, receptionID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []model.Product); ok {
		r0 = rf(ctx, receptionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, receptionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProductRepository creates a new instance of ProductRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProductRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProductRepository {
	mock := &ProductRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
