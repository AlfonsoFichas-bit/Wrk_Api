package handlers

import (
	"net/http"
	"time"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"
	"Wrk_Api/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

type CreateDMRequest struct {
	TargetUserID string `json:"targetUserId" binding:"required"`
}

// GET /:projectId/messages
func GetProjectMessages(c *gin.Context) {
	projectID := c.Param("projectId")
	
	// Optional: Check if user is member of project? 
	// For now, mirroring previous logic but with stricter ID checks if needed.

	var chat models.Chat
	err := database.DB.Preload("Messages.User").First(&chat, "project_id = ?", projectID).Error

	if err != nil {
		// Create if not exists
		chat = models.Chat{
			ID:        utils.GenerateCUID(),
			ProjectID: &projectID,
			Type:      "PROJECT",
		}
		if createErr := database.DB.Create(&chat).Error; createErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear chat"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": chat.Messages})
}

// POST /:projectId/messages
func SendProjectMessage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectID := c.Param("projectId")
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure chat exists
	var chat models.Chat
	err := database.DB.First(&chat, "project_id = ?", projectID).Error
	if err != nil {
		chat = models.Chat{
			ID:        utils.GenerateCUID(),
			ProjectID: &projectID,
			Type:      "PROJECT",
		}
		database.DB.Create(&chat)
	}

	message := models.Message{
		ID:        utils.GenerateCUID(),
		ChatID:    chat.ID,
		UserID:    userID.(string),
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar mensaje"})
		return
	}

	// Re-fetch to include user
	database.DB.Preload("User").First(&message, "id = ?", message.ID)
	c.JSON(http.StatusCreated, gin.H{"data": message})
}

// POST /direct
func CreateOrGetDirectChat(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUserID := userID.(string)

	var req CreateDMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find existing chat with both participants
	var userChats []models.Chat
	database.DB.Joins("JOIN chat_participants cp ON cp.chat_id = chats.id").
		Where("chats.type = ? AND cp.user_id = ?", "DIRECT", currentUserID).
		Preload("Participants").
		Find(&userChats)

	for _, chat := range userChats {
		for _, p := range chat.Participants {
			if p.UserID == req.TargetUserID {
				// Found existing
				c.JSON(http.StatusOK, gin.H{"data": chat})
				return
			}
		}
	}

	// Create new
	chatID := utils.GenerateCUID()
	chat := models.Chat{
		ID:   chatID,
		Type: "DIRECT",
	}
	
	tx := database.DB.Begin()
	if err := tx.Create(&chat).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating chat"})
		return
	}

	p1 := models.ChatParticipant{ChatID: chatID, UserID: currentUserID}
	p2 := models.ChatParticipant{ChatID: chatID, UserID: req.TargetUserID}

	if err := tx.Create(&p1).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding participant 1"})
		return
	}
	if err := tx.Create(&p2).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding participant 2"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": chat})
}

// GET /user/:userId/all
func GetDirectChats(c *gin.Context) {
	userID := c.Param("userId")
	
	// Security check: User can only see their own chats
	authUserID, _ := c.Get("userID")
	if authUserID.(string) != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	
	var chats []models.Chat
	// Find chats where user is participant
	database.DB.Joins("JOIN chat_participants cp ON cp.chat_id = chats.id").
		Where("chats.type = ? AND cp.user_id = ?", "DIRECT", userID).
		Preload("Participants.User").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Find(&chats)

	c.JSON(http.StatusOK, gin.H{"data": chats})
}

// GET /conversation/:chatId/messages
func GetConversationMessages(c *gin.Context) {
	chatID := c.Param("chatId")
	userID, _ := c.Get("userID")

	// Authorization: Check if user is participant
	var count int64
	database.DB.Model(&models.ChatParticipant{}).
		Where("chat_id = ? AND user_id = ?", chatID, userID.(string)).
		Count(&count)
	
	if count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this conversation"})
		return
	}

	var messages []models.Message
	database.DB.Preload("User").Where("chat_id = ?", chatID).Order("created_at asc").Find(&messages)
	c.JSON(http.StatusOK, gin.H{"data": messages})
}

// POST /conversation/:chatId/messages
func SendConversationMessage(c *gin.Context) {
	chatID := c.Param("chatId")
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUserID := userID.(string)

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify participation
	var count int64
	database.DB.Model(&models.ChatParticipant{}).
		Where("chat_id = ? AND user_id = ?", chatID, currentUserID).
		Count(&count)
	
	if count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	message := models.Message{
		ID:        utils.GenerateCUID(),
		ChatID:    chatID,
		UserID:    currentUserID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar mensaje"})
		return
	}

	// Notifications for DM
	var chat models.Chat
	if database.DB.Preload("Participants").First(&chat, "id = ?", chatID).Error == nil {
		if chat.Type == "DIRECT" {
			// Find sender name for message
			var sender models.User
			database.DB.First(&sender, "id = ?", currentUserID)

			for _, p := range chat.Participants {
				if p.UserID != currentUserID {
					notif := models.Notification{
						ID:        utils.GenerateCUID(),
						UserID:    p.UserID,
						Title:     "Nuevo Mensaje Directo",
						Message:   sender.Name + " te ha enviado un mensaje",
						Type:      "MESSAGE",
						CreatedAt: time.Now(),
					}
					database.DB.Create(&notif)
				}
			}
		}
	}

	database.DB.Preload("User").First(&message, "id = ?", message.ID)
	c.JSON(http.StatusCreated, gin.H{"data": message})
}
