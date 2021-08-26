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
	testList = &models.List{
		ID:          1,
		Title:       "hello",
		Description: "world",
	}
)

func TestListCreate(t *testing.T) {
	tests := []struct {
		name              string
		headers           map[string]string
		input             string
		verifyRetUserID   int64
		verifyRetUserUUID string
		verifyRerErr      error
		createRetID       int64
		createRetErr      error
		code              int
		expListID         int64
		errMsg            string
	}{
		{
			name:              "Create return error",
			headers:           map[string]string{"Authorization": "Bearer token"},
			input:             `{"title": "title", "description": "description"}`,
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      nil,
			createRetID:       0,
			createRetErr:      ErrUnknown,
			code:              http.StatusInternalServerError,
			expListID:         0,
			errMsg:            "Internal server error",
		},
		{
			name:              "Create return error",
			headers:           map[string]string{"Authorization": "Bearer token"},
			input:             `{"title": "title", "description": "description"}`,
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      nil,
			createRetID:       1,
			createRetErr:      nil,
			code:              http.StatusOK,
			expListID:         1,
			errMsg:            "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ls := new(mocks.ListService)
			ls.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(
				tc.createRetID, tc.createRetErr,
			)

			tsObj := new(mocks.TokenService)
			tsObj.On("Verify", mock.Anything).Return(
				tc.verifyRetUserID, tc.verifyRetUserUUID, tc.verifyRerErr,
			)

			handler := New(nil, nil, tsObj, ls, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodPost,
				"/api/lists",
				bytes.NewBuffer([]byte(tc.input)),
				tc.headers,
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
				require.Equal(t, tc.expListID, int64(actResp["list_id"].(float64)))
			}
		})
	}
}

func TestGetListByID(t *testing.T) {
	tests := []struct {
		name              string
		headers           map[string]string
		paramListID       string
		verifyRetUserID   int64
		verifyRetUserUUID string
		verifyRerErr      error
		getRetList        *models.List
		getRetErr         error
		code              int
		expList           *models.List
		errMsg            string
	}{
		{
			name:              "Bad param request",
			headers:           map[string]string{"Authorization": "Bearer token"},
			paramListID:       "bad",
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      nil,
			getRetList:        nil,
			getRetErr:         nil,
			code:              http.StatusBadRequest,
			expList:           nil,
			errMsg:            models.ErrBadParam.Error(),
		},
		{
			name:              "Get return unknown error",
			headers:           map[string]string{"Authorization": "Bearer token"},
			paramListID:       "1",
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      nil,
			getRetList:        nil,
			getRetErr:         ErrUnknown,
			code:              http.StatusInternalServerError,
			expList:           nil,
			errMsg:            "Internal server error",
		},
		{
			name:              "List not found",
			headers:           map[string]string{"Authorization": "Bearer token"},
			paramListID:       "1",
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      nil,
			getRetList:        nil,
			getRetErr:         models.ErrNoList,
			code:              http.StatusNotFound,
			expList:           nil,
			errMsg:            models.ErrNoList.Error(),
		},
		{
			name:              "Success get",
			headers:           map[string]string{"Authorization": "Bearer token"},
			paramListID:       "1",
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      nil,
			getRetList:        testList,
			getRetErr:         nil,
			code:              http.StatusOK,
			expList:           testList,
			errMsg:            models.ErrNoList.Error(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ls := new(mocks.ListService)
			ls.On("GetListByID", mock.Anything, mock.Anything).Return(tc.getRetList, tc.getRetErr)
			tsObj := new(mocks.TokenService)
			tsObj.On("Verify", mock.Anything).Return(
				tc.verifyRetUserID, tc.verifyRetUserUUID,
				tc.verifyRerErr,
			)
			handler := New(nil, nil, tsObj, ls, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodGet,
				fmt.Sprintf("/api/lists/%s", tc.paramListID),
				bytes.NewBuffer([]byte{}),
				tc.headers,
			)
			require.Equal(t, tc.code, code)
			if tc.code != 200 {
				errResp := &ErrorResponse{}
				err := json.Unmarshal(data, errResp)
				require.NoError(t, err)
				require.Equal(t, "error", errResp.Status)
				require.Equal(t, tc.errMsg, errResp.Message)
			} else {
				userList := &models.List{}
				err := json.Unmarshal(data, userList)
				require.NoError(t, err)
				require.Equal(t, tc.expList, userList)
			}
		})
	}
}

func TestEditRole(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}
	input := `{"user_id": 1, "is_admin": true}`

	tests := []struct {
		name         string
		code         int
		paramListID  string
		isListAdmRet error
		editRoleRet  error
		errMsg       string
	}{
		{
			name:         "onlyAdminAccess parse error",
			code:         http.StatusBadRequest,
			paramListID:  "s",
			isListAdmRet: nil,
			errMsg:       models.ErrBadParam.Error(),
		},
		{
			name:         "IsListAdmin  return ErrNoList",
			code:         http.StatusNotFound,
			paramListID:  "1",
			isListAdmRet: models.ErrNoList,
			errMsg:       models.ErrNoList.Error(),
		},
		{
			name:         "IsListAdmin  return ErrNoListAccess",
			code:         http.StatusForbidden,
			paramListID:  "1",
			isListAdmRet: models.ErrNoListAccess,
			errMsg:       models.ErrNoListAccess.Error(),
		},
		{
			name:         "IsListAdmin  return Internal error",
			code:         http.StatusInternalServerError,
			paramListID:  "1",
			isListAdmRet: ErrUnknown,
			errMsg:       "Internal server error",
		},
		{
			name:         "EditRole return ErrUserNotFound",
			code:         http.StatusNotFound,
			paramListID:  "1",
			isListAdmRet: nil,
			editRoleRet:  models.ErrUserNotFound,
			errMsg:       models.ErrUserNotFound.Error(),
		},
		{
			name:         "EditRole return ErrUnknown",
			code:         http.StatusInternalServerError,
			paramListID:  "1",
			isListAdmRet: nil,
			editRoleRet:  ErrUnknown,
			errMsg:       "Internal server error",
		},
		{
			name:         "EditRole success",
			code:         http.StatusOK,
			paramListID:  "1",
			isListAdmRet: nil,
			editRoleRet:  nil,
			errMsg:       "",
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
				tc.isListAdmRet,
			)
			ls.On("EditRole", mock.Anything, mock.Anything, mock.Anything).Return(
				tc.editRoleRet,
			)

			handler := New(nil, nil, tsObj, ls, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodPatch,
				fmt.Sprintf("/api/lists/%s/edit-role", tc.paramListID),
				bytes.NewBuffer([]byte(input)),
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
			}
		})
	}
}

func TestDeleteList(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}

	tests := []struct {
		name         string
		code         int
		paramListID  string
		deleteRetErr error
		errMsg       string
	}{
		{
			name:         "Bad list id",
			code:         http.StatusBadRequest,
			paramListID:  "abc",
			deleteRetErr: nil,
			errMsg:       models.ErrBadParam.Error(),
		},
		{
			name:         "Delete list not found",
			code:         http.StatusNotFound,
			paramListID:  "1",
			deleteRetErr: models.ErrNoList,
			errMsg:       models.ErrNoList.Error(),
		},
		{
			name:         "Delete return unknown error",
			code:         http.StatusInternalServerError,
			paramListID:  "1",
			deleteRetErr: ErrUnknown,
			errMsg:       "Internal server error",
		},
		{
			name:         "Success delete",
			code:         http.StatusOK,
			paramListID:  "1",
			deleteRetErr: nil,
			errMsg:       "",
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
			ls.On("Delete", mock.Anything).Return(
				tc.deleteRetErr,
			)

			handler := New(nil, nil, tsObj, ls, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodDelete,
				fmt.Sprintf("/api/lists/%s", tc.paramListID),
				bytes.NewBuffer([]byte{}),
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
			}
		})
	}
}

func TestUpdateList(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}

	tests := []struct {
		name         string
		code         int
		input        string
		updateRetErr error
		errMsg       string
	}{
		{
			name:         "Empty args",
			code:         http.StatusBadRequest,
			input:        `{}`,
			updateRetErr: models.ErrUpdateEmptyArgs,
			errMsg:       models.ErrUpdateEmptyArgs.Error(),
		},
		{
			name:         "Title too short",
			code:         http.StatusBadRequest,
			input:        `{"title": "12"}`,
			updateRetErr: models.ErrTitleTooShort,
			errMsg:       models.ErrTitleTooShort.Error(),
		},
		{
			name:         "Title too short",
			code:         http.StatusNotFound,
			input:        `{"title": "123456", "description": "hello world"}`,
			updateRetErr: models.ErrNoList,
			errMsg:       models.ErrNoList.Error(),
		},
		{
			name:         "Unknown error",
			code:         http.StatusInternalServerError,
			input:        `{"title": "123456", "description": "hello world"}`,
			updateRetErr: ErrUnknown,
			errMsg:       "Internal server error",
		},
		{
			name:         "Success update",
			code:         http.StatusOK,
			input:        `{"title": "123456", "description": "hello world"}`,
			updateRetErr: nil,
			errMsg:       "",
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
			ls.On("Update", mock.Anything, mock.Anything).Return(
				tc.updateRetErr,
			)

			handler := New(nil, nil, tsObj, ls, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodPatch,
				fmt.Sprintf("/api/lists/%s", "1"),
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
			}
		})
	}
}

func TestGetUserLists(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}

	tests := []struct {
		name     string
		code     int
		errMsg   string
		retErr   error
		expLists []*models.List
	}{
		{
			name:     "Internal error",
			code:     http.StatusInternalServerError,
			errMsg:   "Internal server error",
			retErr:   ErrUnknown,
			expLists: nil,
		},
		{
			name:     "Success get",
			code:     http.StatusOK,
			errMsg:   "",
			retErr:   nil,
			expLists: helpers.ExpLists,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tsObj := new(mocks.TokenService)
			tsObj.On("Verify", mock.Anything).Return(
				int64(1), "aaa-aaa-aaa-aaa", nil,
			)

			ls := new(mocks.ListService)
			ls.On("GetUserLists", mock.Anything).Return(
				tc.expLists,
				tc.retErr,
			)

			handler := New(nil, nil, tsObj, ls, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodGet,
				"/api/lists",
				bytes.NewBuffer([]byte{}),
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
				resp := UserListsResponse{}
				require.NoError(t, json.Unmarshal(data, &resp))
				require.Equal(t, "success", resp.Status)
				require.Equal(t, tc.expLists, resp.Result)
			}
		})
	}
}
