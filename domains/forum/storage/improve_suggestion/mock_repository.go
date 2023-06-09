// Code generated by mockery v2.20.0. DO NOT EDIT.

package improve_suggestion_storage

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

type MockRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRepository) EXPECT() *MockRepository_Expecter {
	return &MockRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, data, userID, sourceID, id, now
func (_m *MockRepository) Create(ctx context.Context, data *Core, userID uuid.UUID, sourceID uuid.UUID, id uuid.UUID, now time.Time) (*Model, error) {
	ret := _m.Called(ctx, data, userID, sourceID, id, now)

	var r0 *Model
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *Core, uuid.UUID, uuid.UUID, uuid.UUID, time.Time) (*Model, error)); ok {
		return rf(ctx, data, userID, sourceID, id, now)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *Core, uuid.UUID, uuid.UUID, uuid.UUID, time.Time) *Model); ok {
		r0 = rf(ctx, data, userID, sourceID, id, now)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Model)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *Core, uuid.UUID, uuid.UUID, uuid.UUID, time.Time) error); ok {
		r1 = rf(ctx, data, userID, sourceID, id, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - data *Core
//   - userID uuid.UUID
//   - sourceID uuid.UUID
//   - id uuid.UUID
//   - now time.Time
func (_e *MockRepository_Expecter) Create(ctx interface{}, data interface{}, userID interface{}, sourceID interface{}, id interface{}, now interface{}) *MockRepository_Create_Call {
	return &MockRepository_Create_Call{Call: _e.mock.On("Create", ctx, data, userID, sourceID, id, now)}
}

func (_c *MockRepository_Create_Call) Run(run func(ctx context.Context, data *Core, userID uuid.UUID, sourceID uuid.UUID, id uuid.UUID, now time.Time)) *MockRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*Core), args[2].(uuid.UUID), args[3].(uuid.UUID), args[4].(uuid.UUID), args[5].(time.Time))
	})
	return _c
}

func (_c *MockRepository_Create_Call) Return(_a0 *Model, _a1 error) *MockRepository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_Create_Call) RunAndReturn(run func(context.Context, *Core, uuid.UUID, uuid.UUID, uuid.UUID, time.Time) (*Model, error)) *MockRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, id
func (_m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *MockRepository_Expecter) Delete(ctx interface{}, id interface{}) *MockRepository_Delete_Call {
	return &MockRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *MockRepository_Delete_Call) Run(run func(ctx context.Context, id uuid.UUID)) *MockRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_Delete_Call) Return(_a0 error) *MockRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRepository_Delete_Call) RunAndReturn(run func(context.Context, uuid.UUID) error) *MockRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetPreviews provides a mock function with given fields: ctx, ids
func (_m *MockRepository) GetPreviews(ctx context.Context, ids []uuid.UUID) ([]*Model, error) {
	ret := _m.Called(ctx, ids)

	var r0 []*Model
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []uuid.UUID) ([]*Model, error)); ok {
		return rf(ctx, ids)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []uuid.UUID) []*Model); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Model)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []uuid.UUID) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_GetPreviews_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPreviews'
type MockRepository_GetPreviews_Call struct {
	*mock.Call
}

// GetPreviews is a helper method to define mock.On call
//   - ctx context.Context
//   - ids []uuid.UUID
func (_e *MockRepository_Expecter) GetPreviews(ctx interface{}, ids interface{}) *MockRepository_GetPreviews_Call {
	return &MockRepository_GetPreviews_Call{Call: _e.mock.On("GetPreviews", ctx, ids)}
}

func (_c *MockRepository_GetPreviews_Call) Run(run func(ctx context.Context, ids []uuid.UUID)) *MockRepository_GetPreviews_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_GetPreviews_Call) Return(_a0 []*Model, _a1 error) *MockRepository_GetPreviews_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_GetPreviews_Call) RunAndReturn(run func(context.Context, []uuid.UUID) ([]*Model, error)) *MockRepository_GetPreviews_Call {
	_c.Call.Return(run)
	return _c
}

// IsCreator provides a mock function with given fields: ctx, userID, postID
func (_m *MockRepository) IsCreator(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (bool, error) {
	ret := _m.Called(ctx, userID, postID)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) (bool, error)); ok {
		return rf(ctx, userID, postID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) bool); ok {
		r0 = rf(ctx, userID, postID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, userID, postID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_IsCreator_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsCreator'
type MockRepository_IsCreator_Call struct {
	*mock.Call
}

// IsCreator is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - postID uuid.UUID
func (_e *MockRepository_Expecter) IsCreator(ctx interface{}, userID interface{}, postID interface{}) *MockRepository_IsCreator_Call {
	return &MockRepository_IsCreator_Call{Call: _e.mock.On("IsCreator", ctx, userID, postID)}
}

func (_c *MockRepository_IsCreator_Call) Run(run func(ctx context.Context, userID uuid.UUID, postID uuid.UUID)) *MockRepository_IsCreator_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_IsCreator_Call) Return(_a0 bool, _a1 error) *MockRepository_IsCreator_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_IsCreator_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) (bool, error)) *MockRepository_IsCreator_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, query, limit, offset
func (_m *MockRepository) List(ctx context.Context, query ListQuery, limit int, offset int) ([]*Model, int64, error) {
	ret := _m.Called(ctx, query, limit, offset)

	var r0 []*Model
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, ListQuery, int, int) ([]*Model, int64, error)); ok {
		return rf(ctx, query, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ListQuery, int, int) []*Model); ok {
		r0 = rf(ctx, query, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Model)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ListQuery, int, int) int64); ok {
		r1 = rf(ctx, query, limit, offset)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, ListQuery, int, int) error); ok {
		r2 = rf(ctx, query, limit, offset)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockRepository_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type MockRepository_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - query ListQuery
//   - limit int
//   - offset int
func (_e *MockRepository_Expecter) List(ctx interface{}, query interface{}, limit interface{}, offset interface{}) *MockRepository_List_Call {
	return &MockRepository_List_Call{Call: _e.mock.On("List", ctx, query, limit, offset)}
}

func (_c *MockRepository_List_Call) Run(run func(ctx context.Context, query ListQuery, limit int, offset int)) *MockRepository_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ListQuery), args[2].(int), args[3].(int))
	})
	return _c
}

func (_c *MockRepository_List_Call) Return(_a0 []*Model, _a1 int64, _a2 error) *MockRepository_List_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockRepository_List_Call) RunAndReturn(run func(context.Context, ListQuery, int, int) ([]*Model, int64, error)) *MockRepository_List_Call {
	_c.Call.Return(run)
	return _c
}

// Read provides a mock function with given fields: ctx, id
func (_m *MockRepository) Read(ctx context.Context, id uuid.UUID) (*Model, error) {
	ret := _m.Called(ctx, id)

	var r0 *Model
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*Model, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *Model); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Model)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_Read_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Read'
type MockRepository_Read_Call struct {
	*mock.Call
}

// Read is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *MockRepository_Expecter) Read(ctx interface{}, id interface{}) *MockRepository_Read_Call {
	return &MockRepository_Read_Call{Call: _e.mock.On("Read", ctx, id)}
}

func (_c *MockRepository_Read_Call) Run(run func(ctx context.Context, id uuid.UUID)) *MockRepository_Read_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_Read_Call) Return(_a0 *Model, _a1 error) *MockRepository_Read_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_Read_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*Model, error)) *MockRepository_Read_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, data, id, now
func (_m *MockRepository) Update(ctx context.Context, data *Core, id uuid.UUID, now time.Time) (*Model, error) {
	ret := _m.Called(ctx, data, id, now)

	var r0 *Model
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *Core, uuid.UUID, time.Time) (*Model, error)); ok {
		return rf(ctx, data, id, now)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *Core, uuid.UUID, time.Time) *Model); ok {
		r0 = rf(ctx, data, id, now)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Model)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *Core, uuid.UUID, time.Time) error); ok {
		r1 = rf(ctx, data, id, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - data *Core
//   - id uuid.UUID
//   - now time.Time
func (_e *MockRepository_Expecter) Update(ctx interface{}, data interface{}, id interface{}, now interface{}) *MockRepository_Update_Call {
	return &MockRepository_Update_Call{Call: _e.mock.On("Update", ctx, data, id, now)}
}

func (_c *MockRepository_Update_Call) Run(run func(ctx context.Context, data *Core, id uuid.UUID, now time.Time)) *MockRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*Core), args[2].(uuid.UUID), args[3].(time.Time))
	})
	return _c
}

func (_c *MockRepository_Update_Call) Return(_a0 *Model, _a1 error) *MockRepository_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_Update_Call) RunAndReturn(run func(context.Context, *Core, uuid.UUID, time.Time) (*Model, error)) *MockRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// Validate provides a mock function with given fields: ctx, validated, id
func (_m *MockRepository) Validate(ctx context.Context, validated bool, id uuid.UUID) (*Model, error) {
	ret := _m.Called(ctx, validated, id)

	var r0 *Model
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, bool, uuid.UUID) (*Model, error)); ok {
		return rf(ctx, validated, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, bool, uuid.UUID) *Model); ok {
		r0 = rf(ctx, validated, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Model)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, bool, uuid.UUID) error); ok {
		r1 = rf(ctx, validated, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_Validate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Validate'
type MockRepository_Validate_Call struct {
	*mock.Call
}

// Validate is a helper method to define mock.On call
//   - ctx context.Context
//   - validated bool
//   - id uuid.UUID
func (_e *MockRepository_Expecter) Validate(ctx interface{}, validated interface{}, id interface{}) *MockRepository_Validate_Call {
	return &MockRepository_Validate_Call{Call: _e.mock.On("Validate", ctx, validated, id)}
}

func (_c *MockRepository_Validate_Call) Run(run func(ctx context.Context, validated bool, id uuid.UUID)) *MockRepository_Validate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(bool), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_Validate_Call) Return(_a0 *Model, _a1 error) *MockRepository_Validate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_Validate_Call) RunAndReturn(run func(context.Context, bool, uuid.UUID) (*Model, error)) *MockRepository_Validate_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockRepository creates a new instance of MockRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockRepository(t mockConstructorTestingTNewMockRepository) *MockRepository {
	mock := &MockRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}