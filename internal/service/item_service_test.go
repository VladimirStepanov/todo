package service

import (
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/VladimirStepanov/todo-app/internal/models/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testItem = &models.Item{
		ID:          1,
		ListID:      666,
		Title:       "title",
		Description: "Description",
		Done:        false,
	}
)

func TestItemCreate(t *testing.T) {
	tests := []struct {
		name   string
		idRet  int64
		errRet error
		expID  int64
		expErr error
	}{
		{
			name:   "Return error",
			idRet:  0,
			errRet: ErrSome,
			expID:  0,
			expErr: ErrSome,
		},
		{
			name:   "Success create",
			idRet:  100,
			errRet: nil,
			expID:  100,
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ir := new(mocks.ItemRepository)
			ir.On("Create", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.idRet, tc.expErr)

			is := NewItemService(ir)

			id, err := is.Create("title", "description", 1)
			require.Equal(t, tc.expID, id)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestGetItemByID(t *testing.T) {
	tests := []struct {
		name    string
		retErr  error
		retItem *models.Item
		expErr  error
		expItem *models.Item
	}{
		{
			name:    "Return unknown error",
			retErr:  ErrSome,
			retItem: nil,
			expErr:  ErrSome,
			expItem: nil,
		},
		{
			name:    "Item not found",
			retErr:  models.ErrNoItem,
			retItem: nil,
			expErr:  models.ErrNoItem,
			expItem: nil,
		},
		{
			name:    "Success get",
			retErr:  nil,
			retItem: testItem,
			expErr:  nil,
			expItem: testItem,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ir := new(mocks.ItemRepository)
			ir.On("GetItemByID", mock.Anything, mock.Anything).Return(tc.retItem, tc.retErr)

			is := NewItemService(ir)

			retItem, err := is.GetItemByID(testItem.ListID, testItem.ID)
			require.Equal(t, tc.expItem, retItem)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestDeleteItem(t *testing.T) {
	tests := []struct {
		name   string
		retErr error
		expErr error
	}{
		{
			name:   "Return unknown error",
			retErr: ErrSome,
			expErr: ErrSome,
		},
		{
			name:   "Item not found",
			retErr: models.ErrNoItem,
			expErr: models.ErrNoItem,
		},
		{
			name:   "Success delete",
			retErr: nil,
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ir := new(mocks.ItemRepository)
			ir.On("Delete", mock.Anything, mock.Anything).Return(tc.retErr)

			is := NewItemService(ir)

			err := is.Delete(testItem.ListID, testItem.ID)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestItemUpdate(t *testing.T) {
	goodReq := &models.UpdateItemReq{
		Title:       new(string),
		Description: new(string),
	}

	*goodReq.Title = "helllo"
	*goodReq.Description = "world"

	tests := []struct {
		name   string
		req    func() *models.UpdateItemReq
		retErr error
		expErr error
	}{
		{
			name: "Empty arguments",
			req: func() *models.UpdateItemReq {
				return &models.UpdateItemReq{
					Title:       nil,
					Description: nil,
				}
			},
			retErr: nil,
			expErr: models.ErrUpdateEmptyArgs,
		},
		{
			name: "Title too short error",
			req: func() *models.UpdateItemReq {
				title := "123"
				return &models.UpdateItemReq{
					Title:       &title,
					Description: nil,
				}
			},
			retErr: nil,
			expErr: models.ErrTitleTooShort,
		},
		{
			name: "Return unknown error",
			req: func() *models.UpdateItemReq {
				return goodReq
			},
			retErr: ErrSome,
			expErr: ErrSome,
		},
		{
			name: "Return ErrNoItem error",
			req: func() *models.UpdateItemReq {
				return goodReq
			},
			retErr: models.ErrNoItem,
			expErr: models.ErrNoItem,
		},
		{
			name: "Success update",
			req: func() *models.UpdateItemReq {
				return goodReq
			},
			retErr: nil,
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ir := new(mocks.ItemRepository)
			ir.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(tc.retErr)

			is := NewItemService(ir)

			err := is.Update(1, 1, tc.req())
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestGetItems(t *testing.T) {
	result := []*models.Item{
		{
			ID:          1,
			Title:       "title#1",
			Description: "description#1",
		},
		{
			ID:          2,
			Title:       "title#2",
			Description: "description#2",
		},
	}

	tests := []struct {
		name   string
		retErr error
		expErr error
		expRes []*models.Item
	}{
		{
			name:   "Return unknown error",
			retErr: ErrSome,
			expErr: ErrSome,
			expRes: nil,
		},
		{
			name:   "Success get",
			retErr: nil,
			expErr: nil,
			expRes: result,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ir := new(mocks.ItemRepository)
			ir.On("GetItems", mock.Anything).Return(tc.expRes, tc.retErr)

			is := NewItemService(ir)

			res, err := is.GetItems(1)
			require.Equal(t, tc.expErr, err)
			require.Equal(t, tc.expRes, res)
		})
	}
}
