package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/VladimirStepanov/todo-app/internal/models/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestConfirmHandler(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		code    int
		mockErr error
	}{
		{"Success", "success", http.StatusOK, nil},
		{"Page not found", "not_found", http.StatusNotFound, models.ErrConfirmLinkNotExists},
		{"Internal error", "int_error", http.StatusInternalServerError, fmt.Errorf("unknown error")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			usObj := new(mocks.UserService)
			usObj.On("ConfirmEmail", mock.Anything).Return(tc.mockErr)
			handler := New(usObj, nil, nil, getTestLogger())

			r := handler.InitRoutes(gin.TestMode)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/confirm/%s", tc.link), nil)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			res := w.Result()

			defer res.Body.Close()

			require.Equal(t, tc.code, res.StatusCode)
		})
	}
}
