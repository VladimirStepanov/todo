package it

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/handler"
	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/stretchr/testify/require"
)

func (suite *TestingSuite) makeRequest(method, path string, input *bytes.Buffer) (int, []byte) {
	req := httptest.NewRequest(method, path, input)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	res := w.Result()

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(suite.T(), err)

	return res.StatusCode, data
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
			req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBuffer([]byte(tc.data)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			res := w.Result()

			defer res.Body.Close()
			require.Equalf(
				t, tc.code, res.StatusCode,
				"Error! Expected code: %d, but got %d\n", tc.code, res.StatusCode,
			)

			data, err := ioutil.ReadAll(res.Body)

			require.NoErrorf(t, err, "Error while ReadAll %v", err)

			if tc.errMsg != "" {
				resp := handler.ErrorResponse{}
				err = json.Unmarshal(data, &resp)
				require.Equal(t, resp.Status, "error")
				require.NoErrorf(t, err, "Error while Unmarshal %v", err)
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
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/confirm/%s", tc.link), nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			res := w.Result()

			defer res.Body.Close()
			require.Equalf(
				t, tc.code, res.StatusCode,
				"Error! Expected code: %d, but got %d\n", tc.code, res.StatusCode,
			)

			data, err := ioutil.ReadAll(res.Body)

			require.NoErrorf(t, err, "Error while ReadAll %v", err)

			if tc.errMsg != "" {
				resp := handler.ErrorResponse{}
				err = json.Unmarshal(data, &resp)
				require.Equal(t, resp.Status, "error")
				require.NoErrorf(t, err, "Error while Unmarshal %v", err)
				require.Equal(t, tc.errMsg, resp.Message)
			}
		})
	}
}

func (suite *TestingSuite) TestSignIn() {

	input := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, maxLoggedInUser.Email, defaultPassword)

	for i := 0; i < maxLoggenInCount; i++ {
		code, _ := suite.makeRequest(http.MethodPost, "/auth/sign-in", bytes.NewBuffer([]byte(input)))
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
		code, data := suite.makeRequest(http.MethodPost, "/auth/sign-in", bytes.NewBuffer([]byte(tc.input)))
		require.Equal(suite.T(), tc.code, code)
		if tc.errMsg != "" && tc.code != http.StatusOK {
			resp := handler.ErrorResponse{}
			err := json.Unmarshal(data, &resp)
			require.NoError(suite.T(), err)
			require.Equal(suite.T(), "error", resp.Status)
			require.Equal(suite.T(), tc.errMsg, resp.Message)

		}
	}
}
