package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"
	"Wrk_Api/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateProjectRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	OwnerID     string  `json:"ownerId" binding:"required"`
	StartDate   *string `json:"startDate"`
	EndDate     *string `json:"endDate"`
}

type UpdateProjectRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
	StartDate   *string `json:"startDate"`
	EndDate     *string `json:"endDate"`
}

type AddMemberRequest struct {
	UserID string `json:"userId" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

func GetAllProjects(c *gin.Context) {
	memberID := c.Query("memberId")
	var projects []models.Project

	query := database.DB.Preload("Owner").Preload("Members").Preload("Sprints")

	if memberID != "" {
		// OR: [ { ownerId: memberId }, { members: { some: { userId: memberId } } } ]
		// GORM equivalent:
		query = query.Joins("LEFT JOIN project_members ON project_members.project_id = projects.id").
			Where("projects.owner_id = ? OR project_members.user_id = ?", memberID, memberID).
			Group("projects.id")
	}

	if result := query.Find(&projects); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener proyectos", "details": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func GetProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	// Include owner, members (with user), sprints, userStories, tasks
	if result := database.DB.
		Preload("Owner").
		Preload("Members.User").
		Preload("Sprints").
		Preload("UserStories").
		Preload("Tasks").
		First(&project, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proyecto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": project})
}

func CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var startDate, endDate *time.Time
	if req.StartDate != nil {
		t, _ := time.Parse(time.RFC3339, *req.StartDate)
		startDate = &t
	}
	if req.EndDate != nil {
		t, _ := time.Parse(time.RFC3339, *req.EndDate)
		endDate = &t
	}

	project := models.Project{
		ID:          utils.GenerateCUID(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     req.OwnerID,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      "ACTIVE",
	}

	if result := database.DB.Create(&project); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear proyecto", "details": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": project})
}

func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project models.Project
	if result := database.DB.First(&project, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proyecto no encontrado"})
		return
	}

	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != nil {
		project.Description = req.Description
	}
	if req.Status != "" {
		project.Status = req.Status
	}
	if req.StartDate != nil {
		t, _ := time.Parse(time.RFC3339, *req.StartDate)
		project.StartDate = &t
	}
	if req.EndDate != nil {
		t, _ := time.Parse(time.RFC3339, *req.EndDate)
		project.EndDate = &t
	}

	database.DB.Save(&project)
	c.JSON(http.StatusOK, gin.H{"data": project})
}

func DeleteProject(c *gin.Context) {
	id := c.Param("id")
	if result := database.DB.Delete(&models.Project{}, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar proyecto"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Proyecto eliminado"}})
}

func AddProjectMember(c *gin.Context) {
	projectID := c.Param("id")
	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check existing member
	var existingMember models.ProjectMember
	if result := database.DB.Where("project_id = ? AND user_id = ?", projectID, req.UserID).First(&existingMember); result.Error == nil {
		// Update role
		existingMember.Role = req.Role
		database.DB.Save(&existingMember)
		c.JSON(http.StatusOK, gin.H{"data": existingMember, "message": "Rol actualizado"})
		return
	}

	member := models.ProjectMember{
		ID:        utils.GenerateCUID(),
		ProjectID: projectID,
		UserID:    req.UserID,
		Role:      req.Role,
	}

	if result := database.DB.Create(&member); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al asignar miembro", "details": result.Error.Error()})
		return
	}

	// Notify
	var project models.Project
	database.DB.First(&project, "id = ?", projectID)
	
	notification := models.Notification{
		ID:        utils.GenerateCUID(),
		UserID:    req.UserID,
		Title:     "Nuevo Proyecto Asignado",
		Message:   "Has sido añadido al proyecto \"" + project.Name + "\" como " + req.Role,
		Type:      "PROJECT_ASSIGNED",
		CreatedAt: time.Now(),
	}
	database.DB.Create(&notification)

	c.JSON(http.StatusCreated, gin.H{"data": member})
}

func RemoveProjectMember(c *gin.Context) {
	projectID := c.Param("id")
	userID := c.Param("userId")

	if result := database.DB.Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&models.ProjectMember{}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar miembro"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Miembro eliminado del proyecto"})
}

func generateCUIDProject() string {
	b := make([]byte, 12)
	rand.Read(b)
	return fmt.Sprintf("c%s", hex.EncodeToString(b))
}
