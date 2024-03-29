// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Parser is an autogenerated mock type for the Parser type
type Parser struct {
	mock.Mock
}

// Parse provides a mock function with given fields: dir
func (_m *Parser) Parse(dir string) []string {
	ret := _m.Called(dir)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(dir)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

type mockConstructorTestingTNewParser interface {
	mock.TestingT
	Cleanup(func())
}

// NewParser creates a new instance of Parser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewParser(t mockConstructorTestingTNewParser) *Parser {
	mock := &Parser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
