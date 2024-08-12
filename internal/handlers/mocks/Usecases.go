// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	entity "github.com/DyadyaRodya/go-shortener/internal/domain/entity"

	mock "github.com/stretchr/testify/mock"
)

// Usecases is an autogenerated mock type for the Usecases type
type Usecases struct {
	mock.Mock
}

type Usecases_Expecter struct {
	mock *mock.Mock
}

func (_m *Usecases) EXPECT() *Usecases_Expecter {
	return &Usecases_Expecter{mock: &_m.Mock}
}

// CreateShortURL provides a mock function with given fields: URL
func (_m *Usecases) CreateShortURL(URL string) (*entity.ShortURL, error) {
	ret := _m.Called(URL)

	if len(ret) == 0 {
		panic("no return value specified for CreateShortURL")
	}

	var r0 *entity.ShortURL
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*entity.ShortURL, error)); ok {
		return rf(URL)
	}
	if rf, ok := ret.Get(0).(func(string) *entity.ShortURL); ok {
		r0 = rf(URL)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.ShortURL)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(URL)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Usecases_CreateShortURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateShortURL'
type Usecases_CreateShortURL_Call struct {
	*mock.Call
}

// CreateShortURL is a helper method to define mock.On call
//   - URL string
func (_e *Usecases_Expecter) CreateShortURL(URL interface{}) *Usecases_CreateShortURL_Call {
	return &Usecases_CreateShortURL_Call{Call: _e.mock.On("CreateShortURL", URL)}
}

func (_c *Usecases_CreateShortURL_Call) Run(run func(URL string)) *Usecases_CreateShortURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Usecases_CreateShortURL_Call) Return(_a0 *entity.ShortURL, _a1 error) *Usecases_CreateShortURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Usecases_CreateShortURL_Call) RunAndReturn(run func(string) (*entity.ShortURL, error)) *Usecases_CreateShortURL_Call {
	_c.Call.Return(run)
	return _c
}

// GetShortURL provides a mock function with given fields: ID
func (_m *Usecases) GetShortURL(ID string) (*entity.ShortURL, error) {
	ret := _m.Called(ID)

	if len(ret) == 0 {
		panic("no return value specified for GetShortURL")
	}

	var r0 *entity.ShortURL
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*entity.ShortURL, error)); ok {
		return rf(ID)
	}
	if rf, ok := ret.Get(0).(func(string) *entity.ShortURL); ok {
		r0 = rf(ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.ShortURL)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Usecases_GetShortURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetShortURL'
type Usecases_GetShortURL_Call struct {
	*mock.Call
}

// GetShortURL is a helper method to define mock.On call
//   - ID string
func (_e *Usecases_Expecter) GetShortURL(ID interface{}) *Usecases_GetShortURL_Call {
	return &Usecases_GetShortURL_Call{Call: _e.mock.On("GetShortURL", ID)}
}

func (_c *Usecases_GetShortURL_Call) Run(run func(ID string)) *Usecases_GetShortURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Usecases_GetShortURL_Call) Return(_a0 *entity.ShortURL, _a1 error) *Usecases_GetShortURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Usecases_GetShortURL_Call) RunAndReturn(run func(string) (*entity.ShortURL, error)) *Usecases_GetShortURL_Call {
	_c.Call.Return(run)
	return _c
}

// NewUsecases creates a new instance of Usecases. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUsecases(t interface {
	mock.TestingT
	Cleanup(func())
}) *Usecases {
	mock := &Usecases{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
