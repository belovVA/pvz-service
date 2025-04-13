// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "pvz-service/internal/model"

	uuid "github.com/google/uuid"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// AddNewPvz provides a mock function with given fields: ctx, city
func (_m *Service) AddNewPvz(ctx context.Context, city string) (*model.Pvz, error) {
	ret := _m.Called(ctx, city)

	if len(ret) == 0 {
		panic("no return value specified for AddNewPvz")
	}

	var r0 *model.Pvz
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Pvz, error)); ok {
		return rf(ctx, city)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Pvz); ok {
		r0 = rf(ctx, city)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Pvz)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, city)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddProduct provides a mock function with given fields: ctx, typeProduct, pvzID
func (_m *Service) AddProduct(ctx context.Context, typeProduct string, pvzID uuid.UUID) (*model.Product, error) {
	ret := _m.Called(ctx, typeProduct, pvzID)

	if len(ret) == 0 {
		panic("no return value specified for AddProduct")
	}

	var r0 *model.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uuid.UUID) (*model.Product, error)); ok {
		return rf(ctx, typeProduct, pvzID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, uuid.UUID) *model.Product); ok {
		r0 = rf(ctx, typeProduct, pvzID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, uuid.UUID) error); ok {
		r1 = rf(ctx, typeProduct, pvzID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Authenticate provides a mock function with given fields: ctx, user
func (_m *Service) Authenticate(ctx context.Context, user model.User) (string, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Authenticate")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.User) (string, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.User) string); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CloseReception provides a mock function with given fields: ctx, pvzID
func (_m *Service) CloseReception(ctx context.Context, pvzID uuid.UUID) (*model.Reception, error) {
	ret := _m.Called(ctx, pvzID)

	if len(ret) == 0 {
		panic("no return value specified for CloseReception")
	}

	var r0 *model.Reception
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*model.Reception, error)); ok {
		return rf(ctx, pvzID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *model.Reception); ok {
		r0 = rf(ctx, pvzID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Reception)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, pvzID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateReception provides a mock function with given fields: ctx, pvzID
func (_m *Service) CreateReception(ctx context.Context, pvzID uuid.UUID) (*model.Reception, error) {
	ret := _m.Called(ctx, pvzID)

	if len(ret) == 0 {
		panic("no return value specified for CreateReception")
	}

	var r0 *model.Reception
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*model.Reception, error)); ok {
		return rf(ctx, pvzID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *model.Reception); ok {
		r0 = rf(ctx, pvzID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Reception)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, pvzID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteProduct provides a mock function with given fields: ctx, pvzID
func (_m *Service) DeleteProduct(ctx context.Context, pvzID uuid.UUID) error {
	ret := _m.Called(ctx, pvzID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteProduct")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, pvzID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DummyAuth provides a mock function with given fields: ctx, role
func (_m *Service) DummyAuth(ctx context.Context, role string) (string, error) {
	ret := _m.Called(ctx, role)

	if len(ret) == 0 {
		panic("no return value specified for DummyAuth")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, role)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, role)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, role)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInfoPvz provides a mock function with given fields: ctx, query
func (_m *Service) GetInfoPvz(ctx context.Context, query *model.PvzInfoQuery) ([]*model.Pvz, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for GetInfoPvz")
	}

	var r0 []*model.Pvz
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.PvzInfoQuery) ([]*model.Pvz, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.PvzInfoQuery) []*model.Pvz); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Pvz)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.PvzInfoQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Registration provides a mock function with given fields: ctx, user
func (_m *Service) Registration(ctx context.Context, user model.User) (*model.User, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Registration")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.User) (*model.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.User) *model.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
