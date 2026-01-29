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

func TestProjectMembers(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	// Create owner and another user
	owner := models.User{ID: "owner-1", Name: "Owner", Email: "owner@test.com", Role: "SCRUM_MASTER"}
	member := models.User{ID: "user-2", Name: "Member", Email: "member@test.com", Role: "TEAM_DEVELOPER"}
	database.DB.Create(&owner)
	database.DB.Create(&member)

	// Create Project
	project := models.Project{ID: "proj-1", Name: "Test Project", OwnerID: owner.ID}
	database.DB.Create(&project)

	token := generateTestToken(owner.ID, owner.Email, owner.Role)
	authHeader := "Bearer " + token

	t.Run("AddMember", func(t *testing.T) {
		body := map[string]string{
			"userId": member.ID,
			"role":   "TEAM_DEVELOPER",
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/projects/proj-1/members", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Verify membership
		var pm models.ProjectMember
		result := database.DB.Where("project_id = ? AND user_id = ?", project.ID, member.ID).First(&pm)
		assert.NoError(t, result.Error)
		assert.Equal(t, "TEAM_DEVELOPER", pm.Role)
	})

	t.Run("RemoveMember", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/projects/proj-1/members/user-2", nil)
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify removal
		var count int64
		database.DB.Model(&models.ProjectMember{}).Where("project_id = ? AND user_id = ?", project.ID, member.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}
