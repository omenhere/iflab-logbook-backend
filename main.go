package main

import (
	"auth-backend/controllers"
	"auth-backend/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"	
)

func main() {
	
	utils.InitializeSupabase()

	// Set up Gin router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Ganti dengan origin frontend Anda
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Start server
	r.Run(":8080")
}
