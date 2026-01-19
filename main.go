package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/souvik-cyclic/url-shortener-devops/internal/handler"
	"github.com/souvik-cyclic/url-shortener-devops/internal/service"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// Initialize service
	shortenerService := service.NewShortenerService()
	urlHandler := handler.NewURLHandler(shortenerService)

	// Initialize Gin router
	r := gin.Default()

	// Define routes
	r.POST("/shorten", urlHandler.Shorten)
	r.GET("/r/:code", urlHandler.Redirect)
	r.GET("/health", urlHandler.Health)
	r.GET("/", urlHandler.Hello)

	// Run server
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
