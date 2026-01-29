package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Wrk_Api/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthRegisterAndLogin(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	// Test Register
	registerBody := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
		"role":     "TEAM_DEVELOPER",
	}
	body, _ := json.Marshal(registerBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Usuario registrado exitosamente", response["message"])

	// Test Login
	loginBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ = json.Marshal(loginBody)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "token")
	assert.Equal(t, "Inicio de sesión exitoso", response["message"])
}
