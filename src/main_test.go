package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Set Gin to release mode for tests
	gin.SetMode(gin.ReleaseMode)

	// Setup test environment
	setupTestEnv()

	// Run tests
	code := m.Run()

	// Cleanup
	cleanupTestEnv()

	os.Exit(code)
}

func setupTestEnv() {
	// Create test data directory
	testDir := "test-data"
	os.MkdirAll(testDir, 0755)

	// Set environment variables
	os.Setenv("PORT", "7837")
	os.Setenv("HOST", "127.0.0.1")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_USERS", "testuser")
	os.Setenv("FILE_STORE_PATH", testDir)
}

func cleanupTestEnv() {
	testDir := "test-data"
	os.RemoveAll(testDir)
}

func generateTestToken(userId string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenString
}

func TestIntegration(t *testing.T) {
	r := setupRouter()

	t.Run("Test /test endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "ok", w.Body.String())
	})

	t.Run("Test unauthorized access", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/sync", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})

	t.Run("Test sync workflow", func(t *testing.T) {
		token := generateTestToken("testuser")

		// Test PUT
		testData := map[string]interface{}{
			"test": "data",
		}
		jsonData, _ := json.Marshal(testData)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/sync", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "ok", w.Body.String())

		// Test GET
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/sync", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var responseData map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &responseData)
		assert.NoError(t, err)
		assert.Equal(t, "data", responseData["test"])
	})
}
