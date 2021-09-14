package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/helpers"
	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/VladimirStepanov/todo-app/internal/models/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testItem = &models.Item{
		ID:          1,
		ListID:      testList.ID,
		Title:       "hello",
		Description: "world",
		Done:        true,
	}
)

func TestItemCreate(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}

	tests := []struct {
		name        string
		code        int
		input       string
		listID      string
		listServErr error
		crExpRetID  int64
		crExpRetErr error
		errMsg      string
	}{
		{
			name:        "Bad list input",
			code:        http.StatusBadRequest,
			input:       "",
			listID:      "hello",
			listServErr: nil,
			crExpRetID:  0,
			crExpRetErr: nil,
			errMsg:      models.ErrBadParam.Error(),
		},
		{
			name:        "IsListAdmin return unknown error",
			code:        http.StatusInternalServerError,
			input:       "",
			listID:      "1",
			listServErr: ErrUnknown,
			crExpRetID:  0,
			crExpRetErr: nil,
			errMsg:      "Internal server error",
		},
		{
			name:        "Create return unknown error",
			code:        http.StatusInternalServerError,
			input:       `{"title": "title", "description": "description"}`,
			listID:      "1",
			listServErr: nil,
			crExpRetID:  0,
			crExpRetErr: ErrUnknown,
			errMsg:      "Internal server error",
		},
		{
			name:        "Success create",
			code:        http.StatusOK,
			input:       `{"title": "title", "description": "description"}`,
			listID:      "1",
			listServErr: nil,
			crExpRetID:  777,
			crExpRetErr: nil,
			errMsg:      "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tsObj := new(mocks.TokenService)
			tsObj.On("Verify", mock.Anything).Return(
				int64(1), "aaa-aaa-aaa-aaa", nil,
			)

			ls := new(mocks.ListService)
			ls.On("IsListAdmin", mock.Anything, mock.Anything).Return(
				tc.listServErr,
			)

			is := new(mocks.ItemService)
			is.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(
				tc.crExpRetID, tc.crExpRetErr,
			)

			handler := New(nil, nil, tsObj, ls, is, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodPost,
				fmt.Sprintf("/api/lists/%s/items", tc.listID),
				bytes.NewBuffer([]byte(tc.input)),
				headers,
			)
			require.Equal(t, tc.code, code)
			actResp := map[string]interface{}{}
			err := json.Unmarshal(data, &actResp)
			require.NoError(t, err)
			if tc.code != 200 {
				require.Equal(t, "error", actResp["status"])
				require.Equal(t, tc.errMsg, actResp["message"])
			} else {
				require.Equal(t, "success", actResp["status"])
				require.Equal(t, tc.crExpRetID, int64(actResp["item_id"].(float64)))
			}
		})
	}

}

func TestGetItemByID(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}
	tests := []struct {
		name       string
		listID     string
		itemID     string
		getRetItem *models.Item
		getRetErr  error
		code       int
		expItem    *models.Item
		errMsg     string
	}{
		{
			name:       "Bad itemID",
			listID:     "1",
			itemID:     "bad",
			getRetItem: nil,
			getRetErr:  nil,
			code:       http.StatusBadRequest,
			expItem:    nil,
			errMsg:     models.ErrBadParam.Error(),
		},
		{
			name:       "Item not found",
			listID:     "1",
			itemID:     "2",
			getRetItem: nil,
			getRetErr:  models.ErrNoItem,
			code:       http.StatusNotFound,
			expItem:    nil,
			errMsg:     models.ErrNoItem.Error(),
		},
		{
			name:       "GetItemByID return unknown error",
			listID:     "1",
			itemID:     "2",
			getRetItem: nil,
			getRetErr:  ErrUnknown,
			code:       http.StatusInternalServerError,
			expItem:    nil,
			errMsg:     "Internal server error",
		},
		{
			name:       "Success",
			listID:     fmt.Sprintf("%d", testList.ID),
			itemID:     fmt.Sprintf("%d", testItem.ID),
			getRetItem: testItem,
			getRetErr:  nil,
			code:       http.StatusOK,
			expItem:    testItem,
			errMsg:     "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tsObj := new(mocks.TokenService)
			tsObj.On("Verify", mock.Anything).Return(
				int64(1), "aaa-aaa-aaa-aaa", nil,
			)

			ls := new(mocks.ListService)
			ls.On("IsListAdmin", mock.Anything, mock.Anything).Return(
				nil,
			)

			is := new(mocks.ItemService)
			is.On("GetItemByID", mock.Anything, mock.Anything).Return(
				tc.getRetItem, tc.getRetErr,
			)

			handler := New(nil, nil, tsObj, ls, is, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodGet,
				fmt.Sprintf("/api/lists/%s/items/%s", tc.listID, tc.itemID),
				bytes.NewBuffer([]byte{}),
				headers,
			)
			require.Equal(t, tc.code, code)
			if tc.code != 200 {
				errResp := &ErrorResponse{}
				err := json.Unmarshal(data, errResp)
				require.NoError(t, err)
				require.Equal(t, "error", errResp.Status)
				require.Equal(t, tc.errMsg, errResp.Message)
			} else {
				item := &models.Item{}
				err := json.Unmarshal(data, item)
				require.NoError(t, err)
				require.Equal(t, tc.expItem, item)
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}
	tests := []struct {
		name   string
		listID string
		itemID string
		retErr error
		code   int
		errMsg string
	}{
		{
			name:   "Bad itemID",
			listID: "1",
			itemID: "bad",
			retErr: nil,
			code:   http.StatusBadRequest,
			errMsg: models.ErrBadParam.Error(),
		},
		{
			name:   "Item not found",
			listID: "1",
			itemID: "2",
			retErr: models.ErrNoItem,
			code:   http.StatusNotFound,
			errMsg: models.ErrNoItem.Error(),
		},
		{
			name:   "Delete return unknown error",
			listID: "1",
			itemID: "2",
			retErr: ErrUnknown,
			code:   http.StatusInternalServerError,
			errMsg: "Internal server error",
		},
		{
			name:   "Success",
			listID: fmt.Sprintf("%d", testList.ID),
			itemID: fmt.Sprintf("%d", testItem.ID),
			retErr: nil,
			code:   http.StatusOK,
			errMsg: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tsObj := new(mocks.TokenService)
			tsObj.On("Verify", mock.Anything).Return(
				int64(1), "aaa-aaa-aaa-aaa", nil,
			)

			ls := new(mocks.ListService)
			ls.On("IsListAdmin", mock.Anything, mock.Anything).Return(
				nil,
			)

			is := new(mocks.ItemService)
			is.On("Delete", mock.Anything, mock.Anything).Return(
				tc.retErr,
			)

			handler := New(nil, nil, tsObj, ls, is, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodDelete,
				fmt.Sprintf("/api/lists/%s/items/%s", tc.listID, tc.itemID),
				bytes.NewBuffer([]byte{}),
				headers,
			)
			require.Equal(t, tc.code, code)
			if tc.code != 200 {
				errResp := &ErrorResponse{}
				err := json.Unmarshal(data, errResp)
				require.NoError(t, err)
				require.Equal(t, "error", errResp.Status)
				require.Equal(t, tc.errMsg, errResp.Message)
			}
		})
	}
}

func TestUpdateItem(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}
	tests := []struct {
		name   string
		listID string
		itemID string
		input  string
		retErr error
		code   int
		errMsg string
	}{
		{
			name:   "Bad itemID",
			listID: "1",
			itemID: "bad",
			input:  `{}`,
			retErr: nil,
			code:   http.StatusBadRequest,
			errMsg: models.ErrBadParam.Error(),
		},
		{
			name:   "Empty args",
			listID: "1",
			itemID: "1",
			code:   http.StatusBadRequest,
			input:  `{}`,
			retErr: models.ErrUpdateEmptyArgs,
			errMsg: models.ErrUpdateEmptyArgs.Error(),
		},
		{
			name:   "Title too short",
			listID: "1",
			itemID: "1",
			code:   http.StatusBadRequest,
			input:  `{"title": "12"}`,
			retErr: models.ErrTitleTooShort,
			errMsg: models.ErrTitleTooShort.Error(),
		},
		{
			name:   "Return ErrNoItem",
			listID: "1",
			itemID: "1",
			code:   http.StatusNotFound,
			input:  `{"title": "123456", "description": "hello world"}`,
			retErr: models.ErrNoItem,
			errMsg: models.ErrNoItem.Error(),
		},
		{
			name:   "Unknown error",
			listID: "1",
			itemID: "1",
			code:   http.StatusInternalServerError,
			input:  `{"title": "123456", "description": "hello world"}`,
			retErr: ErrUnknown,
			errMsg: "Internal server error",
		},
		{
			name:   "Success update",
			listID: "1",
			itemID: "1",
			code:   http.StatusOK,
			input:  `{"title": "123456", "description": "hello world"}`,
			retErr: nil,
			errMsg: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tsObj := new(mocks.TokenService)
			tsObj.On("Verify", mock.Anything).Return(
				int64(1), "aaa-aaa-aaa-aaa", nil,
			)

			ls := new(mocks.ListService)
			ls.On("IsListAdmin", mock.Anything, mock.Anything).Return(
				nil,
			)

			is := new(mocks.ItemService)
			is.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(
				tc.retErr,
			)

			handler := New(nil, nil, tsObj, ls, is, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodPatch,
				fmt.Sprintf("/api/lists/%s/items/%s", tc.listID, tc.itemID),
				bytes.NewBuffer([]byte(tc.input)),
				headers,
			)
			require.Equal(t, tc.code, code)
			if tc.code != 200 {
				errResp := &ErrorResponse{}
				err := json.Unmarshal(data, errResp)
				require.NoError(t, err)
				require.Equal(t, "error", errResp.Status)
				require.Equal(t, tc.errMsg, errResp.Message)
			}
		})
	}
}

func TestDoneItem(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}
	tests := []struct {
		name   string
		listID string
		itemID string
		retErr error
		code   int
		errMsg string
	}{
		{
			name:   "Bad itemID",
			listID: "1",
			itemID: "bad",
			retErr: nil,
			code:   http.StatusBadRequest,
			errMsg: models.ErrBadParam.Error(),
		},
		{
			name:   "Empty args",
			listID: "1",
			itemID: "1",
			code:   http.StatusBadRequest,
			retErr: models.ErrUpdateEmptyArgs,
			errMsg: models.ErrUpdateEmptyArgs.Error(),
		},
		{
			name:   "Title too short",
			listID: "1",
			itemID: "1",
			code:   http.StatusBadRequest,
			retErr: models.ErrTitleTooShort,
			errMsg: models.ErrTitleTooShort.Error(),
		},
		{
			name:   "Return ErrNoItem",
			listID: "1",
			itemID: "1",
			code:   http.StatusNotFound,
			retErr: models.ErrNoItem,
			errMsg: models.ErrNoItem.Error(),
		},
		{
			name:   "Unknown error",
			listID: "1",
			itemID: "1",
			code:   http.StatusInternalServerError,
			retErr: ErrUnknown,
			errMsg: "Internal server error",
		},
		{
			name:   "Success done",
			listID: "1",
			itemID: "1",
			code:   http.StatusOK,
			retErr: nil,
			errMsg: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tsObj := new(mocks.TokenService)
			tsObj.On("Verify", mock.Anything).Return(
				int64(1), "aaa-aaa-aaa-aaa", nil,
			)

			ls := new(mocks.ListService)
			ls.On("IsListAdmin", mock.Anything, mock.Anything).Return(
				nil,
			)

			is := new(mocks.ItemService)
			is.On("Done", mock.Anything, mock.Anything).Return(
				tc.retErr,
			)

			handler := New(nil, nil, tsObj, ls, is, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodPatch,
				fmt.Sprintf("/api/lists/%s/items/%s/done", tc.listID, tc.itemID),
				bytes.NewBuffer([]byte{}),
				headers,
			)
			require.Equal(t, tc.code, code)
			if tc.code != 200 {
				errResp := &ErrorResponse{}
				err := json.Unmarshal(data, errResp)
				require.NoError(t, err)
				require.Equal(t, "error", errResp.Status)
				require.Equal(t, tc.errMsg, errResp.Message)
			}
		})
	}
}

func TestGetItems(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}
	tests := []struct {
		name     string
		listID   string
		retErr   error
		retItems []*models.Item
		code     int
		expItems []*models.Item
		errMsg   string
	}{
		{
			name:   "Unknown error",
			listID: "1",
			code:   http.StatusInternalServerError,
			retErr: ErrUnknown,
			errMsg: "Internal server error",
		},
		{
			name:     "Success get items",
			listID:   "1",
			code:     http.StatusOK,
			retItems: helpers.ExpItems,
			expItems: helpers.ExpItems,
			errMsg:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tsObj := new(mocks.TokenService)
			tsObj.On("Verify", mock.Anything).Return(
				int64(1), "aaa-aaa-aaa-aaa", nil,
			)

			ls := new(mocks.ListService)
			ls.On("IsListAdmin", mock.Anything, mock.Anything).Return(
				nil,
			)

			is := new(mocks.ItemService)
			is.On("GetItems", mock.Anything).Return(
				tc.retItems, tc.retErr,
			)

			handler := New(nil, nil, tsObj, ls, is, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodGet,
				fmt.Sprintf("/api/lists/%s/items", tc.listID),
				bytes.NewBuffer([]byte{}),
				headers,
			)
			require.Equal(t, tc.code, code)
			if tc.code != 200 {
				errResp := &ErrorResponse{}
				err := json.Unmarshal(data, errResp)
				require.NoError(t, err)
				require.Equal(t, "error", errResp.Status)
				require.Equal(t, tc.errMsg, errResp.Message)
			} else {
				resp := UserItemsResponse{}
				require.NoError(t, json.Unmarshal(data, &resp))
				require.Equal(t, "success", resp.Status)
				require.Equal(t, tc.expItems, resp.Result)
			}
		})
	}
}
