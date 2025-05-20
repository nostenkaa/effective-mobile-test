package main

import (
	"effective-mobile-test/logger"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"effective-mobile-test/config"
	_ "effective-mobile-test/docs"
	"effective-mobile-test/handlers"
)

// @title Effective Mobile Test API
// @version 1.0
// @description API for enriching person data with age, gender, and nationality.

// @host localhost:8080
// @BasePath /

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	logger.Init()
	config.InitDB()

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	handlers.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
