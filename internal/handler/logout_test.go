package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/helpers"
	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/VladimirStepanov/todo-app/internal/models/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogout(t *testing.T) {
	tests := []struct {
		name              string
		headers           map[string]string
		verifyRetUserID   int64
		verifyRetUserUUID string
		verifyRerErr      error
		logoutRetErr      error
		code              int
		errMsg            string
	}{
		{
			name:              "No auth header",
			headers:           map[string]string{},
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      nil,
			logoutRetErr:      nil,
			code:              http.StatusUnauthorized,
			errMsg:            models.ErrNoAuthHeader.Error(),
		},
		{
			name: "No Bearer",
			headers: map[string]string{
				"Authorization": "one two",
			},
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      nil,
			logoutRetErr:      nil,
			code:              http.StatusUnauthorized,
			errMsg:            models.ErrInvalidAuthHeader.Error(),
		},
		{
			name: "More than two values in header value",
			headers: map[string]string{
				"Authorization": "one two three",
			},
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      nil,
			logoutRetErr:      nil,
			code:              http.StatusUnauthorized,
			errMsg:            models.ErrInvalidAuthHeader.Error(),
		},
		{
			name: "Verify return models.ErrBadToken",
			headers: map[string]string{
				"Authorization": "Bearer token",
			},
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      models.ErrBadToken,
			logoutRetErr:      nil,
			code:              http.StatusForbidden,
			errMsg:            models.ErrBadToken.Error(),
		},
		{
			name: "Verify return models.ErrUserUnauthorized",
			headers: map[string]string{
				"Authorization": "Bearer token",
			},
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      models.ErrUserUnauthorized,
			logoutRetErr:      nil,
			code:              http.StatusUnauthorized,
			errMsg:            models.ErrUserUnauthorized.Error(),
		},
		{
			name: "Verify return models.ErrTokenExpired",
			headers: map[string]string{
				"Authorization": "Bearer token",
			},
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      models.ErrTokenExpired,
			logoutRetErr:      nil,
			code:              http.StatusUnauthorized,
			errMsg:            models.ErrTokenExpired.Error(),
		},
		{
			name: "Verify return unknown error",
			headers: map[string]string{
				"Authorization": "Bearer token",
			},
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      ErrUnknown,
			logoutRetErr:      nil,
			code:              http.StatusInternalServerError,
			errMsg:            "Internal server error",
		},
		{
			name: "Logout return unknown error",
			headers: map[string]string{
				"Authorization": "Bearer token",
			},
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      nil,
			logoutRetErr:      ErrUnknown,
			code:              http.StatusInternalServerError,
			errMsg:            "Internal server error",
		},
		{
			name: "Success logout",
			headers: map[string]string{
				"Authorization": "Bearer token",
			},
			verifyRetUserID:   0,
			verifyRetUserUUID: "",
			verifyRerErr:      nil,
			logoutRetErr:      nil,
			code:              http.StatusOK,
			errMsg:            "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tsObj := new(mocks.TokenService)
			tsObj.On("Verify", mock.Anything).Return(tc.verifyRetUserID, tc.verifyRetUserUUID, tc.verifyRerErr)
			tsObj.On("Logout", mock.Anything, mock.Anything).Return(tc.logoutRetErr)

			handler := New(nil, nil, tsObj, nil, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodGet,
				"/auth/logout",
				bytes.NewBuffer([]byte{}),
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
			}
		})
	}
}
