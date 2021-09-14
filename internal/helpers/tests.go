package helpers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/stretchr/testify/require"
)

var (
	ExpLists = []*models.List{
		{
			ID:          1,
			Title:       "title#1",
			Description: "description#1",
		},
		{
			ID:          2,
			Title:       "title#2",
			Description: "description#2",
		},
	}
)

var (
	ExpItems = []*models.Item{
		{
			ID:          1,
			ListID:      20,
			Title:       "title#1",
			Description: "description#1",
			Done:        true,
		},
		{
			ID:          2,
			ListID:      10,
			Title:       "title#2",
			Description: "description#2",
			Done:        false,
		},
	}
)

func MakeRequest(router http.Handler, t *testing.T, method, path string, input *bytes.Buffer, headers map[string]string) (int, []byte) {
	req := httptest.NewRequest(method, path, input)
	req.Header.Set("Content-Type", "application/json")

	for h, v := range headers {
		req.Header.Set(h, v)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res := w.Result()

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	return res.StatusCode, data
}
