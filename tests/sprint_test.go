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
	"github.com/stretchr/testify/assert"
)

func TestSprintHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	// User & Project Setup
	user := models.User{ID: "u1", Name: "Test", Email: "test@sprint.com", Role: "SCRUM_MASTER"}
	database.DB.Create(&user)
	project := models.Project{ID: "p1", Name: "Sprint Project", OwnerID: user.ID}
	database.DB.Create(&project)

	token := generateTestToken(user.ID, user.Email, user.Role)
	authHeader := "Bearer " + token

	var createdSprintID string

	t.Run("CreateSprint", func(t *testing.T) {
		body := map[string]interface{}{
			"name":      "Sprint 1",
			"projectId": project.ID,
			"startDate": time.Now().Format(time.RFC3339),
			"endDate":   time.Now().Add(time.Hour * 24 * 14).Format(time.RFC3339),
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/sprints/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		createdSprintID = resp["data"]["ID"].(string)
		assert.Equal(t, "Sprint 1", resp["data"]["Name"])
	})

	t.Run("AddStoryToSprint", func(t *testing.T) {
		// Create a story first
		story := models.UserStory{
			ID:          "us1",
			Title:       "Story 1",
			Description: "Desc",
			ProjectID:   project.ID,
			Status:      "BACKLOG",
		}
		database.DB.Create(&story)

		body := map[string]string{"userStoryId": story.ID}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/sprints/"+createdSprintID+"/add-story", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Verify story updated
		var updatedStory models.UserStory
		database.DB.First(&updatedStory, "id = ?", story.ID)
		assert.Equal(t, createdSprintID, *updatedStory.SprintID)
	})
}
