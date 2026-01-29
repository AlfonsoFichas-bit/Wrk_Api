package handlers

import (
	"net/http"
	"time"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"

	"github.com/gin-gonic/gin"
)

type CreateSprintRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	ProjectID   string  `json:"projectId" binding:"required"`
	StartDate   *string `json:"startDate"`
	EndDate     *string `json:"endDate"`
	Status      string  `json:"status"`
}

type UpdateSprintRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	StartDate   *string `json:"startDate"`
	EndDate     *string `json:"endDate"`
	Status      string  `json:"status"`
}

type AddStoryRequest struct {
	UserStoryID string `json:"userStoryId" binding:"required"`
}

func GetAllSprints(c *gin.Context) {
	var sprints []models.Sprint
	if result := database.DB.Preload("Project").Preload("Tasks").Preload("UserStories").Preload("Evaluations").Find(&sprints); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener sprints", "details": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sprints})
}

func GetSprint(c *gin.Context) {
	id := c.Param("id")
	var sprint models.Sprint
	if result := database.DB.Preload("Project").Preload("Tasks").Preload("UserStories").Preload("Evaluations").First(&sprint, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sprint no encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sprint})
}

func CreateSprint(c *gin.Context) {
	var req CreateSprintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var startDate, endDate time.Time
	if req.StartDate != nil {
		startDate, _ = time.Parse(time.RFC3339, *req.StartDate)
	}
	if req.EndDate != nil {
		endDate, _ = time.Parse(time.RFC3339, *req.EndDate)
	}

	status := "PLANNING"
	if req.Status != "" {
		status = req.Status
	}

	sprint := models.Sprint{
		ID:          generateCUIDSprint(),
		Name:        req.Name,
		Description: req.Description,
		ProjectID:   req.ProjectID,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      status,
	}

	if result := database.DB.Create(&sprint); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear sprint", "details": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": sprint})
}

func UpdateSprint(c *gin.Context) {
	id := c.Param("id")
	var req UpdateSprintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sprint models.Sprint
	if result := database.DB.First(&sprint, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sprint no encontrado"})
		return
	}

	if req.Name != "" {
		sprint.Name = req.Name
	}
	if req.Description != nil {
		sprint.Description = req.Description
	}
	if req.Status != "" {
		sprint.Status = req.Status
	}
	if req.StartDate != nil {
		t, _ := time.Parse(time.RFC3339, *req.StartDate)
		sprint.StartDate = t
	}
	if req.EndDate != nil {
		t, _ := time.Parse(time.RFC3339, *req.EndDate)
		sprint.EndDate = t
	}

	database.DB.Save(&sprint)
	c.JSON(http.StatusOK, gin.H{"data": sprint})
}

func AddStoryToSprint(c *gin.Context) {
	sprintID := c.Param("id")
	var req AddStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userStory models.UserStory
	if result := database.DB.First(&userStory, "id = ?", req.UserStoryID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Historia no encontrada"})
		return
	}

	if userStory.SprintID != nil && *userStory.SprintID == sprintID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La historia ya está en el sprint"})
		return
	}

	userStory.SprintID = &sprintID
	database.DB.Save(&userStory)

	c.JSON(http.StatusCreated, gin.H{"data": userStory})
}

func DeleteSprint(c *gin.Context) {
	id := c.Param("id")
	if result := database.DB.Delete(&models.Sprint{}, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar sprint"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Sprint eliminado"}})
}

// Reuse or duplicate generator for independence
func generateCUIDSprint() string {
	return generateCUID() // Assuming generateCUID is accessible or duplicated. Ideally move to utils.
}
