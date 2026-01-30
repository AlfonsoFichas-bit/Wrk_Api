package tests

import (
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

func TestMetricHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	// User, Project, Sprint, Stories
	user := models.User{ID: "u1", Name: "User", Email: "test@metrics.com"}
	database.DB.Create(&user)
	project := models.Project{ID: "p1", Name: "Metric Project", OwnerID: user.ID}
	database.DB.Create(&project)
	
	start := time.Now().Add(-time.Hour * 24 * 5)
	end := time.Now().Add(time.Hour * 24 * 5)
	sprint := models.Sprint{ID: "s1", Name: "Sprint 1", ProjectID: project.ID, StartDate: start, EndDate: end}
	database.DB.Create(&sprint)

	points := 5
	story := models.UserStory{ID: "us1", ProjectID: project.ID, SprintID: &sprint.ID, StoryPoints: &points}
	database.DB.Create(&story)

	token := generateTestToken(user.ID, user.Email, user.Role)
	authHeader := "Bearer " + token

	t.Run("GetBurndown", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/metrics/sprints/"+sprint.ID+"/burndown", nil)
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
