package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"
	"Wrk_Api/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserHandlers(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	// Create a test user
	user := models.User{
		ID:       "user-123",
		Name:     "Test User",
		Email:    "test-users@example.com", // Unique email for this test
		Password: "hashedpassword",
		Role:     "TEAM_DEVELOPER",
		Active:   true,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Mock Auth Token
	token := generateTestToken(user.ID, user.Email, user.Role)
	authHeader := "Bearer " + token

	t.Run("GetAllUsers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/users", nil)
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var users []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &users)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1)
		
		// Verify we find our user
		found := false
		for _, u := range users {
			if u["id"] == "user-123" {
				found = true
				break
			}
		}
		assert.True(t, found, "User should be in list")
	})

	t.Run("GetUserByID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/users/user-123", nil)
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Test User", resp["name"])
	})

	t.Run("GetUserNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/users/non-existent-id", nil)
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
