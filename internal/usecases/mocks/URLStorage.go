// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	entity "github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	mock "github.com/stretchr/testify/mock"
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

// AddURL provides a mock function with given fields: _a0
func (_m *URLStorage) AddURL(_a0 *entity.ShortURL) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for AddURL")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*entity.ShortURL) error); ok {
		r0 = rf(_a0)
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
//   - _a0 *entity.ShortURL
func (_e *URLStorage_Expecter) AddURL(_a0 interface{}) *URLStorage_AddURL_Call {
	return &URLStorage_AddURL_Call{Call: _e.mock.On("AddURL", _a0)}
}

func (_c *URLStorage_AddURL_Call) Run(run func(_a0 *entity.ShortURL)) *URLStorage_AddURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*entity.ShortURL))
	})
	return _c
}

func (_c *URLStorage_AddURL_Call) Return(_a0 error) *URLStorage_AddURL_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *URLStorage_AddURL_Call) RunAndReturn(run func(*entity.ShortURL) error) *URLStorage_AddURL_Call {
	_c.Call.Return(run)
	return _c
}

// GetURLByID provides a mock function with given fields: _a0
func (_m *URLStorage) GetURLByID(_a0 string) (*entity.ShortURL, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetURLByID")
	}

	var r0 *entity.ShortURL
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*entity.ShortURL, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) *entity.ShortURL); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.ShortURL)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
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
//   - _a0 string
func (_e *URLStorage_Expecter) GetURLByID(_a0 interface{}) *URLStorage_GetURLByID_Call {
	return &URLStorage_GetURLByID_Call{Call: _e.mock.On("GetURLByID", _a0)}
}

func (_c *URLStorage_GetURLByID_Call) Run(run func(_a0 string)) *URLStorage_GetURLByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *URLStorage_GetURLByID_Call) Return(_a0 *entity.ShortURL, _a1 error) *URLStorage_GetURLByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *URLStorage_GetURLByID_Call) RunAndReturn(run func(string) (*entity.ShortURL, error)) *URLStorage_GetURLByID_Call {
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
