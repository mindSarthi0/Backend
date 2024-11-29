package routers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT secret key
var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

// UserCredentials represents the JSON structure for user authentication
type UserCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Authenticate handles user authentication and generates a JWT token
func Authenticate(c *gin.Context) {
	println("Authenticate")
	var creds UserCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if creds.Username != os.Getenv("SP_ADMIN") || creds.Password != os.Getenv("SP_ADMIN_PASSWORD") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	durationStr := os.Getenv("DURATION_HOURS")

	durationHours, err := strconv.Atoi(durationStr)
	if err != nil {
		fmt.Println("Invalid DURATION_HOURS, defaulting to 72 hours")
		durationHours = 72
	}

	// Add the duration to the current time
	futureTime := time.Now().Add(time.Duration(durationHours) * time.Hour)

	// Convert to Unix timestamp
	unixTime := futureTime.Unix()

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": creds.Username,
		"exp":      unixTime,
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
