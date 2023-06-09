// Code generated by mockery v2.16.0. DO NOT EDIT.

package identity_service

import (
	context "context"
	"github.com/a-novel/agora-backend/domains/user/storage/identity"
	"github.com/a-novel/agora-backend/models"

	mock "github.com/stretchr/testify/mock"
	time "time"

	uuid "github.com/google/uuid"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

type MockService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockService) EXPECT() *MockService_Expecter {
	return &MockService_Expecter{mock: &_m.Mock}
}

// Age provides a mock function with given fields: source, now
func (_m *MockService) Age(source *models.UserIdentity, now time.Time) int {
	ret := _m.Called(source, now)

	var r0 int
	if rf, ok := ret.Get(0).(func(*models.UserIdentity, time.Time) int); ok {
		r0 = rf(source, now)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockService_Age_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Age'
type MockService_Age_Call struct {
	*mock.Call
}

// Age is a helper method to define mock.On call
//   - source *UserIdentity
//   - now time.Time
func (_e *MockService_Expecter) Age(source interface{}, now interface{}) *MockService_Age_Call {
	return &MockService_Age_Call{Call: _e.mock.On("Age", source, now)}
}

func (_c *MockService_Age_Call) Run(run func(source *models.UserIdentity, now time.Time)) *MockService_Age_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.UserIdentity), args[1].(time.Time))
	})
	return _c
}

func (_c *MockService_Age_Call) Return(_a0 int) *MockService_Age_Call {
	_c.Call.Return(_a0)
	return _c
}

// PrepareRegistration provides a mock function with given fields: ctx, data, now
func (_m *MockService) PrepareRegistration(ctx context.Context, data *models.UserIdentityUpdateForm, now time.Time) (*models.UserIdentityRegistrationForm, error) {
	ret := _m.Called(ctx, data, now)

	var r0 *models.UserIdentityRegistrationForm
	if rf, ok := ret.Get(0).(func(context.Context, *models.UserIdentityUpdateForm, time.Time) *models.UserIdentityRegistrationForm); ok {
		r0 = rf(ctx, data, now)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserIdentityRegistrationForm)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.UserIdentityUpdateForm, time.Time) error); ok {
		r1 = rf(ctx, data, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_PrepareRegistration_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PrepareRegistration'
type MockService_PrepareRegistration_Call struct {
	*mock.Call
}

// PrepareRegistration is a helper method to define mock.On call
//   - ctx context.Context
//   - data *UserIdentityUpdateForm
//   - now time.Time
func (_e *MockService_Expecter) PrepareRegistration(ctx interface{}, data interface{}, now interface{}) *MockService_PrepareRegistration_Call {
	return &MockService_PrepareRegistration_Call{Call: _e.mock.On("PrepareRegistration", ctx, data, now)}
}

func (_c *MockService_PrepareRegistration_Call) Run(run func(ctx context.Context, data *models.UserIdentityUpdateForm, now time.Time)) *MockService_PrepareRegistration_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.UserIdentityUpdateForm), args[2].(time.Time))
	})
	return _c
}

func (_c *MockService_PrepareRegistration_Call) Return(_a0 *models.UserIdentityRegistrationForm, _a1 error) *MockService_PrepareRegistration_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Read provides a mock function with given fields: ctx, id
func (_m *MockService) Read(ctx context.Context, id uuid.UUID) (*models.UserIdentity, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.UserIdentity
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *models.UserIdentity); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserIdentity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_Read_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Read'
type MockService_Read_Call struct {
	*mock.Call
}

// Read is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *MockService_Expecter) Read(ctx interface{}, id interface{}) *MockService_Read_Call {
	return &MockService_Read_Call{Call: _e.mock.On("Read", ctx, id)}
}

func (_c *MockService_Read_Call) Run(run func(ctx context.Context, id uuid.UUID)) *MockService_Read_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockService_Read_Call) Return(_a0 *models.UserIdentity, _a1 error) *MockService_Read_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// StorageToModel provides a mock function with given fields: source
func (_m *MockService) StorageToModel(source *identity_storage.Model) *models.UserIdentity {
	ret := _m.Called(source)

	var r0 *models.UserIdentity
	if rf, ok := ret.Get(0).(func(*identity_storage.Model) *models.UserIdentity); ok {
		r0 = rf(source)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserIdentity)
		}
	}

	return r0
}

// MockService_StorageToModel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StorageToModel'
type MockService_StorageToModel_Call struct {
	*mock.Call
}

// StorageToModel is a helper method to define mock.On call
//   - source *identity_storage.Model
func (_e *MockService_Expecter) StorageToModel(source interface{}) *MockService_StorageToModel_Call {
	return &MockService_StorageToModel_Call{Call: _e.mock.On("StorageToModel", source)}
}

func (_c *MockService_StorageToModel_Call) Run(run func(source *identity_storage.Model)) *MockService_StorageToModel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*identity_storage.Model))
	})
	return _c
}

func (_c *MockService_StorageToModel_Call) Return(_a0 *models.UserIdentity) *MockService_StorageToModel_Call {
	_c.Call.Return(_a0)
	return _c
}

// Update provides a mock function with given fields: ctx, data, id, now
func (_m *MockService) Update(ctx context.Context, data *models.UserIdentityUpdateForm, id uuid.UUID, now time.Time) (*models.UserIdentity, error) {
	ret := _m.Called(ctx, data, id, now)

	var r0 *models.UserIdentity
	if rf, ok := ret.Get(0).(func(context.Context, *models.UserIdentityUpdateForm, uuid.UUID, time.Time) *models.UserIdentity); ok {
		r0 = rf(ctx, data, id, now)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserIdentity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.UserIdentityUpdateForm, uuid.UUID, time.Time) error); ok {
		r1 = rf(ctx, data, id, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockService_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - data *UserIdentityUpdateForm
//   - id uuid.UUID
//   - now time.Time
func (_e *MockService_Expecter) Update(ctx interface{}, data interface{}, id interface{}, now interface{}) *MockService_Update_Call {
	return &MockService_Update_Call{Call: _e.mock.On("Update", ctx, data, id, now)}
}

func (_c *MockService_Update_Call) Run(run func(ctx context.Context, data *models.UserIdentityUpdateForm, id uuid.UUID, now time.Time)) *MockService_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.UserIdentityUpdateForm), args[2].(uuid.UUID), args[3].(time.Time))
	})
	return _c
}

func (_c *MockService_Update_Call) Return(_a0 *models.UserIdentity, _a1 error) *MockService_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewMockService interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockService creates a new instance of MockService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockService(t mockConstructorTestingTNewMockService) *MockService {
	mock := &MockService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
