// Code generated by mockery v2.20.0. DO NOT EDIT.

package improve_post_storage

import (
	context "context"
	"github.com/a-novel/agora-backend/domains/bookmark/storage"

	mock "github.com/stretchr/testify/mock"

	time "time"

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

// Bookmark provides a mock function with given fields: ctx, userID, requestID, target, level, now
func (_m *MockRepository) Bookmark(ctx context.Context, userID uuid.UUID, requestID uuid.UUID, target BookmarkTarget, level bookmark_storage.Level, now time.Time) (*Model, error) {
	ret := _m.Called(ctx, userID, requestID, target, level, now)

	var r0 *Model
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, BookmarkTarget, bookmark_storage.Level, time.Time) (*Model, error)); ok {
		return rf(ctx, userID, requestID, target, level, now)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, BookmarkTarget, bookmark_storage.Level, time.Time) *Model); ok {
		r0 = rf(ctx, userID, requestID, target, level, now)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Model)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID, BookmarkTarget, bookmark_storage.Level, time.Time) error); ok {
		r1 = rf(ctx, userID, requestID, target, level, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_Bookmark_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Bookmark'
type MockRepository_Bookmark_Call struct {
	*mock.Call
}

// Bookmark is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - requestID uuid.UUID
//   - target BookmarkTarget
//   - level bookmark_storage.Level
//   - now time.Time
func (_e *MockRepository_Expecter) Bookmark(ctx interface{}, userID interface{}, requestID interface{}, target interface{}, level interface{}, now interface{}) *MockRepository_Bookmark_Call {
	return &MockRepository_Bookmark_Call{Call: _e.mock.On("Bookmark", ctx, userID, requestID, target, level, now)}
}

func (_c *MockRepository_Bookmark_Call) Run(run func(ctx context.Context, userID uuid.UUID, requestID uuid.UUID, target BookmarkTarget, level bookmark_storage.Level, now time.Time)) *MockRepository_Bookmark_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID), args[3].(BookmarkTarget), args[4].(bookmark_storage.Level), args[5].(time.Time))
	})
	return _c
}

func (_c *MockRepository_Bookmark_Call) Return(_a0 *Model, _a1 error) *MockRepository_Bookmark_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_Bookmark_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID, BookmarkTarget, bookmark_storage.Level, time.Time) (*Model, error)) *MockRepository_Bookmark_Call {
	_c.Call.Return(run)
	return _c
}

// IsBookmarked provides a mock function with given fields: ctx, userID, requestID, target
func (_m *MockRepository) IsBookmarked(ctx context.Context, userID uuid.UUID, requestID uuid.UUID, target BookmarkTarget) (*bookmark_storage.Level, error) {
	ret := _m.Called(ctx, userID, requestID, target)

	var r0 *bookmark_storage.Level
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, BookmarkTarget) (*bookmark_storage.Level, error)); ok {
		return rf(ctx, userID, requestID, target)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, BookmarkTarget) *bookmark_storage.Level); ok {
		r0 = rf(ctx, userID, requestID, target)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bookmark_storage.Level)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID, BookmarkTarget) error); ok {
		r1 = rf(ctx, userID, requestID, target)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_IsBookmarked_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsBookmarked'
type MockRepository_IsBookmarked_Call struct {
	*mock.Call
}

// IsBookmarked is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - requestID uuid.UUID
//   - target BookmarkTarget
func (_e *MockRepository_Expecter) IsBookmarked(ctx interface{}, userID interface{}, requestID interface{}, target interface{}) *MockRepository_IsBookmarked_Call {
	return &MockRepository_IsBookmarked_Call{Call: _e.mock.On("IsBookmarked", ctx, userID, requestID, target)}
}

func (_c *MockRepository_IsBookmarked_Call) Run(run func(ctx context.Context, userID uuid.UUID, requestID uuid.UUID, target BookmarkTarget)) *MockRepository_IsBookmarked_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID), args[3].(BookmarkTarget))
	})
	return _c
}

func (_c *MockRepository_IsBookmarked_Call) Return(_a0 *bookmark_storage.Level, _a1 error) *MockRepository_IsBookmarked_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_IsBookmarked_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID, BookmarkTarget) (*bookmark_storage.Level, error)) *MockRepository_IsBookmarked_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, userID, level, target, limit, offset
func (_m *MockRepository) List(ctx context.Context, userID uuid.UUID, level bookmark_storage.Level, target BookmarkTarget, limit int, offset int) ([]*Model, int64, error) {
	ret := _m.Called(ctx, userID, level, target, limit, offset)

	var r0 []*Model
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, bookmark_storage.Level, BookmarkTarget, int, int) ([]*Model, int64, error)); ok {
		return rf(ctx, userID, level, target, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, bookmark_storage.Level, BookmarkTarget, int, int) []*Model); ok {
		r0 = rf(ctx, userID, level, target, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Model)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, bookmark_storage.Level, BookmarkTarget, int, int) int64); ok {
		r1 = rf(ctx, userID, level, target, limit, offset)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, uuid.UUID, bookmark_storage.Level, BookmarkTarget, int, int) error); ok {
		r2 = rf(ctx, userID, level, target, limit, offset)
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
//   - userID uuid.UUID
//   - level bookmark_storage.Level
//   - target BookmarkTarget
//   - limit int
//   - offset int
func (_e *MockRepository_Expecter) List(ctx interface{}, userID interface{}, level interface{}, target interface{}, limit interface{}, offset interface{}) *MockRepository_List_Call {
	return &MockRepository_List_Call{Call: _e.mock.On("List", ctx, userID, level, target, limit, offset)}
}

func (_c *MockRepository_List_Call) Run(run func(ctx context.Context, userID uuid.UUID, level bookmark_storage.Level, target BookmarkTarget, limit int, offset int)) *MockRepository_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(bookmark_storage.Level), args[3].(BookmarkTarget), args[4].(int), args[5].(int))
	})
	return _c
}

func (_c *MockRepository_List_Call) Return(_a0 []*Model, _a1 int64, _a2 error) *MockRepository_List_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockRepository_List_Call) RunAndReturn(run func(context.Context, uuid.UUID, bookmark_storage.Level, BookmarkTarget, int, int) ([]*Model, int64, error)) *MockRepository_List_Call {
	_c.Call.Return(run)
	return _c
}

// UnBookmark provides a mock function with given fields: ctx, userID, requestID, target
func (_m *MockRepository) UnBookmark(ctx context.Context, userID uuid.UUID, requestID uuid.UUID, target BookmarkTarget) error {
	ret := _m.Called(ctx, userID, requestID, target)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, BookmarkTarget) error); ok {
		r0 = rf(ctx, userID, requestID, target)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRepository_UnBookmark_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnBookmark'
type MockRepository_UnBookmark_Call struct {
	*mock.Call
}

// UnBookmark is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - requestID uuid.UUID
//   - target BookmarkTarget
func (_e *MockRepository_Expecter) UnBookmark(ctx interface{}, userID interface{}, requestID interface{}, target interface{}) *MockRepository_UnBookmark_Call {
	return &MockRepository_UnBookmark_Call{Call: _e.mock.On("UnBookmark", ctx, userID, requestID, target)}
}

func (_c *MockRepository_UnBookmark_Call) Run(run func(ctx context.Context, userID uuid.UUID, requestID uuid.UUID, target BookmarkTarget)) *MockRepository_UnBookmark_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID), args[3].(BookmarkTarget))
	})
	return _c
}

func (_c *MockRepository_UnBookmark_Call) Return(_a0 error) *MockRepository_UnBookmark_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRepository_UnBookmark_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID, BookmarkTarget) error) *MockRepository_UnBookmark_Call {
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
