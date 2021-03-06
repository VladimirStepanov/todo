// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	models "github.com/VladimirStepanov/todo-app/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// ListService is an autogenerated mock type for the ListService type
type ListService struct {
	mock.Mock
}

// Create provides a mock function with given fields: title, description, userID
func (_m *ListService) Create(title string, description string, userID int64) (int64, error) {
	ret := _m.Called(title, description, userID)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, string, int64) int64); ok {
		r0 = rf(title, description, userID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, int64) error); ok {
		r1 = rf(title, description, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: listID
func (_m *ListService) Delete(listID int64) error {
	ret := _m.Called(listID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(listID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EditRole provides a mock function with given fields: listID, userID, role
func (_m *ListService) EditRole(listID int64, userID int64, role bool) error {
	ret := _m.Called(listID, userID, role)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, int64, bool) error); ok {
		r0 = rf(listID, userID, role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetListByID provides a mock function with given fields: listID, userID
func (_m *ListService) GetListByID(listID int64, userID int64) (*models.List, error) {
	ret := _m.Called(listID, userID)

	var r0 *models.List
	if rf, ok := ret.Get(0).(func(int64, int64) *models.List); ok {
		r0 = rf(listID, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.List)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, int64) error); ok {
		r1 = rf(listID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserLists provides a mock function with given fields: userID
func (_m *ListService) GetUserLists(userID int64) ([]*models.List, error) {
	ret := _m.Called(userID)

	var r0 []*models.List
	if rf, ok := ret.Get(0).(func(int64) []*models.List); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.List)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsListAdmin provides a mock function with given fields: ListID, userID
func (_m *ListService) IsListAdmin(ListID int64, userID int64) error {
	ret := _m.Called(ListID, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, int64) error); ok {
		r0 = rf(ListID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: listID, list
func (_m *ListService) Update(listID int64, list *models.UpdateListReq) error {
	ret := _m.Called(listID, list)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, *models.UpdateListReq) error); ok {
		r0 = rf(listID, list)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
