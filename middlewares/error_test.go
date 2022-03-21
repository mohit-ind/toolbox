package middlewares

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	logrusTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestRespondWithError(t *testing.T) {
	req := require.New(t)

	// hook := logrusTest.NewGlobal()

	w := httptest.NewRecorder()

	RespondWithError(w, http.StatusInternalServerError, "Internal Server Error", "500")

	resp := w.Result()

	body, err := ioutil.ReadAll(resp.Body)

	req.NoError(err, "The response body should be readable")

	req.Equal(
		http.StatusInternalServerError,
		resp.StatusCode,
		"The response status should be Internal Server Error 500",
	)

	req.Equal(
		"application/json",
		resp.Header.Get("Content-Type"),
		"The response's content type should be application/json",
	)

	req.Equal(
		[]byte(`{"error":{"code":500,"message":"Internal Server Error","messageCode":"500"}}`),
		body,
		`The response body should be "error":{"code":500,"message":"Internal Server Error","messageCode":"500"}}`,
	)
}

type FailWriter struct{}

func (fw *FailWriter) Header() http.Header {
	return http.Header{}
}

func (fw *FailWriter) Write(body []byte) (int, error) {
	return 0, errors.New("Failed to write")
}

func (fw *FailWriter) WriteHeader(statusCode int) {}

func TestRespondWithError_fail_to_write(t *testing.T) {
	req := require.New(t)

	hook := logrusTest.NewGlobal()

	RespondWithError(&FailWriter{}, http.StatusInternalServerError, "Internal Server Error", "500")

	req.Equal("Failed to write HTTP response", hook.LastEntry().Message, "An error should have been logged")
}
