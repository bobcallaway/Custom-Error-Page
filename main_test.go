package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCodeHeader(t *testing.T) {
	gin := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Code", "200")
	gin.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Code", "503")
	gin.ServeHTTP(w, req)
	assert.Equal(t, 503, w.Code)
}

func TestServerName(t *testing.T) {
	_ = os.Setenv("SERVER_NAME", "test-server")
	gin := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)
	gin.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "test-server")
}
