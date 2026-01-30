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

func TestRubricHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	// User & Project
	user := models.User{ID: "u1", Name: "User", Email: "test@rubric.com"}
	database.DB.Create(&user)
	project := models.Project{ID: "p1", Name: "Rubric Project", OwnerID: user.ID}
	database.DB.Create(&project)

	token := generateTestToken(user.ID, user.Email, user.Role)
	authHeader := "Bearer " + token

	t.Run("CreateRubric", func(t *testing.T) {
		body := map[string]interface{}{
			"projectId":   project.ID,
			"name":        "Code Review",
			"description": "Standard rubric",
			"criteria": []map[string]interface{}{
				{"name": "Code Style", "maxScore": 10, "weight": 1},
				{"name": "Logic", "maxScore": 20, "weight": 2},
			},
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/rubrics/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Verify DB
		var rubric models.Rubric
		database.DB.Preload("Criteria").Where("name = ?", "Code Review").First(&rubric)
		assert.NotEmpty(t, rubric.ID)
		assert.Equal(t, 2, len(rubric.Criteria))
	})
}
