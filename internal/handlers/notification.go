package handlers

import (
	"net/http"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"

	"github.com/gin-gonic/gin"
)

// GET /api/notifications
func GetNotifications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var notifications []models.Notification
	// Return all, or filter by unread? Usually show latest.
	// Let's return latest 50
	if result := database.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(50).Find(&notifications); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener notificaciones"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notifications})
}

// PUT /api/notifications/:id/read
func MarkNotificationRead(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	var notification models.Notification
	if result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&notification); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notificación no encontrada"})
		return
	}

	notification.Read = true
	database.DB.Save(&notification)

	c.JSON(http.StatusOK, gin.H{"data": notification})
}
