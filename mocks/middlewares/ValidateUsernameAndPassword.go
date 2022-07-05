// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ValidateUsernameAndPassword is an autogenerated mock type for the ValidateUsernameAndPassword type
type ValidateUsernameAndPassword struct {
	mock.Mock
}

// Execute provides a mock function with given fields: username, password
func (_m *ValidateUsernameAndPassword) Execute(username string, password string) bool {
	ret := _m.Called(username, password)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(username, password)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewValidateUsernameAndPassword interface {
	mock.TestingT
	Cleanup(func())
}

// NewValidateUsernameAndPassword creates a new instance of ValidateUsernameAndPassword. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewValidateUsernameAndPassword(t mockConstructorTestingTNewValidateUsernameAndPassword) *ValidateUsernameAndPassword {
	mock := &ValidateUsernameAndPassword{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}