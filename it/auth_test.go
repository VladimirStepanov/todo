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
