package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/helpers"
	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/VladimirStepanov/todo-app/internal/models/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSendMailReturnForSignUp(t *testing.T) {
	reqData := `{"email": "test@test.com", "password": "123456789"}`
	usObj := new(mocks.UserService)
	msObj := new(mocks.MailService)
	msObj.On("SendConfirmationsEmail", mock.Anything).Return(errors.New("Send mail error"))
	usObj.On("Create", mock.Anything, mock.Anything).Return(nil, nil)

	handler := New(usObj, msObj, nil, nil, getTestLogger())
	r := handler.InitRoutes(gin.TestMode)
	code, _ := helpers.MakeRequest(
		r,
		t,
		http.MethodPost,
		"/auth/sign-up",
		bytes.NewBuffer([]byte(reqData)),
		nil,
	)

	require.Equal(t, http.StatusInternalServerError, code)
	usObj.AssertExpectations(t)
}

func TestCreateErrorForSignUp(t *testing.T) {
	tests := []struct {
		name   string
		retErr error
		code   int
	}{
		{"User already exists", models.ErrUserAlreadyExists, http.StatusConflict},
		{"Internal server error", errors.New("Internal error"), http.StatusInternalServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reqData := `{"email": "test@test.com", "password": "123456789"}`
			usObj := new(mocks.UserService)
			msObj := new(mocks.MailService)
			msObj.On("SendConfirmationsEmail", mock.Anything).Return(nil)
			usObj.On("Create", mock.Anything, mock.Anything).Return(nil, tc.retErr)

			handler := New(usObj, msObj, nil, nil, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, _ := helpers.MakeRequest(
				r,
				t,
				http.MethodPost,
				"/auth/sign-up",
				bytes.NewBuffer([]byte(reqData)),
				nil,
			)

			require.Equal(t, tc.code, code)
			usObj.AssertExpectations(t)
		})
	}
}

func TestBadContentType(t *testing.T) {
	usObj := new(mocks.UserService)
	usObj.On("Create", mock.Anything, mock.Anything).Return(nil, nil)
	msObj := new(mocks.MailService)
	msObj.On("SendConfirmationsEmail", mock.Anything).Return(nil)

	handler := New(usObj, msObj, nil, nil, getTestLogger())
	r := handler.InitRoutes(gin.TestMode)
	req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()

	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	expResp := map[string]interface{}{
		"error": "/auth/sign-up only accepts Content-Type application/json",
	}

	actResp := map[string]interface{}{}

	data, err := ioutil.ReadAll(res.Body)

	require.NoError(t, err)
	err = json.Unmarshal(data, &actResp)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(expResp, actResp))

	// usObj.AssertExpectations(t)
}

func TestSignUpInput(t *testing.T) {

	cases := []struct {
		name   string
		code   int
		expArg *invalidArgument
		data   string
	}{
		{"Sucess", http.StatusOK, nil, `{"email": "test@test.com", "password": "123456789"}`},
		{
			"Empty email", http.StatusBadRequest,
			&invalidArgument{"Email", "", "required", ""},
			`{"password": "12345678"}`,
		},
		{
			"Bad email", http.StatusBadRequest,
			&invalidArgument{"Email", "test", "email", ""},
			`{"email": "test"}`,
		},
		{
			"Empty password", http.StatusBadRequest,
			&invalidArgument{"Password", "", "required", ""},
			`{"email": "test@mail.ru"}`,
		},
		{
			"Short password", http.StatusBadRequest,
			&invalidArgument{"Password", "123", "gte", "8"},
			`{"email": "test@test.com", "password": "123"}`,
		},
		{
			"Long password", http.StatusBadRequest,
			&invalidArgument{"Password", strings.Repeat("1", 48), "lte", "32"},
			fmt.Sprintf(`{"email": "test@test.com", "password": "%s"}`, strings.Repeat("1", 48)),
		},
		{
			"Empty data", http.StatusBadRequest, nil, "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			usObj := new(mocks.UserService)
			usObj.On("Create", mock.Anything, mock.Anything).Return(nil, nil)
			msObj := new(mocks.MailService)
			msObj.On("SendConfirmationsEmail", mock.Anything).Return(nil)

			handler := New(usObj, msObj, nil, nil, getTestLogger())
			r := handler.InitRoutes(gin.TestMode)
			code, data := helpers.MakeRequest(
				r,
				t,
				http.MethodPost,
				"/auth/sign-up",
				bytes.NewBuffer([]byte(tc.data)),
				nil,
			)

			require.Equal(t, tc.code, code)

			if tc.expArg != nil {
				resp := BindDataError{}
				err := json.Unmarshal(data, &resp)

				require.NoError(t, err)

				require.True(t, reflect.DeepEqual(tc.expArg, &resp.InvalidArgs[0]))
			}

			if tc.code == http.StatusOK {
				usObj.AssertExpectations(t)
			}

		})
	}

}
