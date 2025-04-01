package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yan-cerqueira-unvoid/url-shortener/internal/handlers"
	"github.com/yan-cerqueira-unvoid/url-shortener/internal/parser"
	"github.com/yan-cerqueira-unvoid/url-shortener/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "url_shortener"
	}

	db := client.Database(dbName)

	router := gin.Default()

	urlService := services.NewURLService(db, &ctx)
	urlParser := parser.NewURLParser()

	router.GET("/", handlers.HomeHandler())
	router.GET("/:shortCode", handlers.RedirectHandler(urlService))

	router.POST("/shorten", handlers.ShortenURLHandler(urlService, urlParser))

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
