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

func TestSignIn(t *testing.T) {
	retUser := &models.User{ID: 1}

	tests := []struct {
		name      string
		usRetUser *models.User
		usRetErr  error
		tsRetErr  error
		tsRetTd   *models.TokenDetails
		code      int
		errMsg    string
	}{
		{
			name:      "Internal error for SignIn",
			usRetUser: nil,
			usRetErr:  ErrUnknown,
			tsRetErr:  nil,
			tsRetTd:   nil,
			code:      http.StatusInternalServerError,
			errMsg:    "Internal server error",
		},
		{
			name:      "Not activated user error",
			usRetUser: nil,
			usRetErr:  models.ErrUserNotActivated,
			tsRetErr:  nil,
			tsRetTd:   nil,
			code:      http.StatusForbidden,
			errMsg:    models.ErrUserNotActivated.Error(),
		},
		{
			name:      "User not found",
			usRetUser: nil,
			usRetErr:  models.ErrBadUser,
			tsRetErr:  nil,
			tsRetTd:   nil,
			code:      http.StatusNotFound,
			errMsg:    models.ErrBadUser.Error(),
		},
		{
			name:      "Internal error for NewTokenPair",
			usRetUser: retUser,
			usRetErr:  nil,
			tsRetErr:  ErrUnknown,
			tsRetTd:   nil,
			code:      http.StatusInternalServerError,
			errMsg:    "Internal server error",
		},
		{
			name:      "Max user logged in error",
			usRetUser: retUser,
			usRetErr:  nil,
			tsRetErr:  models.ErrMaxLoggedIn,
			tsRetTd:   nil,
			code:      http.StatusForbidden,
			errMsg:    models.ErrMaxLoggedIn.Error(),
		},
		{
			name:      "Success",
			usRetUser: retUser,
			usRetErr:  nil,
			tsRetErr:  nil,
			tsRetTd: &models.TokenDetails{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
				UUID:         "",
			},
			code:   http.StatusOK,
			errMsg: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reqData := `{"email": "test@test.com", "password": "123456789"}`
			usObj := new(mocks.UserService)
			usObj.On("SignIn", mock.Anything, mock.Anything).Return(tc.usRetUser, tc.usRetErr)
			tsObj := new(mocks.TokenService)
			tsObj.On("NewTokenPair", mock.Anything).Return(tc.tsRetTd, tc.tsRetErr)

			handler := New(usObj, nil, tsObj, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodPost,
				"/auth/sign-in",
				bytes.NewBuffer([]byte(reqData)),
			)
			require.Equal(t, tc.code, code)
			actResp := map[string]interface{}{}
			err := json.Unmarshal(data, &actResp)
			require.NoError(t, err)
			if tc.code != 200 {
				require.Equal(t, "error", actResp["status"])
				require.Equal(t, tc.errMsg, actResp["message"])
			} else {
				require.NotEmpty(t, actResp["access_token"])
				require.NotEmpty(t, actResp["refresh_token"])
			}
		})
	}
}
