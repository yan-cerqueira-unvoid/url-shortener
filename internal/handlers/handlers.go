package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yan-cerqueira-unvoid/url-shortener/internal/models"
	"github.com/yan-cerqueira-unvoid/url-shortener/internal/parser"
)

type ShortenURLRequest struct {
	URL        string `json:"url" binding:"required"`
	CustomCode string `json:"custom_code,omitempty"`
}

type URLServiceInterface interface {
	ShortenURL(originalURL string, customCode string) (*models.URL, error)
	GetURL(shortCode string) (*models.URL, error)
}

type URLParserInterface interface {
	Parse(rawURL string) (*parser.URLParseResult, error)
}

func HomeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "URL Shortener API",
			"version": "1.0.0",
			"endpoints": []string{
				"POST /shorten",
				"GET /:shortCode",
			},
		})
	}
}

func ShortenURLHandler(urlService URLServiceInterface, urlParser URLParserInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request ShortenURLRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		parseResult, err := urlParser.Parse(request.URL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		url, err := urlService.ShortenURL(parseResult.Normalized, request.CustomCode)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"original_url": url.OriginalURL,
			"short_code":   url.ShortCode,
			"short_url":    c.Request.Host + "/" + url.ShortCode,
			"expires_at":   url.ExpiresAt,
		})
	}
}

func RedirectHandler(urlService URLServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")

		url, err := urlService.GetURL(shortCode)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
	}
}
