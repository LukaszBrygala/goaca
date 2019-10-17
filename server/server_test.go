package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestServer() Server {
	s := New()
	s.routes()
	return s
}

func TestServer(t *testing.T) {
	s := newTestServer()

	req, _ := http.NewRequest("GET", "/", nil)

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, `{"name":"goaca-server","version":"0.0.1"}`, string(rec.Body.Bytes()))
}
