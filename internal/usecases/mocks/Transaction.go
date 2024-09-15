// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	mock "github.com/stretchr/testify/mock"
)

// Transaction is an autogenerated mock type for the Transaction type
type Transaction struct {
	mock.Mock
}

type Transaction_Expecter struct {
	mock *mock.Mock
}

func (_m *Transaction) EXPECT() *Transaction_Expecter {
	return &Transaction_Expecter{mock: &_m.Mock}
}

// AddURL provides a mock function with given fields: ctx, ShortURL
func (_m *Transaction) AddURL(ctx context.Context, ShortURL *entity.ShortURL) error {
	ret := _m.Called(ctx, ShortURL)

	if len(ret) == 0 {
		panic("no return value specified for AddURL")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.ShortURL) error); ok {
		r0 = rf(ctx, ShortURL)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Transaction_AddURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddURL'
type Transaction_AddURL_Call struct {
	*mock.Call
}

// AddURL is a helper method to define mock.On call
//   - ctx context.Context
//   - ShortURL *entity.ShortURL
func (_e *Transaction_Expecter) AddURL(ctx interface{}, ShortURL interface{}) *Transaction_AddURL_Call {
	return &Transaction_AddURL_Call{Call: _e.mock.On("AddURL", ctx, ShortURL)}
}

func (_c *Transaction_AddURL_Call) Run(run func(ctx context.Context, ShortURL *entity.ShortURL)) *Transaction_AddURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.ShortURL))
	})
	return _c
}

func (_c *Transaction_AddURL_Call) Return(_a0 error) *Transaction_AddURL_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Transaction_AddURL_Call) RunAndReturn(run func(context.Context, *entity.ShortURL) error) *Transaction_AddURL_Call {
	_c.Call.Return(run)
	return _c
}

// CheckIDs provides a mock function with given fields: ctx, IDs
func (_m *Transaction) CheckIDs(ctx context.Context, IDs []string) ([]string, error) {
	ret := _m.Called(ctx, IDs)

	if len(ret) == 0 {
		panic("no return value specified for CheckIDs")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) ([]string, error)); ok {
		return rf(ctx, IDs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string) []string); ok {
		r0 = rf(ctx, IDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, IDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Transaction_CheckIDs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckIDs'
type Transaction_CheckIDs_Call struct {
	*mock.Call
}

// CheckIDs is a helper method to define mock.On call
//   - ctx context.Context
//   - IDs []string
func (_e *Transaction_Expecter) CheckIDs(ctx interface{}, IDs interface{}) *Transaction_CheckIDs_Call {
	return &Transaction_CheckIDs_Call{Call: _e.mock.On("CheckIDs", ctx, IDs)}
}

func (_c *Transaction_CheckIDs_Call) Run(run func(ctx context.Context, IDs []string)) *Transaction_CheckIDs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *Transaction_CheckIDs_Call) Return(_a0 []string, _a1 error) *Transaction_CheckIDs_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Transaction_CheckIDs_Call) RunAndReturn(run func(context.Context, []string) ([]string, error)) *Transaction_CheckIDs_Call {
	_c.Call.Return(run)
	return _c
}

// Commit provides a mock function with given fields: ctx
func (_m *Transaction) Commit(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Commit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Transaction_Commit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Commit'
type Transaction_Commit_Call struct {
	*mock.Call
}

// Commit is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Transaction_Expecter) Commit(ctx interface{}) *Transaction_Commit_Call {
	return &Transaction_Commit_Call{Call: _e.mock.On("Commit", ctx)}
}

func (_c *Transaction_Commit_Call) Run(run func(ctx context.Context)) *Transaction_Commit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Transaction_Commit_Call) Return(_a0 error) *Transaction_Commit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Transaction_Commit_Call) RunAndReturn(run func(context.Context) error) *Transaction_Commit_Call {
	_c.Call.Return(run)
	return _c
}

// GetByURLs provides a mock function with given fields: ctx, URLs
func (_m *Transaction) GetByURLs(ctx context.Context, URLs []string) (map[string]*entity.ShortURL, error) {
	ret := _m.Called(ctx, URLs)

	if len(ret) == 0 {
		panic("no return value specified for GetByURLs")
	}

	var r0 map[string]*entity.ShortURL
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) (map[string]*entity.ShortURL, error)); ok {
		return rf(ctx, URLs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string) map[string]*entity.ShortURL); ok {
		r0 = rf(ctx, URLs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*entity.ShortURL)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, URLs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Transaction_GetByURLs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByURLs'
type Transaction_GetByURLs_Call struct {
	*mock.Call
}

// GetByURLs is a helper method to define mock.On call
//   - ctx context.Context
//   - URLs []string
func (_e *Transaction_Expecter) GetByURLs(ctx interface{}, URLs interface{}) *Transaction_GetByURLs_Call {
	return &Transaction_GetByURLs_Call{Call: _e.mock.On("GetByURLs", ctx, URLs)}
}

func (_c *Transaction_GetByURLs_Call) Run(run func(ctx context.Context, URLs []string)) *Transaction_GetByURLs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *Transaction_GetByURLs_Call) Return(_a0 map[string]*entity.ShortURL, _a1 error) *Transaction_GetByURLs_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Transaction_GetByURLs_Call) RunAndReturn(run func(context.Context, []string) (map[string]*entity.ShortURL, error)) *Transaction_GetByURLs_Call {
	_c.Call.Return(run)
	return _c
}

// Rollback provides a mock function with given fields: ctx
func (_m *Transaction) Rollback(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Rollback")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Transaction_Rollback_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Rollback'
type Transaction_Rollback_Call struct {
	*mock.Call
}

// Rollback is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Transaction_Expecter) Rollback(ctx interface{}) *Transaction_Rollback_Call {
	return &Transaction_Rollback_Call{Call: _e.mock.On("Rollback", ctx)}
}

func (_c *Transaction_Rollback_Call) Run(run func(ctx context.Context)) *Transaction_Rollback_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Transaction_Rollback_Call) Return(_a0 error) *Transaction_Rollback_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Transaction_Rollback_Call) RunAndReturn(run func(context.Context) error) *Transaction_Rollback_Call {
	_c.Call.Return(run)
	return _c
}

// NewTransaction creates a new instance of Transaction. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransaction(t interface {
	mock.TestingT
	Cleanup(func())
}) *Transaction {
	mock := &Transaction{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
