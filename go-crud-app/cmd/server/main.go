package main

import (
	"log"
	"os"
	"time"

	"go-crud-app/internal/database"
	"go-crud-app/internal/handlers"
	"go-crud-app/internal/middleware"
	"go-crud-app/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Database configuration
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "gocrud"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Connect to database
	if err := database.Connect(dbConfig); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// JWT configuration
	jwtConfig := utils.JWTConfig{
		SecretKey:       getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
		ExpirationHours: 24, // 24 hours
	}

	// Initialize Gin router
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{getEnv("CORS_ORIGIN", "*")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rate limiters
	authLimiter := middleware.NewRateLimiter(5, 1*time.Minute)      // 5 requests per minute for auth
	registerLimiter := middleware.NewRateLimiter(3, 1*time.Minute)  // 3 requests per minute for registration
	generalLimiter := middleware.NewRateLimiter(100, 1*time.Minute) // 100 requests per minute for general endpoints

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Authentication routes (with rate limiting)
		auth := api.Group("/auth")
		{
			auth.POST("/register", middleware.RateLimitMiddleware(registerLimiter), handlers.Register(jwtConfig))
			auth.POST("/login", middleware.RateLimitMiddleware(authLimiter), handlers.Login(jwtConfig))
		}

		// Protected user routes (require authentication)
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(jwtConfig.SecretKey))
		users.Use(middleware.RateLimitMiddleware(generalLimiter))
		{
			users.GET("", handlers.GetAllUsers)       // List all users except current user
			users.GET("/me", handlers.GetCurrentUser) // Get current user profile
			users.GET("/:id", handlers.GetUserByID)   // Get user by ID
			users.PUT("/:id", handlers.UpdateUser)    // Update user (own profile only)
			users.DELETE("/:id", handlers.DeleteUser) // Delete user (own profile only)
		}
	}

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
