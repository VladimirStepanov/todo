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

func TestRefresh(t *testing.T) {
	tests := []struct {
		name     string
		tsRetErr error
		tsRetTd  *models.TokenDetails
		code     int
		errMsg   string
	}{
		{
			name:     "Return bad token error",
			tsRetErr: models.ErrBadToken,
			tsRetTd:  nil,
			code:     http.StatusBadRequest,
			errMsg:   models.ErrBadToken.Error(),
		},
		{
			name:     "Return bad token expired error",
			tsRetErr: models.ErrTokenExpired,
			tsRetTd:  nil,
			code:     http.StatusUnauthorized,
			errMsg:   models.ErrTokenExpired.Error(),
		},
		{
			name:     "Return user unauthorized error",
			tsRetErr: models.ErrUserUnauthorized,
			tsRetTd:  nil,
			code:     http.StatusUnauthorized,
			errMsg:   models.ErrUserUnauthorized.Error(),
		},
		{
			name:     "Return unknown error",
			tsRetErr: ErrUnknown,
			tsRetTd:  nil,
			code:     http.StatusInternalServerError,
			errMsg:   "Internal server error",
		},
		{
			name:     "Success",
			tsRetErr: nil,
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
			reqData := `{"refresh_token": "token"}`
			tsObj := new(mocks.TokenService)
			tsObj.On("Refresh", mock.Anything).Return(tc.tsRetTd, tc.tsRetErr)

			handler := New(nil, nil, tsObj, nil, nil, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodPost,
				"/auth/refresh",
				bytes.NewBuffer([]byte(reqData)),
				nil,
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
