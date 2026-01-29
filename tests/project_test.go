package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"
	"Wrk_Api/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// Helper to generate a valid JWT for tests
// We need to match the secret used in AuthMiddleware
func generateTestToken(userID, email, role string) string {
	// Must match the fallback or env var in middleware/auth.go
	secret := []byte("default_secret_key")
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"email":  email,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, _ := token.SignedString(secret)
	return tokenString
}

func TestProjectHandlers(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	// Create user
	user := models.User{
		ID:       "owner-123",
		Name:     "Owner User",
		Email:    "owner@example.com",
		Password: "hashedpassword",
		Role:     "SCRUM_MASTER",
		Active:   true,
	}
	database.DB.Create(&user)
	token := generateTestToken(user.ID, user.Email, user.Role)
	authHeader := "Bearer " + token

	var createdProjectID string

	t.Run("CreateProject", func(t *testing.T) {
		body := map[string]interface{}{
			"name":        "New Project",
			"description": "Test Description",
			"ownerId":     user.ID,
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/projects/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"]
		
		createdProjectID = data["ID"].(string)
		assert.Equal(t, "New Project", data["Name"])
		assert.Equal(t, user.ID, data["OwnerID"])
	})

	t.Run("GetProject", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/projects/"+createdProjectID, nil)
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp map[string]map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, createdProjectID, resp["data"]["ID"])
	})

	t.Run("UpdateProject", func(t *testing.T) {
		body := map[string]interface{}{
			"name": "Updated Project Name",
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/projects/"+createdProjectID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify update in DB
		var project models.Project
		database.DB.First(&project, "id = ?", createdProjectID)
		assert.Equal(t, "Updated Project Name", project.Name)
	})

	t.Run("DeleteProject", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/projects/"+createdProjectID, nil)
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify deletion
		var count int64
		database.DB.Model(&models.Project{}).Where("id = ?", createdProjectID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}
