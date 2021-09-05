// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	models "github.com/VladimirStepanov/todo-app/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// ItemService is an autogenerated mock type for the ItemService type
type ItemService struct {
	mock.Mock
}

// Create provides a mock function with given fields: title, description, listID
func (_m *ItemService) Create(title string, description string, listID int64) (int64, error) {
	ret := _m.Called(title, description, listID)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, string, int64) int64); ok {
		r0 = rf(title, description, listID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, int64) error); ok {
		r1 = rf(title, description, listID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: listID, itemID
func (_m *ItemService) Delete(listID int64, itemID int64) error {
	ret := _m.Called(listID, itemID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, int64) error); ok {
		r0 = rf(listID, itemID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Done provides a mock function with given fields: listID, itemID
func (_m *ItemService) Done(listID int64, itemID int64) error {
	ret := _m.Called(listID, itemID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, int64) error); ok {
		r0 = rf(listID, itemID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetItemBydID provides a mock function with given fields: listID, itemID
func (_m *ItemService) GetItemBydID(listID int64, itemID int64) (*models.Item, error) {
	ret := _m.Called(listID, itemID)

	var r0 *models.Item
	if rf, ok := ret.Get(0).(func(int64, int64) *models.Item); ok {
		r0 = rf(listID, itemID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Item)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, int64) error); ok {
		r1 = rf(listID, itemID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetItems provides a mock function with given fields: listID
func (_m *ItemService) GetItems(listID int64) ([]*models.Item, error) {
	ret := _m.Called(listID)

	var r0 []*models.Item
	if rf, ok := ret.Get(0).(func(int64) []*models.Item); ok {
		r0 = rf(listID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Item)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(listID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: listID, itemID, item
func (_m *ItemService) Update(listID int64, itemID int64, item *models.UpdateItemReq) error {
	ret := _m.Called(listID, itemID, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, int64, *models.UpdateItemReq) error); ok {
		r0 = rf(listID, itemID, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
