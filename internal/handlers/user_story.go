package handlers

import (
	"net/http"
	"time"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"

	"github.com/gin-gonic/gin"
)

type CreateUserStoryRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Acceptance  string `json:"acceptance"`
	ProjectID   string `json:"projectId" binding:"required"`
	AssigneeID  string `json:"assigneeId"`
	Priority    string `json:"priority"`
	StoryPoints int    `json:"storyPoints"`
}

type UpdateUserStoryRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Acceptance  string `json:"acceptance"`
	Priority    string `json:"priority"`
	StoryPoints *int   `json:"storyPoints"`
	AssigneeID  string `json:"assigneeId"`
	SprintID    string `json:"sprintId"`
	Status      string `json:"status"`
}

func GetAllUserStories(c *gin.Context) {
	var stories []models.UserStory
	if result := database.DB.Preload("Project").Preload("Assignee").Find(&stories); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener user stories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stories})
}

func GetUserStory(c *gin.Context) {
	id := c.Param("id")
	var story models.UserStory
	if result := database.DB.Preload("Project").Preload("Assignee").Preload("Tasks").First(&story, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User story no encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": story})
}

func CreateUserStory(c *gin.Context) {
	var req CreateUserStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	priority := "MEDIUM"
	if req.Priority != "" {
		priority = req.Priority
	}

	story := models.UserStory{
		ID:          generateCUIDStory(),
		Title:       req.Title,
		Description: req.Description,
		Acceptance:  &req.Acceptance,
		ProjectID:   req.ProjectID,
		AssigneeID:  nil,
		Priority:    priority,
		StoryPoints: &req.StoryPoints,
		Status:      "BACKLOG",
	}

	if req.AssigneeID != "" {
		story.AssigneeID = &req.AssigneeID
	}

	if result := database.DB.Create(&story); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear user story", "details": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": story})
}

func UpdateUserStory(c *gin.Context) {
	id := c.Param("id")
	var req UpdateUserStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var story models.UserStory
	if result := database.DB.Preload("Project").First(&story, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User story no encontrado"})
		return
	}

	previousAssigneeID := story.AssigneeID

	if req.Title != "" {
		story.Title = req.Title
	}
	if req.Description != "" {
		story.Description = req.Description
	}
	if req.Acceptance != "" {
		story.Acceptance = &req.Acceptance
	}
	if req.Priority != "" {
		story.Priority = req.Priority
	}
	if req.StoryPoints != nil {
		story.StoryPoints = req.StoryPoints
	}
	if req.SprintID != "" {
		story.SprintID = &req.SprintID
	}
	if req.AssigneeID != "" {
		story.AssigneeID = &req.AssigneeID
	}

	// Status logic
	if req.Status != "" {
		story.Status = req.Status
		if req.Status == "COMPLETED" || req.Status == "DONE" {
			now := time.Now()
			story.CompletedAt = &now
		} else if req.Status == "BACKLOG" || req.Status == "TODO" {
			story.CompletedAt = nil
		}
	}

	database.DB.Save(&story)

	// Notification Logic
	if req.AssigneeID != "" && (previousAssigneeID == nil || *previousAssigneeID != req.AssigneeID) {
		notification := models.Notification{
			ID:        generateCUIDStory(),
			UserID:    req.AssigneeID,
			Title:     "Historia de Usuario Asignada",
			Message:   "Se te ha asignado la historia \"" + story.Title + "\" en el proyecto " + story.Project.Name,
			Type:      "TASK_ASSIGNED",
			CreatedAt: time.Now(),
		}
		database.DB.Create(&notification)
	}

	c.JSON(http.StatusOK, gin.H{"data": story})
}

func DeleteUserStory(c *gin.Context) {
	id := c.Param("id")
	if result := database.DB.Delete(&models.UserStory{}, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar user story"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User story eliminado"})
}

func generateCUIDStory() string {
	return generateCUID()
}
