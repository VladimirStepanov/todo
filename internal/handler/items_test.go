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

func TestAccessListMiddleware(t *testing.T) {
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
