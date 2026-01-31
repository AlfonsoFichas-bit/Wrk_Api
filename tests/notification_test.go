package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"
	"Wrk_Api/internal/routes"
	"Wrk_Api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNotificationHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	// User
	user := models.User{ID: "u1", Name: "User", Email: "test@notif.com"}
	database.DB.Create(&user)

	// Create Notification manually (simulating system event)
	notif := models.Notification{
		ID:        utils.GenerateCUID(),
		UserID:    user.ID,
		Title:     "Test Notif",
		Message:   "Something happened",
		Type:      "TEST",
		Read:      false,
		CreatedAt: time.Now(),
	}
	database.DB.Create(&notif)

	token := generateTestToken(user.ID, user.Email, user.Role)
	authHeader := "Bearer " + token

	t.Run("GetNotifications", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/notifications/", nil)
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("MarkRead", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/notifications/"+notif.ID+"/read", nil)
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var n models.Notification
		database.DB.First(&n, "id = ?", notif.ID)
		assert.True(t, n.Read)
	})
}
