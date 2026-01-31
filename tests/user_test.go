package tests

import (
	"bytes"
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
		Email:    "test-users@example.com", 
		Password: "hashedpassword",
		Role:     "ADMIN", // Needed for admin actions if protected (mock check)
		Active:   true,
	}
	database.DB.Create(&user)

	// Mock Auth Token
	token := generateTestToken(user.ID, user.Email, user.Role)
	authHeader := "Bearer " + token

	t.Run("GetAllUsers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/users", nil)
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("CreateUser_Admin", func(t *testing.T) {
		body := map[string]string{
			"name":     "New Admin",
			"email":    "admin@new.com",
			"password": "password123",
			"role":     "ADMIN",
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/users/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		body := map[string]string{
			"name": "Updated Name",
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/users/user-123", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Updated Name", resp["data"]["Name"])
	})

	t.Run("DeleteUser", func(t *testing.T) {
		// Create temp user to delete
		tempUser := models.User{ID: "del-user", Name: "Del", Email: "del@test.com"}
		database.DB.Create(&tempUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/users/del-user", nil)
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
