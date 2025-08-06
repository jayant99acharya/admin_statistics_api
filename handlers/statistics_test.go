package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new handler
	handler := NewStatisticsHandler()

	// Create a new Gin router
	router := gin.New()
	router.GET("/health", handler.HealthCheck)

	// Create a test request
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
	assert.Contains(t, w.Body.String(), "Admin Statistics API")
}

func TestAuthenticationRequired(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new handler
	handler := NewStatisticsHandler()

	// Create a new Gin router with auth middleware
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Simple auth check for testing
		auth := c.GetHeader("Authorization")
		if auth != "test-token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	})
	router.GET("/gross_gaming_rev", handler.GetGrossGamingRevenue)

	// Test without auth header
	req, err := http.NewRequest("GET", "/gross_gaming_rev?from=2024-01-01&to=2024-01-31", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Test with auth header
	req.Header.Set("Authorization", "test-token")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Note: This will fail with database connection error in test environment
	// but it shows the auth middleware is working
	assert.NotEqual(t, http.StatusUnauthorized, w.Code)
}

func TestInvalidDateFormat(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new handler
	handler := NewStatisticsHandler()

	// Create a new Gin router
	router := gin.New()
	router.GET("/gross_gaming_rev", handler.GetGrossGamingRevenue)

	// Test with invalid date format
	req, err := http.NewRequest("GET", "/gross_gaming_rev?from=invalid-date&to=2024-01-31", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid date parameters")
}

func TestMissingDateParameters(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new handler
	handler := NewStatisticsHandler()

	// Create a new Gin router
	router := gin.New()
	router.GET("/gross_gaming_rev", handler.GetGrossGamingRevenue)

	// Test without date parameters
	req, err := http.NewRequest("GET", "/gross_gaming_rev", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}