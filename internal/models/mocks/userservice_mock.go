// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	models "github.com/VladimirStepanov/todo-app/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// ConfirmEmail provides a mock function with given fields: Link
func (_m *UserService) ConfirmEmail(Link string) error {
	ret := _m.Called(Link)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(Link)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: Email, Password
func (_m *UserService) Create(Email string, Password string) (*models.User, error) {
	ret := _m.Called(Email, Password)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string, string) *models.User); ok {
		r0 = rf(Email, Password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(Email, Password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignIn provides a mock function with given fields: Email, Password
func (_m *UserService) SignIn(Email string, Password string) (*models.User, error) {
	ret := _m.Called(Email, Password)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string, string) *models.User); ok {
		r0 = rf(Email, Password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(Email, Password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
