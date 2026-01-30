package handlers

import (
	"net/http"
	"time"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"
	"Wrk_Api/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateRetroItemRequest struct {
	SprintID string `json:"sprintId" binding:"required"`
	Type     string `json:"type" binding:"required"` // GOOD, BAD, ACTION
	Content  string `json:"content" binding:"required"`
	UserID   string `json:"userId" binding:"required"`
}

func GetSprintRetrospective(c *gin.Context) {
	sprintID := c.Param("sprintId")
	var items []models.RetrospectiveItem
	
	if result := database.DB.Preload("User").Where("sprint_id = ?", sprintID).Find(&items); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener retrospectiva"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func CreateRetrospectiveItem(c *gin.Context) {
	var req CreateRetroItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := models.RetrospectiveItem{
		ID:        utils.GenerateCUID(),
		SprintID:  req.SprintID,
		Type:      req.Type,
		Content:   req.Content,
		UserID:    req.UserID,
		CreatedAt: time.Now(),
	}

	if result := database.DB.Create(&item); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear item"})
		return
	}
	
	database.DB.Preload("User").First(&item, "id = ?", item.ID)
	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func DeleteRetrospectiveItem(c *gin.Context) {
	id := c.Param("id")
	if result := database.DB.Delete(&models.RetrospectiveItem{}, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item eliminado"})
}
