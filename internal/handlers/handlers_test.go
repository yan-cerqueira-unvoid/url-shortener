package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yan-cerqueira-unvoid/url-shortener/internal/mocks"
	"github.com/yan-cerqueira-unvoid/url-shortener/internal/models"
	"github.com/yan-cerqueira-unvoid/url-shortener/internal/parser"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestHomeHandler(t *testing.T) {
	router := setupRouter()
	router.GET("/", HomeHandler())

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]any
	err := json.Unmarshal(resp.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, "URL Shortener API", response["message"])
	assert.Equal(t, "1.0.0", response["version"])
	assert.Contains(t, response, "endpoints")
}

func TestShortenURLHandler_Success(t *testing.T) {
	mockURLService := new(mocks.URLService)
	mockURLParser := new(mocks.URLParser)

	router := setupRouter()
	router.POST("/shorten", ShortenURLHandler(mockURLService, mockURLParser))

	expiresAt := time.Now().Add(24 * time.Hour)
	validURL := "https://example.com"
	normalizedURL := "https://example.com"
	shortCode := "abc123"

	mockURLParser.On("Parse", validURL).Return(&parser.URLParseResult{
		OriginalURL: validURL,
		Normalized:  normalizedURL,
		Domain:      "example.com",
		Path:        "",
		Params:      map[string]string{},
		IsValid:     true,
	}, nil)

	mockURLService.On("ShortenURL", normalizedURL, "").Return(&models.URL{
		OriginalURL: validURL,
		ShortCode:   shortCode,
		ExpiresAt:   expiresAt,
	}, nil)

	requestBody := ShortenURLRequest{
		URL: validURL,
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]any
	err := json.Unmarshal(resp.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, validURL, response["original_url"])
	assert.Equal(t, shortCode, response["short_code"])

	mockURLParser.AssertExpectations(t)
	mockURLService.AssertExpectations(t)
}

func TestShortenURLHandler_CustomCode(t *testing.T) {
	mockURLService := new(mocks.URLService)
	mockURLParser := new(mocks.URLParser)

	router := setupRouter()
	router.POST("/shorten", ShortenURLHandler(mockURLService, mockURLParser))

	expiresAt := time.Now().Add(24 * time.Hour)
	validURL := "https://example.com"
	normalizedURL := "https://example.com"
	customCode := "custom"

	mockURLParser.On("Parse", validURL).Return(&parser.URLParseResult{
		OriginalURL: validURL,
		Normalized:  normalizedURL,
		Domain:      "example.com",
		Path:        "",
		Params:      map[string]string{},
		IsValid:     true,
	}, nil)

	mockURLService.On("ShortenURL", normalizedURL, customCode).Return(&models.URL{
		OriginalURL: validURL,
		ShortCode:   customCode,
		ExpiresAt:   expiresAt,
	}, nil)

	requestBody := ShortenURLRequest{
		URL:        validURL,
		CustomCode: customCode,
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]any
	err := json.Unmarshal(resp.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, validURL, response["original_url"])
	assert.Equal(t, customCode, response["short_code"])

	mockURLParser.AssertExpectations(t)
	mockURLService.AssertExpectations(t)
}

func TestShortenURLHandler_InvalidJSON(t *testing.T) {
	mockURLService := new(mocks.URLService)
	mockURLParser := new(mocks.URLParser)

	router := setupRouter()
	router.POST("/shorten", ShortenURLHandler(mockURLService, mockURLParser))

	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer([]byte(`{invalid json}`)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestShortenURLHandler_URLParseError(t *testing.T) {
	mockURLService := new(mocks.URLService)
	mockURLParser := new(mocks.URLParser)

	router := setupRouter()
	router.POST("/shorten", ShortenURLHandler(mockURLService, mockURLParser))

	invalidURL := "invalid-url"

	mockURLParser.On("Parse", invalidURL).Return(nil, errors.New("invalid URL format"))

	requestBody := ShortenURLRequest{
		URL: invalidURL,
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	mockURLParser.AssertExpectations(t)
}

func TestShortenURLHandler_ShortenURLError(t *testing.T) {
	mockURLService := new(mocks.URLService)
	mockURLParser := new(mocks.URLParser)

	router := setupRouter()
	router.POST("/shorten", ShortenURLHandler(mockURLService, mockURLParser))

	validURL := "https://example.com"
	normalizedURL := "https://example.com"

	mockURLParser.On("Parse", validURL).Return(&parser.URLParseResult{
		OriginalURL: validURL,
		Normalized:  normalizedURL,
		Domain:      "example.com",
		Path:        "",
		Params:      map[string]string{},
		IsValid:     true,
	}, nil)

	mockURLService.On("ShortenURL", normalizedURL, "").Return(nil, errors.New("database error"))

	requestBody := ShortenURLRequest{
		URL: validURL,
	}

	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	mockURLParser.AssertExpectations(t)
	mockURLService.AssertExpectations(t)
}

func TestRedirectHandler_Success(t *testing.T) {
	mockURLService := new(mocks.URLService)

	router := setupRouter()
	router.GET("/:shortCode", RedirectHandler(mockURLService))

	shortCode := "abc123"
	originalURL := "https://example.com"

	mockURLService.On("GetURL", shortCode).Return(&models.URL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
	}, nil)

	req, _ := http.NewRequest("GET", "/"+shortCode, nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusMovedPermanently, resp.Code)
	assert.Equal(t, originalURL, resp.Header().Get("Location"))

	mockURLService.AssertExpectations(t)
}

func TestRedirectHandler_NotFound(t *testing.T) {
	mockURLService := new(mocks.URLService)

	router := setupRouter()
	router.GET("/:shortCode", RedirectHandler(mockURLService))

	shortCode := "nonexistent"

	mockURLService.On("GetURL", shortCode).Return(nil, errors.New("short URL not found"))

	req, _ := http.NewRequest("GET", "/"+shortCode, nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	mockURLService.AssertExpectations(t)
}
