package handlers

import (
	"net/http"
	"sort"
	"time"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"
	"Wrk_Api/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateGenericEvaluationRequest struct {
	ProjectID      string          `json:"projectId" binding:"required"`
	TaskID         *string         `json:"taskId"`
	SprintID       *string         `json:"sprintId"`
	EvaluatorID    string          `json:"evaluatorId" binding:"required"`
	Feedback       string          `json:"feedback"`
	Score          int             `json:"score"`
	CriteriaScores []CriteriaScore `json:"criteriaScores" binding:"required"`
}

type UpdateEvaluationRequest struct {
	Feedback       string          `json:"feedback"`
	Score          int             `json:"score"`
	CriteriaScores []CriteriaScore `json:"criteriaScores" binding:"required"`
}

func GetEvaluation(c *gin.Context) {
	id := c.Param("id")
	var eval models.Evaluation
	if result := database.DB.Preload("Criteria.Criteria").Preload("Evaluator").First(&eval, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Evaluación no encontrada"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": eval})
}

func GetTaskEvaluations(c *gin.Context) {
	taskID := c.Param("taskId")
	var evals []models.Evaluation
	database.DB.Preload("Evaluator").Preload("Criteria").Where("task_id = ?", taskID).Order("created_at desc").Find(&evals)
	c.JSON(http.StatusOK, gin.H{"data": evals})
}

func GetSprintEvaluations(c *gin.Context) {
	sprintID := c.Param("sprintId")
	var evals []models.Evaluation
	database.DB.Preload("Evaluator").Preload("Criteria").Where("sprint_id = ?", sprintID).Order("created_at desc").Find(&evals)
	c.JSON(http.StatusOK, gin.H{"data": evals})
}

func GetProjectEvaluations(c *gin.Context) {
	projectID := c.Param("projectId")
	var evals []models.Evaluation
	// General evaluations (no task, no sprint)
	database.DB.Preload("Evaluator").Preload("Criteria").
		Where("project_id = ? AND task_id IS NULL AND sprint_id IS NULL", projectID).
		Order("created_at desc").Find(&evals)
	c.JSON(http.StatusOK, gin.H{"data": evals})
}

func GetStudentEvaluations(c *gin.Context) {
	studentID := c.Param("studentId")

	// 1. Task Evaluations (assigned to student)
	var taskEvals []models.Evaluation
	database.DB.Joins("JOIN tasks ON tasks.id = evaluations.task_id").
		Where("tasks.assignee_id = ?", studentID).
		Preload("Project").Preload("Task").Preload("Sprint").Preload("Evaluator").
		Find(&taskEvals)

	// 2. Team Evaluations (Project/Sprint level where student is member)
	var projectIDs []string
	database.DB.Model(&models.ProjectMember{}).Where("user_id = ?", studentID).Pluck("project_id", &projectIDs)

	var teamEvals []models.Evaluation
	if len(projectIDs) > 0 {
		database.DB.Where("project_id IN ? AND task_id IS NULL", projectIDs).
			Preload("Project").Preload("Sprint").Preload("Evaluator").
			Find(&teamEvals)
	}

	// Merge and Sort
	allEvals := append(taskEvals, teamEvals...)
	sort.Slice(allEvals, func(i, j int) bool {
		return allEvals[i].CreatedAt.After(allEvals[j].CreatedAt)
	})

	c.JSON(http.StatusOK, gin.H{"data": allEvals})
}

func CreateEvaluation(c *gin.Context) {
	var req CreateGenericEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		eval := models.Evaluation{
			ID:          utils.GenerateCUID(),
			ProjectID:   req.ProjectID,
			TaskID:      req.TaskID,
			SprintID:    req.SprintID,
			EvaluatorID: req.EvaluatorID,
			Feedback:    &req.Feedback,
			Score:       &req.Score,
			Status:      "COMPLETED",
			CreatedAt:   time.Now(),
		}

		if err := tx.Create(&eval).Error; err != nil {
			return err
		}

		for _, cs := range req.CriteriaScores {
			ec := models.EvaluationCriteria{
				ID:           utils.GenerateCUID(),
				EvaluationID: eval.ID,
				CriteriaID:   cs.CriteriaID,
				Score:        cs.Score,
			}
			if err := tx.Create(&ec).Error; err != nil {
				return err
			}
		}
		
		// Return with ID
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar evaluación"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Evaluación creada"})
}

func UpdateEvaluation(c *gin.Context) {
	id := c.Param("id")
	var req UpdateEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Update Evaluation
		if err := tx.Model(&models.Evaluation{}).Where("id = ?", id).Updates(map[string]interface{}{
			"feedback": req.Feedback,
			"score":    req.Score,
		}).Error; err != nil {
			return err
		}

		// Replace Criteria
		if err := tx.Delete(&models.EvaluationCriteria{}, "evaluation_id = ?", id).Error; err != nil {
			return err
		}

		for _, cs := range req.CriteriaScores {
			ec := models.EvaluationCriteria{
				ID:           utils.GenerateCUID(),
				EvaluationID: id,
				CriteriaID:   cs.CriteriaID,
				Score:        cs.Score,
			}
			if err := tx.Create(&ec).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar evaluación"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Evaluación actualizada"})
}
