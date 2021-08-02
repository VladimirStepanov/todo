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
		ls := new(mocks.ListService)
		ls.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(tc.createRetID, tc.createRetErr)
		tsObj := new(mocks.TokenService)
		tsObj.On("Verify", mock.Anything).Return(tc.verifyRetUserID, tc.verifyRetUserUUID, tc.verifyRerErr)
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
		ls := new(mocks.ListService)
		ls.On("GetListByID", mock.Anything, mock.Anything).Return(tc.getRetList, tc.getRetErr)
		tsObj := new(mocks.TokenService)
		tsObj.On("Verify", mock.Anything).Return(tc.verifyRetUserID, tc.verifyRetUserUUID, tc.verifyRerErr)
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
	}
}
