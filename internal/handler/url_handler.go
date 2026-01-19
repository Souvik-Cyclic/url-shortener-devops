package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik-cyclic/url-shortener-devops/internal/service"
)

type URLHandler struct {
	svc *service.ShortenerService
}

func NewURLHandler(svc *service.ShortenerService) *URLHandler {
	return &URLHandler{svc: svc}
}

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

func (h *URLHandler) Shorten(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	code := h.svc.Shorten(req.URL)
	shortURL := "http://" + c.Request.Host + "/r/" + code

	c.JSON(http.StatusOK, gin.H{
		"short_url": shortURL,
		"code":      code,
	})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	code := c.Param("code")
	originalURL, exists := h.svc.GetOriginalURL(code)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}

func (h *URLHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}

func (h *URLHandler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
