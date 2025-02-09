// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	dto "github.com/DyadyaRodya/go-shortener/internal/usecases/dto"

	mock "github.com/stretchr/testify/mock"

	usecases "github.com/DyadyaRodya/go-shortener/internal/usecases"
)

// URLStorage is an autogenerated mock type for the URLStorage type
type URLStorage struct {
	mock.Mock
}

type URLStorage_Expecter struct {
	mock *mock.Mock
}

func (_m *URLStorage) EXPECT() *URLStorage_Expecter {
	return &URLStorage_Expecter{mock: &_m.Mock}
}

// AddURL provides a mock function with given fields: ctx, ShortURL, OwnerUUID
func (_m *URLStorage) AddURL(ctx context.Context, ShortURL *entity.ShortURL, OwnerUUID string) error {
	ret := _m.Called(ctx, ShortURL, OwnerUUID)

	if len(ret) == 0 {
		panic("no return value specified for AddURL")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.ShortURL, string) error); ok {
		r0 = rf(ctx, ShortURL, OwnerUUID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// URLStorage_AddURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddURL'
type URLStorage_AddURL_Call struct {
	*mock.Call
}

// AddURL is a helper method to define mock.On call
//   - ctx context.Context
//   - ShortURL *entity.ShortURL
//   - OwnerUUID string
func (_e *URLStorage_Expecter) AddURL(ctx interface{}, ShortURL interface{}, OwnerUUID interface{}) *URLStorage_AddURL_Call {
	return &URLStorage_AddURL_Call{Call: _e.mock.On("AddURL", ctx, ShortURL, OwnerUUID)}
}

func (_c *URLStorage_AddURL_Call) Run(run func(ctx context.Context, ShortURL *entity.ShortURL, OwnerUUID string)) *URLStorage_AddURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.ShortURL), args[2].(string))
	})
	return _c
}

func (_c *URLStorage_AddURL_Call) Return(_a0 error) *URLStorage_AddURL_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *URLStorage_AddURL_Call) RunAndReturn(run func(context.Context, *entity.ShortURL, string) error) *URLStorage_AddURL_Call {
	_c.Call.Return(run)
	return _c
}

// Begin provides a mock function with given fields: ctx
func (_m *URLStorage) Begin(ctx context.Context) (usecases.Transaction, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Begin")
	}

	var r0 usecases.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (usecases.Transaction, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) usecases.Transaction); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(usecases.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// URLStorage_Begin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Begin'
type URLStorage_Begin_Call struct {
	*mock.Call
}

// Begin is a helper method to define mock.On call
//   - ctx context.Context
func (_e *URLStorage_Expecter) Begin(ctx interface{}) *URLStorage_Begin_Call {
	return &URLStorage_Begin_Call{Call: _e.mock.On("Begin", ctx)}
}

func (_c *URLStorage_Begin_Call) Run(run func(ctx context.Context)) *URLStorage_Begin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *URLStorage_Begin_Call) Return(_a0 usecases.Transaction, _a1 error) *URLStorage_Begin_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *URLStorage_Begin_Call) RunAndReturn(run func(context.Context) (usecases.Transaction, error)) *URLStorage_Begin_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteUserURLs provides a mock function with given fields: ctx, requests
func (_m *URLStorage) DeleteUserURLs(ctx context.Context, requests ...*dto.DeleteUserShortURLsRequest) error {
	_va := make([]interface{}, len(requests))
	for _i := range requests {
		_va[_i] = requests[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUserURLs")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...*dto.DeleteUserShortURLsRequest) error); ok {
		r0 = rf(ctx, requests...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// URLStorage_DeleteUserURLs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteUserURLs'
type URLStorage_DeleteUserURLs_Call struct {
	*mock.Call
}

// DeleteUserURLs is a helper method to define mock.On call
//   - ctx context.Context
//   - requests ...*dto.DeleteUserShortURLsRequest
func (_e *URLStorage_Expecter) DeleteUserURLs(ctx interface{}, requests ...interface{}) *URLStorage_DeleteUserURLs_Call {
	return &URLStorage_DeleteUserURLs_Call{Call: _e.mock.On("DeleteUserURLs",
		append([]interface{}{ctx}, requests...)...)}
}

func (_c *URLStorage_DeleteUserURLs_Call) Run(run func(ctx context.Context, requests ...*dto.DeleteUserShortURLsRequest)) *URLStorage_DeleteUserURLs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*dto.DeleteUserShortURLsRequest, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(*dto.DeleteUserShortURLsRequest)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *URLStorage_DeleteUserURLs_Call) Return(_a0 error) *URLStorage_DeleteUserURLs_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *URLStorage_DeleteUserURLs_Call) RunAndReturn(run func(context.Context, ...*dto.DeleteUserShortURLsRequest) error) *URLStorage_DeleteUserURLs_Call {
	_c.Call.Return(run)
	return _c
}

// GetStats provides a mock function with given fields: ctx
func (_m *URLStorage) GetStats(ctx context.Context) (*dto.StatsResponse, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetStats")
	}

	var r0 *dto.StatsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*dto.StatsResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *dto.StatsResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.StatsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// URLStorage_GetStats_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetStats'
type URLStorage_GetStats_Call struct {
	*mock.Call
}

// GetStats is a helper method to define mock.On call
//   - ctx context.Context
func (_e *URLStorage_Expecter) GetStats(ctx interface{}) *URLStorage_GetStats_Call {
	return &URLStorage_GetStats_Call{Call: _e.mock.On("GetStats", ctx)}
}

func (_c *URLStorage_GetStats_Call) Run(run func(ctx context.Context)) *URLStorage_GetStats_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *URLStorage_GetStats_Call) Return(_a0 *dto.StatsResponse, _a1 error) *URLStorage_GetStats_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *URLStorage_GetStats_Call) RunAndReturn(run func(context.Context) (*dto.StatsResponse, error)) *URLStorage_GetStats_Call {
	_c.Call.Return(run)
	return _c
}

// GetURLByID provides a mock function with given fields: ctx, ID
func (_m *URLStorage) GetURLByID(ctx context.Context, ID string) (*entity.ShortURL, error) {
	ret := _m.Called(ctx, ID)

	if len(ret) == 0 {
		panic("no return value specified for GetURLByID")
	}

	var r0 *entity.ShortURL
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*entity.ShortURL, error)); ok {
		return rf(ctx, ID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.ShortURL); ok {
		r0 = rf(ctx, ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.ShortURL)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// URLStorage_GetURLByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetURLByID'
type URLStorage_GetURLByID_Call struct {
	*mock.Call
}

// GetURLByID is a helper method to define mock.On call
//   - ctx context.Context
//   - ID string
func (_e *URLStorage_Expecter) GetURLByID(ctx interface{}, ID interface{}) *URLStorage_GetURLByID_Call {
	return &URLStorage_GetURLByID_Call{Call: _e.mock.On("GetURLByID", ctx, ID)}
}

func (_c *URLStorage_GetURLByID_Call) Run(run func(ctx context.Context, ID string)) *URLStorage_GetURLByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *URLStorage_GetURLByID_Call) Return(_a0 *entity.ShortURL, _a1 error) *URLStorage_GetURLByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *URLStorage_GetURLByID_Call) RunAndReturn(run func(context.Context, string) (*entity.ShortURL, error)) *URLStorage_GetURLByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserUrls provides a mock function with given fields: ctx, UserUUID
func (_m *URLStorage) GetUserUrls(ctx context.Context, UserUUID string) (map[string]*entity.ShortURL, error) {
	ret := _m.Called(ctx, UserUUID)

	if len(ret) == 0 {
		panic("no return value specified for GetUserUrls")
	}

	var r0 map[string]*entity.ShortURL
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (map[string]*entity.ShortURL, error)); ok {
		return rf(ctx, UserUUID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) map[string]*entity.ShortURL); ok {
		r0 = rf(ctx, UserUUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*entity.ShortURL)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, UserUUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// URLStorage_GetUserUrls_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserUrls'
type URLStorage_GetUserUrls_Call struct {
	*mock.Call
}

// GetUserUrls is a helper method to define mock.On call
//   - ctx context.Context
//   - UserUUID string
func (_e *URLStorage_Expecter) GetUserUrls(ctx interface{}, UserUUID interface{}) *URLStorage_GetUserUrls_Call {
	return &URLStorage_GetUserUrls_Call{Call: _e.mock.On("GetUserUrls", ctx, UserUUID)}
}

func (_c *URLStorage_GetUserUrls_Call) Run(run func(ctx context.Context, UserUUID string)) *URLStorage_GetUserUrls_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *URLStorage_GetUserUrls_Call) Return(_a0 map[string]*entity.ShortURL, _a1 error) *URLStorage_GetUserUrls_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *URLStorage_GetUserUrls_Call) RunAndReturn(run func(context.Context, string) (map[string]*entity.ShortURL, error)) *URLStorage_GetUserUrls_Call {
	_c.Call.Return(run)
	return _c
}

// TestConnection provides a mock function with given fields: ctx
func (_m *URLStorage) TestConnection(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for TestConnection")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// URLStorage_TestConnection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TestConnection'
type URLStorage_TestConnection_Call struct {
	*mock.Call
}

// TestConnection is a helper method to define mock.On call
//   - ctx context.Context
func (_e *URLStorage_Expecter) TestConnection(ctx interface{}) *URLStorage_TestConnection_Call {
	return &URLStorage_TestConnection_Call{Call: _e.mock.On("TestConnection", ctx)}
}

func (_c *URLStorage_TestConnection_Call) Run(run func(ctx context.Context)) *URLStorage_TestConnection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *URLStorage_TestConnection_Call) Return(_a0 error) *URLStorage_TestConnection_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *URLStorage_TestConnection_Call) RunAndReturn(run func(context.Context) error) *URLStorage_TestConnection_Call {
	_c.Call.Return(run)
	return _c
}

// NewURLStorage creates a new instance of URLStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewURLStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *URLStorage {
	mock := &URLStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
