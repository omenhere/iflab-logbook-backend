package main

import (
	"auth-backend/controllers"
	"auth-backend/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	
	utils.InitializeSupabase()

	// Set up Gin router
	r := gin.Default()

	
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Start server
	r.Run(":8080")
}
