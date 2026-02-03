package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/meshackkirwachumba/golang-postgres-todos/internal/config"
	"github.com/meshackkirwachumba/golang-postgres-todos/internal/database"
)

func main() {
	//Load environment variables from .env file if it exists
	env_configurations, err := config.LoadEnvironmentalVariables()
	if err != nil {
		log.Fatal("Failed to load environment variables: " + err.Error())
	}
	log.Println("Environment variables loaded successfully:")

	//connect to the database
	connectionPool, err := database.ConnectToDatabase(env_configurations.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database: " + err.Error())
	}
	log.Println("Database connection established successfully")

	defer connectionPool.Close()

	// Initialize Gin router
	var router *gin.Engine = gin.Default()

	router.SetTrustedProxies(nil)
	

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message":"Todo api is running",
			"status":"success",
			"database": "connected",	
		})
	})

	router.Run(":" + env_configurations.Port)
}