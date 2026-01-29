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

func TestTaskHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	// User & Project Setup
	user := models.User{ID: "u1", Name: "Test", Email: "test@task.com", Role: "SCRUM_MASTER"}
	assignee := models.User{ID: "u2", Name: "Dev", Email: "dev@task.com", Role: "TEAM_DEVELOPER"}
	database.DB.Create(&user)
	database.DB.Create(&assignee)
	
	project := models.Project{ID: "p1", Name: "Task Project", OwnerID: user.ID}
	database.DB.Create(&project)

	token := generateTestToken(user.ID, user.Email, user.Role)
	authHeader := "Bearer " + token

	var taskID string

	t.Run("CreateTask", func(t *testing.T) {
		body := map[string]interface{}{
			"title":      "New Task",
			"projectId":  project.ID,
			"assigneeId": assignee.ID,
			"deadline":   time.Now().Add(time.Hour * 24).Format(time.RFC3339),
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/tasks/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		taskID = resp["data"]["ID"].(string)
		
		// Check Notification Created
		var notif models.Notification
		database.DB.Where("user_id = ? AND type = ?", assignee.ID, "TASK_ASSIGNED").First(&notif)
		assert.NotEmpty(t, notif.ID)
	})

	t.Run("FilterTasks", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/tasks/?assigneeId="+assignee.ID, nil)
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string][]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp["data"])
	})

	t.Run("UpdateTaskStatus", func(t *testing.T) {
		body := map[string]interface{}{
			"status": "COMPLETED",
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/tasks/"+taskID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var task models.Task
		database.DB.First(&task, "id = ?", taskID)
		assert.Equal(t, "COMPLETED", task.Status)
		assert.NotNil(t, task.CompletedAt)
	})

	t.Run("EvaluateTask", func(t *testing.T) {
		body := map[string]interface{}{
			"score":       85,
			"feedback":    "Good job",
			"evaluatorId": user.ID,
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/tasks/"+taskID+"/evaluate", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Check Evaluation in DB
		var eval models.Evaluation
		database.DB.Where("task_id = ?", taskID).First(&eval)
		assert.Equal(t, 85, *eval.Score)

		// Check Notification for Assignee
		var notif models.Notification
		database.DB.Where("user_id = ? AND type = ?", assignee.ID, "EVALUATION_COMPLETED").First(&notif)
		assert.NotEmpty(t, notif.ID)
	})
}
