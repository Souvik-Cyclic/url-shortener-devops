package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/souvik-cyclic/url-shortener-devops/internal/handler"
	"github.com/souvik-cyclic/url-shortener-devops/internal/service"
	"github.com/stretchr/testify/assert"
)

func getTestURL() string {
	return "https://google.com"
}

func setupRouter() *gin.Engine {
	svc := service.NewShortenerService()
	h := handler.NewURLHandler(svc)
	r := gin.Default()
	r.POST("/shorten", h.Shorten)
	r.GET("/r/:code", h.Redirect)
	r.GET("/health", h.Health)
	return r
}

func TestHealthCheck(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status": "up"}`, w.Body.String())
}

func TestShortenURL(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	testURL := getTestURL()
	reqBody := map[string]string{"url": testURL}
	jsonBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp["short_url"])
	assert.NotEmpty(t, resp["code"])
}

func TestRedirectURL(t *testing.T) {
	r := setupRouter()

	// Shorten a URL first
	w := httptest.NewRecorder()
	testURL := getTestURL()
	reqBody := map[string]string{"url": testURL}
	jsonBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonBytes))
	r.ServeHTTP(w, req)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	code := resp["code"]

	// Now try to redirect
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/r/"+code, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, testURL, w.Result().Header.Get("Location"))
}

func TestRedirectNotFound(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/r/nonexistent", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSameURL(t *testing.T) {
	r := setupRouter()

	// Shorten first time
	w1 := httptest.NewRecorder()
	testURL := getTestURL()
	reqBody := map[string]string{"url": testURL}
	jsonBytes, _ := json.Marshal(reqBody)
	req1, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonBytes))
	r.ServeHTTP(w1, req1)

	var resp1 map[string]string
	err := json.Unmarshal(w1.Body.Bytes(), &resp1)
	assert.NoError(t, err)
	code1 := resp1["code"]

	// Shorten second time
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonBytes))
	r.ServeHTTP(w2, req2)

	var resp2 map[string]string
	err = json.Unmarshal(w2.Body.Bytes(), &resp2)
	assert.NoError(t, err)
	code2 := resp2["code"]

	assert.Equal(t, code1, code2, "Shortening the same URL should return the same code")
}

func TestWhitespaceURL(t *testing.T) {
	r := setupRouter()

	// Shorten clean URL
	w1 := httptest.NewRecorder()
	testURL := getTestURL()
	reqBody1 := map[string]string{"url": testURL}
	jsonBytes1, _ := json.Marshal(reqBody1)
	req1, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonBytes1))
	r.ServeHTTP(w1, req1)

	var resp1 map[string]string
	err := json.Unmarshal(w1.Body.Bytes(), &resp1)
	assert.NoError(t, err)
	code1 := resp1["code"]

	// Shorten URL with whitespace
	w2 := httptest.NewRecorder()
	reqBody2 := map[string]string{"url": "   " + testURL + "   "}
	jsonBytes2, _ := json.Marshal(reqBody2)
	req2, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonBytes2))
	r.ServeHTTP(w2, req2)

	var resp2 map[string]string
	err = json.Unmarshal(w2.Body.Bytes(), &resp2)
	assert.NoError(t, err)
	code2 := resp2["code"]

	assert.Equal(t, code1, code2, "Shortening URL with whitespace should return the same code")
}
