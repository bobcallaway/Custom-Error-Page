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
	_ = os.Setenv(ServerNameVar, "test-server")
	gin := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Ingress-Name", "test-ingress")
	gin.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "test-server")

	w = httptest.NewRecorder()
	_ = os.Setenv(ServerNameVar, "")
	gin = setupRouter()
	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Ingress-Name", "test-ingress")
	gin.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "test-ingress")
}

func TestContentType(t *testing.T) {
	gin := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Format", "text/html")
	gin.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "html")
	assert.Equal(t, w.Header().Get(ContentType), "text/html")

	w = httptest.NewRecorder()
	gin = setupRouter()
	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Format", "application/json")
	gin.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "{")
	assert.Equal(t, w.Header().Get(ContentType), "application/json")
}
