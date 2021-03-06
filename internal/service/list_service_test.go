package service

import (
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/VladimirStepanov/todo-app/internal/models/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testList = &models.List{
		ID:          1,
		Title:       "hello",
		Description: "world",
	}
)

func TestListCreate(t *testing.T) {
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
			lr := new(mocks.ListRepository)
			lr.On("Create", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.idRet, tc.expErr)

			ls := NewListService(lr)

			id, err := ls.Create("title", "description", 1)
			require.Equal(t, tc.expID, id)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestListGetByID(t *testing.T) {
	tests := []struct {
		name    string
		retErr  error
		retList *models.List
		expErr  error
		expList *models.List
	}{
		{
			name:    "Return unknown error",
			retErr:  ErrSome,
			retList: nil,
			expErr:  ErrSome,
			expList: nil,
		},
		{
			name:    "List not found",
			retErr:  models.ErrNoList,
			retList: nil,
			expErr:  models.ErrNoList,
			expList: nil,
		},
		{
			name:    "Success get",
			retErr:  nil,
			retList: testList,
			expErr:  nil,
			expList: testList,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lr := new(mocks.ListRepository)
			lr.On("GetListByID", mock.Anything, mock.Anything).Return(tc.retList, tc.retErr)

			ls := NewListService(lr)

			retList, err := ls.GetListByID(1, 1)
			require.Equal(t, tc.expList, retList)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestIsListAdmin(t *testing.T) {
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
			name:   "List not found",
			retErr: models.ErrNoList,
			expErr: models.ErrNoList,
		},
		{
			name:   "No access",
			retErr: models.ErrNoListAccess,
			expErr: models.ErrNoListAccess,
		},
		{
			name:   "Success admin access",
			retErr: nil,
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lr := new(mocks.ListRepository)
			lr.On("IsListAdmin", mock.Anything, mock.Anything).Return(tc.retErr)

			ls := NewListService(lr)

			err := ls.IsListAdmin(1, 1)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestEditRole(t *testing.T) {
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
			name:   "Success grant role",
			retErr: nil,
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lr := new(mocks.ListRepository)
			lr.On("EditRole", mock.Anything, mock.Anything, mock.Anything).Return(tc.retErr)

			ls := NewListService(lr)

			err := ls.EditRole(1, 1, true)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestDeleteList(t *testing.T) {
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
			name:   "Return ErrNoList error",
			retErr: models.ErrNoList,
			expErr: models.ErrNoList,
		},
		{
			name:   "Success grant role",
			retErr: nil,
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lr := new(mocks.ListRepository)
			lr.On("Delete", mock.Anything).Return(tc.retErr)

			ls := NewListService(lr)

			err := ls.Delete(1)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestUpdate(t *testing.T) {
	goodReq := &models.UpdateListReq{
		Title:       new(string),
		Description: new(string),
	}

	*goodReq.Title = "helllo"
	*goodReq.Description = "world"

	tests := []struct {
		name   string
		req    func() *models.UpdateListReq
		retErr error
		expErr error
	}{
		{
			name: "Empty arguments",
			req: func() *models.UpdateListReq {
				return &models.UpdateListReq{
					Title:       nil,
					Description: nil,
				}
			},
			retErr: nil,
			expErr: models.ErrUpdateEmptyArgs,
		},
		{
			name: "Title too short error",
			req: func() *models.UpdateListReq {
				title := "123"
				return &models.UpdateListReq{
					Title:       &title,
					Description: nil,
				}
			},
			retErr: nil,
			expErr: models.ErrTitleTooShort,
		},
		{
			name: "Return unknown error",
			req: func() *models.UpdateListReq {
				return goodReq
			},
			retErr: ErrSome,
			expErr: ErrSome,
		},
		{
			name: "Return ErrNoList error",
			req: func() *models.UpdateListReq {
				return goodReq
			},
			retErr: models.ErrNoList,
			expErr: models.ErrNoList,
		},
		{
			name: "Success update",
			req: func() *models.UpdateListReq {
				return goodReq
			},
			retErr: nil,
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lr := new(mocks.ListRepository)
			lr.On("Update", mock.Anything, mock.Anything).Return(tc.retErr)

			ls := NewListService(lr)

			err := ls.Update(1, tc.req())
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestGetUserLists(t *testing.T) {
	result := []*models.List{
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
		expRes []*models.List
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
			lr := new(mocks.ListRepository)
			lr.On("GetUserLists", mock.Anything).Return(tc.expRes, tc.retErr)

			ls := NewListService(lr)

			res, err := ls.GetUserLists(1)
			require.Equal(t, tc.expErr, err)
			require.Equal(t, tc.expRes, res)
		})
	}
}
