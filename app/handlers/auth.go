package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Hardcoded admin credentials - CHANGE THESE IN PRODUCTION!
var (
	adminUsername     = "admin"
	adminPasswordHash = "$2a$12$.m2mfJAZaMaNy8wNyRvxqe9CiZxtXy0D4CGaD6J6kPV1d84Msvn0e"

	// Active tokens with expiration
	activeTokens = make(map[string]time.Time)
	tokenMutex   sync.RWMutex

	// Token validity duration
	tokenDuration = 24 * time.Hour
)

func init() {
	// Allow environment variables to override credentials
	if envUser := os.Getenv("ADMIN_USERNAME"); envUser != "" {
		adminUsername = envUser
	}
	if envPass := os.Getenv("ADMIN_PASSWORD"); envPass != "" {
		// Hash the password from env
		hash, err := bcrypt.GenerateFromPassword([]byte(envPass), 12)
		if err == nil {
			adminPasswordHash = string(hash)
		}
	}
}

// generateToken creates a secure random token
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Add timestamp to make it unique
	timestamp := time.Now().UnixNano()
	data := append(bytes, byte(timestamp))

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// AdminLogin handles admin authentication
func AdminLogin(c *gin.Context) {
	var creds AdminCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Verify username
	if creds.Username != adminUsername {
		// Use constant time comparison to prevent timing attacks
		time.Sleep(100 * time.Millisecond)
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(adminPasswordHash), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// Generate token
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Error generating session",
		})
		return
	}

	// Store token with expiration
	tokenMutex.Lock()
	activeTokens[token] = time.Now().Add(tokenDuration)

	// Cleanup expired tokens
	for t, exp := range activeTokens {
		if time.Now().After(exp) {
			delete(activeTokens, t)
		}
	}
	tokenMutex.Unlock()

	c.JSON(http.StatusOK, LoginResponse{
		Token:   token,
		Message: "Login successful",
	})
}

// AuthMiddleware validates the admin token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "Invalid authorization format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		tokenMutex.RLock()
		expiration, exists := activeTokens[token]
		tokenMutex.RUnlock()

		if !exists || time.Now().After(expiration) {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Logout invalidates a token
func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 {
		token := parts[1]
		tokenMutex.Lock()
		delete(activeTokens, token)
		tokenMutex.Unlock()
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}
