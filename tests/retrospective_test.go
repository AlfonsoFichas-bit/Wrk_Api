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

func TestRetrospectiveHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	// User & Project & Sprint
	user := models.User{ID: "u1", Name: "User", Email: "test@retro.com"}
	database.DB.Create(&user)
	project := models.Project{ID: "p1", Name: "Retro Project", OwnerID: user.ID}
	database.DB.Create(&project)
	sprint := models.Sprint{ID: "s1", Name: "Sprint 1", ProjectID: project.ID}
	database.DB.Create(&sprint)

	token := generateTestToken(user.ID, user.Email, user.Role)
	authHeader := "Bearer " + token

	t.Run("CreateRetroItem", func(t *testing.T) {
		body := map[string]string{
			"sprintId": sprint.ID,
			"type":     "GOOD",
			"content":  "We did well",
			"userId":   user.ID,
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/retrospectives/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var items []models.RetrospectiveItem
		database.DB.Where("sprint_id = ?", sprint.ID).Find(&items)
		assert.Equal(t, 1, len(items))
		assert.Equal(t, "GOOD", items[0].Type)
	})
}
