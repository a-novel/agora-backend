// Code generated by mockery v2.20.0. DO NOT EDIT.

package votes_service

import (
	context "context"
	"github.com/a-novel/agora-backend/domains/forum/storage/votes"
	"github.com/a-novel/agora-backend/models"
	time "time"

	mock "github.com/stretchr/testify/mock"

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

// GetVotedPosts provides a mock function with given fields: ctx, userID, target, limit, offset
func (_m *MockService) GetVotedPosts(ctx context.Context, userID uuid.UUID, target models.VoteTarget, limit int, offset int) ([]*models.VotedPost, int64, error) {
	ret := _m.Called(ctx, userID, target, limit, offset)

	var r0 []*models.VotedPost
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.VoteTarget, int, int) ([]*models.VotedPost, int64, error)); ok {
		return rf(ctx, userID, target, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.VoteTarget, int, int) []*models.VotedPost); ok {
		r0 = rf(ctx, userID, target, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.VotedPost)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, models.VoteTarget, int, int) int64); ok {
		r1 = rf(ctx, userID, target, limit, offset)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, uuid.UUID, models.VoteTarget, int, int) error); ok {
		r2 = rf(ctx, userID, target, limit, offset)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockService_GetVotedPosts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetVotedPosts'
type MockService_GetVotedPosts_Call struct {
	*mock.Call
}

// GetVotedPosts is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - target VoteTarget
//   - limit int
//   - offset int
func (_e *MockService_Expecter) GetVotedPosts(ctx interface{}, userID interface{}, target interface{}, limit interface{}, offset interface{}) *MockService_GetVotedPosts_Call {
	return &MockService_GetVotedPosts_Call{Call: _e.mock.On("GetVotedPosts", ctx, userID, target, limit, offset)}
}

func (_c *MockService_GetVotedPosts_Call) Run(run func(ctx context.Context, userID uuid.UUID, target models.VoteTarget, limit int, offset int)) *MockService_GetVotedPosts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(models.VoteTarget), args[3].(int), args[4].(int))
	})
	return _c
}

func (_c *MockService_GetVotedPosts_Call) Return(_a0 []*models.VotedPost, _a1 int64, _a2 error) *MockService_GetVotedPosts_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockService_GetVotedPosts_Call) RunAndReturn(run func(context.Context, uuid.UUID, models.VoteTarget, int, int) ([]*models.VotedPost, int64, error)) *MockService_GetVotedPosts_Call {
	_c.Call.Return(run)
	return _c
}

// HasVoted provides a mock function with given fields: ctx, postID, userID, target
func (_m *MockService) HasVoted(ctx context.Context, postID uuid.UUID, userID uuid.UUID, target models.VoteTarget) (models.VoteValue, error) {
	ret := _m.Called(ctx, postID, userID, target)

	var r0 models.VoteValue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, models.VoteTarget) (models.VoteValue, error)); ok {
		return rf(ctx, postID, userID, target)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, models.VoteTarget) models.VoteValue); ok {
		r0 = rf(ctx, postID, userID, target)
	} else {
		r0 = ret.Get(0).(models.VoteValue)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID, models.VoteTarget) error); ok {
		r1 = rf(ctx, postID, userID, target)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_HasVoted_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HasVoted'
type MockService_HasVoted_Call struct {
	*mock.Call
}

// HasVoted is a helper method to define mock.On call
//   - ctx context.Context
//   - postID uuid.UUID
//   - userID uuid.UUID
//   - target VoteTarget
func (_e *MockService_Expecter) HasVoted(ctx interface{}, postID interface{}, userID interface{}, target interface{}) *MockService_HasVoted_Call {
	return &MockService_HasVoted_Call{Call: _e.mock.On("HasVoted", ctx, postID, userID, target)}
}

func (_c *MockService_HasVoted_Call) Run(run func(ctx context.Context, postID uuid.UUID, userID uuid.UUID, target models.VoteTarget)) *MockService_HasVoted_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID), args[3].(models.VoteTarget))
	})
	return _c
}

func (_c *MockService_HasVoted_Call) Return(_a0 models.VoteValue, _a1 error) *MockService_HasVoted_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockService_HasVoted_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID, models.VoteTarget) (models.VoteValue, error)) *MockService_HasVoted_Call {
	_c.Call.Return(run)
	return _c
}

// StorageToModel provides a mock function with given fields: source
func (_m *MockService) StorageToModel(source *votes_storage.Model) *models.Vote {
	ret := _m.Called(source)

	var r0 *models.Vote
	if rf, ok := ret.Get(0).(func(*votes_storage.Model) *models.Vote); ok {
		r0 = rf(source)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Vote)
		}
	}

	return r0
}

// MockService_StorageToModel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StorageToModel'
type MockService_StorageToModel_Call struct {
	*mock.Call
}

// StorageToModel is a helper method to define mock.On call
//   - source *votes_storage.Model
func (_e *MockService_Expecter) StorageToModel(source interface{}) *MockService_StorageToModel_Call {
	return &MockService_StorageToModel_Call{Call: _e.mock.On("StorageToModel", source)}
}

func (_c *MockService_StorageToModel_Call) Run(run func(source *votes_storage.Model)) *MockService_StorageToModel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*votes_storage.Model))
	})
	return _c
}

func (_c *MockService_StorageToModel_Call) Return(_a0 *models.Vote) *MockService_StorageToModel_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockService_StorageToModel_Call) RunAndReturn(run func(*votes_storage.Model) *models.Vote) *MockService_StorageToModel_Call {
	_c.Call.Return(run)
	return _c
}

// Vote provides a mock function with given fields: ctx, postID, userID, target, vote, now
func (_m *MockService) Vote(ctx context.Context, postID uuid.UUID, userID uuid.UUID, target models.VoteTarget, vote models.VoteValue, now time.Time) (models.VoteValue, error) {
	ret := _m.Called(ctx, postID, userID, target, vote, now)

	var r0 models.VoteValue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, models.VoteTarget, models.VoteValue, time.Time) (models.VoteValue, error)); ok {
		return rf(ctx, postID, userID, target, vote, now)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, models.VoteTarget, models.VoteValue, time.Time) models.VoteValue); ok {
		r0 = rf(ctx, postID, userID, target, vote, now)
	} else {
		r0 = ret.Get(0).(models.VoteValue)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID, models.VoteTarget, models.VoteValue, time.Time) error); ok {
		r1 = rf(ctx, postID, userID, target, vote, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_Vote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VoteValue'
type MockService_Vote_Call struct {
	*mock.Call
}

// Vote is a helper method to define mock.On call
//   - ctx context.Context
//   - postID uuid.UUID
//   - userID uuid.UUID
//   - target VoteTarget
//   - vote Vote
//   - now time.Time
func (_e *MockService_Expecter) Vote(ctx interface{}, postID interface{}, userID interface{}, target interface{}, vote interface{}, now interface{}) *MockService_Vote_Call {
	return &MockService_Vote_Call{Call: _e.mock.On("VoteValue", ctx, postID, userID, target, vote, now)}
}

func (_c *MockService_Vote_Call) Run(run func(ctx context.Context, postID uuid.UUID, userID uuid.UUID, target models.VoteTarget, vote models.VoteValue, now time.Time)) *MockService_Vote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID), args[3].(models.VoteTarget), args[4].(models.VoteValue), args[5].(time.Time))
	})
	return _c
}

func (_c *MockService_Vote_Call) Return(_a0 models.VoteValue, _a1 error) *MockService_Vote_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockService_Vote_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID, models.VoteTarget, models.VoteValue, time.Time) (models.VoteValue, error)) *MockService_Vote_Call {
	_c.Call.Return(run)
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