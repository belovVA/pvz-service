// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"
	time "time"

	model "pvz-service/internal/model"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// ReceptionRepository is an autogenerated mock type for the ReceptionRepository type
type ReceptionRepository struct {
	mock.Mock
}

// CloseReception provides a mock function with given fields: ctx, receptionID
func (_m *ReceptionRepository) CloseReception(ctx context.Context, receptionID uuid.UUID) error {
	ret := _m.Called(ctx, receptionID)

	if len(ret) == 0 {
		panic("no return value specified for CloseReception")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, receptionID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateReception provides a mock function with given fields: ctx, pvzID
func (_m *ReceptionRepository) CreateReception(ctx context.Context, pvzID uuid.UUID) (uuid.UUID, error) {
	ret := _m.Called(ctx, pvzID)

	if len(ret) == 0 {
		panic("no return value specified for CreateReception")
	}

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (uuid.UUID, error)); ok {
		return rf(ctx, pvzID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) uuid.UUID); ok {
		r0 = rf(ctx, pvzID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, pvzID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastReception provides a mock function with given fields: ctx, pvzID
func (_m *ReceptionRepository) GetLastReception(ctx context.Context, pvzID uuid.UUID) (*model.Reception, error) {
	ret := _m.Called(ctx, pvzID)

	if len(ret) == 0 {
		panic("no return value specified for GetLastReception")
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

// GetReceptionByID provides a mock function with given fields: ctx, id
func (_m *ReceptionRepository) GetReceptionByID(ctx context.Context, id uuid.UUID) (*model.Reception, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetReceptionByID")
	}

	var r0 *model.Reception
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*model.Reception, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *model.Reception); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Reception)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetReceptionsSliceWithTimeRange provides a mock function with given fields: ctx, begin, end
func (_m *ReceptionRepository) GetReceptionsSliceWithTimeRange(ctx context.Context, begin time.Time, end time.Time) ([]model.Reception, error) {
	ret := _m.Called(ctx, begin, end)

	if len(ret) == 0 {
		panic("no return value specified for GetReceptionsSliceWithTimeRange")
	}

	var r0 []model.Reception
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, time.Time) ([]model.Reception, error)); ok {
		return rf(ctx, begin, end)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, time.Time) []model.Reception); ok {
		r0 = rf(ctx, begin, end)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Reception)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Time, time.Time) error); ok {
		r1 = rf(ctx, begin, end)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewReceptionRepository creates a new instance of ReceptionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReceptionRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ReceptionRepository {
	mock := &ReceptionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
