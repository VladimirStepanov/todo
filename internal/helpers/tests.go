package helpers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
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
