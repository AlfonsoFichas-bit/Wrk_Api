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

func TestUserStoryHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	// User & Project Setup
	user := models.User{ID: "u1", Name: "Test", Email: "test@story.com", Role: "SCRUM_MASTER"}
	assignee := models.User{ID: "u2", Name: "Dev", Email: "dev@story.com", Role: "TEAM_DEVELOPER"}
	database.DB.Create(&user)
	database.DB.Create(&assignee)
	
	project := models.Project{ID: "p1", Name: "Story Project", OwnerID: user.ID}
	database.DB.Create(&project)

	token := generateTestToken(user.ID, user.Email, user.Role)
	authHeader := "Bearer " + token

	var storyID string

	t.Run("CreateUserStory", func(t *testing.T) {
		body := map[string]interface{}{
			"title":       "New Story",
			"description": "As a user...",
			"projectId":   project.ID,
			"priority":    "HIGH",
			"storyPoints": 5,
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/user-stories/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		storyID = resp["data"]["ID"].(string)
		assert.Equal(t, "New Story", resp["data"]["Title"])
	})

	t.Run("UpdateStoryStatusAndAssign", func(t *testing.T) {
		body := map[string]interface{}{
			"status":     "COMPLETED",
			"assigneeId": assignee.ID,
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/user-stories/"+storyID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify updates
		var story models.UserStory
		database.DB.First(&story, "id = ?", storyID)
		assert.Equal(t, "COMPLETED", story.Status)
		assert.NotNil(t, story.CompletedAt)
		assert.Equal(t, assignee.ID, *story.AssigneeID)

		// Verify notification generated
		var notif models.Notification
		database.DB.Where("user_id = ? AND type = ?", assignee.ID, "TASK_ASSIGNED").First(&notif)
		assert.NotEmpty(t, notif.ID)
	})
}
