package it

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/VladimirStepanov/todo-app/internal/handler"
	"github.com/VladimirStepanov/todo-app/internal/helpers"
	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/VladimirStepanov/todo-app/internal/service"
	"github.com/stretchr/testify/require"
)

func makeLogout(t *testing.T, r http.Handler, data []byte) {
	tokenResp := &handler.TokensResponse{}
	err := json.Unmarshal(data, tokenResp)
	require.NoError(t, err)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", tokenResp.AccessToken),
	}
	code, _ := helpers.MakeRequest(
		r,
		t,
		http.MethodGet,
		"/auth/logout",
		bytes.NewBuffer([]byte{}),
		headers,
	)
	require.Equal(t, http.StatusOK, code)
}

func (suite *TestingSuite) TestSignUp() {

	tests := []struct {
		name   string
		data   string
		code   int
		errMsg string
	}{
		{
			"Success create",
			`{"email": "new@test.com", "password": "123456789"}`,
			http.StatusOK,
			"",
		},
		{
			"User already exists",
			fmt.Sprintf(`{"email": "%s", "password": "123456789"}`, userForCreate.Email),
			http.StatusConflict,
			models.ErrUserAlreadyExists.Error(),
		},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			code, data := helpers.MakeRequest(
				suite.router,
				suite.T(),
				http.MethodPost,
				"/auth/sign-up",
				bytes.NewBuffer([]byte(tc.data)),
				nil,
			)
			require.Equal(t, tc.code, code)
			if tc.errMsg != "" {
				resp := handler.ErrorResponse{}
				err := json.Unmarshal(data, &resp)
				require.NoError(t, err)
				require.Equal(t, "error", resp.Status)
				require.Equal(t, tc.errMsg, resp.Message)
			}
		})
	}
}

func (suite *TestingSuite) TestEmailConfirmation() {
	tests := []struct {
		name   string
		link   string
		code   int
		errMsg string
	}{
		{
			"Success confirmation",
			notConfirmedUser.ActivatedLink,
			http.StatusOK,
			"",
		},
		{
			"Already confirmation",
			confirmedUser.ActivatedLink,
			http.StatusNotFound,
			"Page not found",
		},
		{
			"Unknown confirmation link",
			unknownConfLink,
			http.StatusNotFound,
			"Page not found",
		},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			code, data := helpers.MakeRequest(
				suite.router,
				suite.T(),
				http.MethodGet,
				fmt.Sprintf("/auth/confirm/%s", tc.link),
				bytes.NewBuffer([]byte{}),
				nil,
			)
			require.Equal(t, tc.code, code)
			if tc.errMsg != "" && tc.code != http.StatusOK {
				resp := handler.ErrorResponse{}
				err := json.Unmarshal(data, &resp)
				require.NoError(t, err)
				require.Equal(t, "error", resp.Status)
				require.Equal(t, tc.errMsg, resp.Message)

			}
		})
	}
}

func (suite *TestingSuite) TestSignIn() {

	input := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, maxLoggedInUser.Email, defaultPassword)

	for i := 0; i < maxLoggenInCount; i++ {
		code, _ := helpers.MakeRequest(
			suite.router,
			suite.T(),
			http.MethodPost,
			"/auth/sign-in",
			bytes.NewBuffer([]byte(input)),
			nil,
		)
		require.Equal(suite.T(), http.StatusOK, code)
	}

	tests := []struct {
		name   string
		input  string
		code   int
		errMsg string
	}{
		{
			name:   "Max logged in users",
			input:  fmt.Sprintf(`{"email": "%s", "password": "%s"}`, maxLoggedInUser.Email, defaultPassword),
			code:   http.StatusForbidden,
			errMsg: models.ErrMaxLoggedIn.Error(),
		},
		{
			name:   "Check ErrBadUser error",
			input:  fmt.Sprintf(`{"email": "%s", "password": "%s"}`, "bad@user.com", defaultPassword),
			code:   http.StatusNotFound,
			errMsg: models.ErrBadUser.Error(),
		},
		{
			name:   "Check ErrUserNotActivated error",
			input:  fmt.Sprintf(`{"email": "%s", "password": "%s"}`, authNotConfirmedUser.Email, defaultPassword),
			code:   http.StatusForbidden,
			errMsg: models.ErrUserNotActivated.Error(),
		},
		{
			name:   "Success auth",
			input:  fmt.Sprintf(`{"email": "%s", "password": "%s"}`, authUser.Email, defaultPassword),
			code:   http.StatusOK,
			errMsg: "",
		},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			code, data := helpers.MakeRequest(
				suite.router,
				suite.T(),
				http.MethodPost,
				"/auth/sign-in",
				bytes.NewBuffer([]byte(tc.input)),
				nil,
			)
			require.Equal(t, tc.code, code)
			if tc.errMsg != "" && tc.code != http.StatusOK {
				resp := handler.ErrorResponse{}
				err := json.Unmarshal(data, &resp)
				require.NoError(t, err)
				require.Equal(t, "error", resp.Status)
				require.Equal(t, tc.errMsg, resp.Message)
			} else if tc.code == http.StatusOK {
				makeLogout(t, suite.router, data)
			}
		})
	}
}

func (suite *TestingSuite) TestRefresh() {
	inputF := `{"refresh_token": "%s"}`

	siginInput := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, authUser.Email, defaultPassword)
	code, data := helpers.MakeRequest(
		suite.router,
		suite.T(),
		http.MethodPost,
		"/auth/sign-in",
		bytes.NewBuffer([]byte(siginInput)),
		nil,
	)
	require.Equal(suite.T(), http.StatusOK, code)

	authResp := &handler.TokensResponse{}
	err := json.Unmarshal(data, authResp)
	require.NoError(suite.T(), err)

	tests := []struct {
		name   string
		input  func() string
		code   int
		errMsg string
	}{
		{
			name: "Bad token error",
			input: func() string {
				return fmt.Sprintf(inputF, "bad.bad.bad")
			},
			code:   http.StatusForbidden,
			errMsg: models.ErrBadToken.Error(),
		},
		{
			name: "Expired token error",
			input: func() string {
				token, err := service.GenerateToken(testUUID, authUser.ID, 100, 103, refreshKey)
				require.NoError(suite.T(), err)
				return fmt.Sprintf(inputF, token)
			},
			code:   http.StatusUnauthorized,
			errMsg: models.ErrTokenExpired.Error(),
		},
		{
			name: "User unauthorized",
			input: func() string {
				token, err := service.GenerateToken(
					testUUID, notConfirmedUser.ID,
					time.Now().Unix(),
					time.Now().Add(time.Hour).Unix(), refreshKey,
				)
				require.NoError(suite.T(), err)
				return fmt.Sprintf(inputF, token)
			},
			code:   http.StatusUnauthorized,
			errMsg: models.ErrUserUnauthorized.Error(),
		},
		{
			name: "Success refresh",
			input: func() string {
				return fmt.Sprintf(inputF, authResp.RefreshToken)
			},
			code:   http.StatusOK,
			errMsg: "",
		},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			code, data := helpers.MakeRequest(
				suite.router,
				suite.T(),
				http.MethodPost,
				"/auth/refresh",
				bytes.NewBuffer([]byte(tc.input())),
				nil,
			)
			require.Equal(t, tc.code, code)
			if tc.errMsg != "" && tc.code != http.StatusOK {
				resp := handler.ErrorResponse{}
				err := json.Unmarshal(data, &resp)
				require.NoError(t, err)
				require.Equal(t, "error", resp.Status)
				require.Equal(t, tc.errMsg, resp.Message)
			} else if tc.code == http.StatusOK {
				makeLogout(t, suite.router, data)
			}
		})
	}
}
