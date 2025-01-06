package main

import (
	"auth-backend/controllers"
	"auth-backend/middleware"
	"auth-backend/utils"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	utils.InitializeSupabase()

	// Set up Gin router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://iflab-logbook-backend.onrender.com"}, // Pastikan ini sesuai dengan origin frontend Anda
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, // Izinkan cookie atau kredensial
		MaxAge:           12 * time.Hour,
	}))

	r.Use(func(c *gin.Context) {
		log.Println("CORS middleware invoked for:", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	// Auth routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)

	// User routes with middleware
	userGroup := r.Group("/users")
	userGroup.Use(middleware.AuthMiddleware()) // Middleware untuk memeriksa autentikasi
	{
		userGroup.GET("/profile", controllers.GetUser) // Endpoint GetUser
	}

	// Logbook routes with middleware
	protected := r.Group("/logbooks")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/", controllers.GetLogbooks)   // Get logbook by user
		protected.POST("/", controllers.AddLogbook)  // Add new logbook
	}

	// Start server
	r.Run(":8080")
}

