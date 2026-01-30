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

type CreateTaskRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	ProjectID   string  `json:"projectId" binding:"required"`
	AssigneeID  string  `json:"assigneeId"`
	Priority    string  `json:"priority"`
	Deadline    *string `json:"deadline"`
	Status      string  `json:"status"`
	SprintID    string  `json:"sprintId"`
	UserStoryID string  `json:"userStoryId"`
}

type UpdateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	AssigneeID  string  `json:"assigneeId"`
	Priority    string  `json:"priority"`
	Deadline    *string `json:"deadline"`
	Status      string  `json:"status"`
	SprintID    string  `json:"sprintId"`
	UserStoryID string  `json:"userStoryId"`
}

type CriteriaScore struct {
	CriteriaID string `json:"criteriaId"`
	Score      int    `json:"score"`
}

type EvaluateTaskRequest struct {
	Score          int             `json:"score" binding:"required"`
	Feedback       string          `json:"feedback"`
	EvaluatorID    string          `json:"evaluatorId" binding:"required"`
	CriteriaScores []CriteriaScore `json:"criteriaScores"`
}

func GetAllTasks(c *gin.Context) {
	assigneeID := c.Query("assigneeId")
	projectID := c.Query("projectId")

	var tasks []models.Task
	query := database.DB.Preload("Assignee").Preload("Project").Preload("Evaluations")

	if assigneeID != "" {
		query = query.Where("assignee_id = ?", assigneeID)
	}
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	if result := query.Find(&tasks); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener tareas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func GetTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if result := database.DB.Preload("Assignee").Preload("Project").Preload("Evaluations").First(&task, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarea no encontrada"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	priority := "MEDIUM"
	if req.Priority != "" {
		priority = req.Priority
	}
	status := "TODO"
	if req.Status != "" {
		status = req.Status
	}

	var deadline *time.Time
	if req.Deadline != nil {
		t, _ := time.Parse(time.RFC3339, *req.Deadline)
		deadline = &t
	}

	task := models.Task{
		ID:          utils.GenerateCUID(),
		Title:       req.Title,
		Description: &req.Description,
		ProjectID:   req.ProjectID,
		AssigneeID:  nil,
		Priority:    priority,
		Deadline:    deadline,
		Status:      status,
	}

	if req.AssigneeID != "" {
		task.AssigneeID = &req.AssigneeID
	}
	if req.SprintID != "" {
		task.SprintID = &req.SprintID
	}
	if req.UserStoryID != "" {
		task.UserStoryID = &req.UserStoryID
	}

	if result := database.DB.Create(&task); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear tarea"})
		return
	}

	// Notificación
	if req.AssigneeID != "" {
		notification := models.Notification{
			ID:        utils.GenerateCUID(),
			UserID:    req.AssigneeID,
			Title:     "Nueva Tarea Asignada",
			Message:   "Se te ha asignado la tarea: " + task.Title,
			Type:      "TASK_ASSIGNED",
			CreatedAt: time.Now(),
		}
		database.DB.Create(&notification)
	}

	c.JSON(http.StatusCreated, gin.H{"data": task})
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task models.Task
	if result := database.DB.First(&task, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarea no encontrada"})
		return
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = &req.Description
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.AssigneeID != "" {
		task.AssigneeID = &req.AssigneeID
	}
	if req.SprintID != "" {
		task.SprintID = &req.SprintID
	}
	if req.UserStoryID != "" {
		task.UserStoryID = &req.UserStoryID
	}
	if req.Deadline != nil {
		t, _ := time.Parse(time.RFC3339, *req.Deadline)
		task.Deadline = &t
	}

	if req.Status != "" {
		task.Status = req.Status
		if req.Status == "COMPLETED" || req.Status == "DONE" {
			now := time.Now()
			task.CompletedAt = &now
		} else if req.Status == "TODO" || req.Status == "IN_PROGRESS" || req.Status == "PENDING" {
			task.CompletedAt = nil
		}
	}

	database.DB.Save(&task)
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if result := database.DB.Delete(&models.Task{}, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar tarea"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Tarea eliminada"}})
}

func EvaluateTask(c *gin.Context) {
	taskID := c.Param("id")
	var req EvaluateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task models.Task
	if result := database.DB.First(&task, "id = ?", taskID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarea no encontrada"})
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		evaluation := models.Evaluation{
			ID:          utils.GenerateCUID(),
			TaskID:      &taskID,
			ProjectID:   task.ProjectID,
			EvaluatorID: req.EvaluatorID,
			Score:       &req.Score,
			Feedback:    &req.Feedback,
			Status:      "COMPLETED",
			CreatedAt:   time.Now(),
		}

		if err := tx.Create(&evaluation).Error; err != nil {
			return err
		}

		if len(req.CriteriaScores) > 0 {
			for _, cs := range req.CriteriaScores {
				ec := models.EvaluationCriteria{
					ID:           utils.GenerateCUID(),
					EvaluationID: evaluation.ID,
					CriteriaID:   cs.CriteriaID,
					Score:        cs.Score,
				}
				if err := tx.Create(&ec).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la evaluación"})
		return
	}

	// Notify Assignee
	if task.AssigneeID != nil {
		notification := models.Notification{
			ID:        utils.GenerateCUID(),
			UserID:    *task.AssigneeID,
			Title:     "Tarea Evaluada",
			Message:   "Tu tarea \"" + task.Title + "\" ha sido evaluada",
			Type:      "EVALUATION_COMPLETED",
			CreatedAt: time.Now(),
		}
		database.DB.Create(&notification)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Evaluación guardada"})
}
