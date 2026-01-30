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

func TestChatHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	// User & Project Setup
	user1 := models.User{ID: "u1", Name: "User 1", Email: "u1@chat.com", Role: "TEAM_DEVELOPER"}
	user2 := models.User{ID: "u2", Name: "User 2", Email: "u2@chat.com", Role: "TEAM_DEVELOPER"}
	user3 := models.User{ID: "u3", Name: "User 3", Email: "u3@chat.com", Role: "TEAM_DEVELOPER"}
	database.DB.Create(&user1)
	database.DB.Create(&user2)
	database.DB.Create(&user3)
	
	project := models.Project{ID: "p1", Name: "Chat Project", OwnerID: user1.ID}
	database.DB.Create(&project)

	token1 := generateTestToken(user1.ID, user1.Email, user1.Role)
	authHeader1 := "Bearer " + token1
	
	token3 := generateTestToken(user3.ID, user3.Email, user3.Role)
	authHeader3 := "Bearer " + token3

	t.Run("ProjectChat_AutoCreate", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/chat/"+project.ID+"/messages", nil)
		req.Header.Set("Authorization", authHeader1)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Should have created a chat and returned empty messages
		var resp map[string][]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, len(resp["data"]))

		var chat models.Chat
		database.DB.Where("project_id = ?", project.ID).First(&chat)
		assert.NotEmpty(t, chat.ID)
	})

	t.Run("ProjectChat_SendMessage", func(t *testing.T) {
		body := map[string]string{
			"content": "Hello Project",
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/chat/"+project.ID+"/messages", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader1)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Hello Project", resp["data"]["Content"])
		assert.Equal(t, user1.ID, resp["data"]["UserID"])
	})

	var dmChatID string

	t.Run("DirectChat_Create", func(t *testing.T) {
		body := map[string]string{
			"targetUserId": user2.ID,
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/chat/direct", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader1)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		dmChatID = resp["data"]["ID"].(string)
		assert.Equal(t, "DIRECT", resp["data"]["Type"])
	})

	t.Run("DirectChat_SendMessage_Notification", func(t *testing.T) {
		body := map[string]string{
			"content": "Hello DM",
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/chat/conversation/"+dmChatID+"/messages", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader1)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Check Notification for User 2
		var notif models.Notification
		database.DB.Where("user_id = ? AND type = ?", user2.ID, "MESSAGE").First(&notif)
		assert.NotEmpty(t, notif.ID)
		assert.Contains(t, notif.Message, "User 1")
	})
	
	t.Run("DirectChat_Security_IDOR", func(t *testing.T) {
		// User 3 tries to access User 1 & 2's chat
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/chat/conversation/"+dmChatID+"/messages", nil)
		req.Header.Set("Authorization", authHeader3)
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusForbidden, w.Code)
		
		// User 3 tries to send message to User 1 & 2's chat
		body := map[string]string{
			"content": "Hacking in",
		}
		jsonBody, _ := json.Marshal(body)
		
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("POST", "/api/chat/conversation/"+dmChatID+"/messages", bytes.NewBuffer(jsonBody))
		req2.Header.Set("Authorization", authHeader3)
		r.ServeHTTP(w2, req2)
		
		assert.Equal(t, http.StatusForbidden, w2.Code)
	})
}
