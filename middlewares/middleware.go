package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// JWT secret key
var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

// Rate Limiting Middleware
func RateLimitingMiddleware() gin.HandlerFunc {
	limiter := rate.NewLimiter(1, 5) // 1 request per second with a burst of 5
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}
		c.Next()
	}
}

// CORSMiddleware sets up CORS configuration
func CORSMiddleware() gin.HandlerFunc {
	// Get allowed origins from environment variable or use default
	allowOrigins := os.Getenv("WEBAPP_DOMAIN")
	if allowOrigins == "" {
		allowOrigins = "http://localhost:3000" // Default value if not set in the environment
	}

	// Split the comma-separated allowed origins into a slice
	allowOriginsSlice := strings.Split(allowOrigins, ",")

	// Return CORS configuration using the gin-cors middleware
	return cors.New(cors.Config{
		AllowOrigins:     allowOriginsSlice,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// JWTAuthMiddleware verifies the JWT token
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request is for /auth or /health with POST method

		// Get excluded paths from environment variable
		excludedPaths := os.Getenv("EXCLUDED_PATHS")
		if excludedPaths == "" {
			// Default paths if not set in .env
			excludedPaths = ""
		}

		// Split the comma-separated excluded paths into a slice
		excludedPathsSlice := strings.Split(excludedPaths, ",")

		// Check if the request path and method are in the excluded paths
		requestMethod := c.Request.Method
		requestPath := c.Request.URL.Path

		var excludedPath = false
		for _, path := range excludedPathsSlice {
			// Split the path into method and path components
			parts := strings.SplitN(path, " ", 2)

			// If there are not exactly two parts (method and path), skip this entry
			if len(parts) != 2 {
				continue
			}

			method := parts[0]
			path := parts[1]

			// If the request method and path match the excluded ones, proceed with the request
			if strings.EqualFold(requestMethod, method) && strings.Contains(requestPath, path) {
				c.Next() // Skip further checks and process the request
				excludedPath = true
				return
			}
		}

		if !excludedPath {
			// Retrieve Authorization header
			authHeader := c.GetHeader("Authorization")

			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
				c.Abort()
				return
			}

			log.Println("AUTH HEADER:" + authHeader)

			// Check if the token is in the correct format
			if !strings.HasPrefix(authHeader, "Bearer ") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
				c.Abort()
				return
			}

			// Extract the token from the "Bearer " prefix
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			log.Println("TOKEN STRING:" + tokenString)
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Ensure the token signing method is correct
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return jwtSecret, nil
			})

			// Handle errors or invalid token
			if err != nil || !token.Valid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}

			// Proceed with the request
			c.Next()
		}
	}

}

// Error Handling Middleware
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			log.Printf("Error: %v", c.Errors.String())
			c.JSON(c.Writer.Status(), gin.H{"error": "An error occurred, please try again later"})
		}
	}
}

// Input Validation Middleware
func InputValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if strings.Contains(c.Request.URL.Path, "/questions") && c.Request.Method == "POST" {
			// Add input validation logic here
		}
	}
}
