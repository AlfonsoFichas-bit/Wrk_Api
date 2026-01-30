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

type CreateCriteriaRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	MaxScore    int    `json:"maxScore"`
	Weight      int    `json:"weight"`
}

type CreateRubricRequest struct {
	ProjectID   string                  `json:"projectId"`
	Name        string                  `json:"name" binding:"required"`
	Description string                  `json:"description"`
	Criteria    []CreateCriteriaRequest `json:"criteria"`
}

func GetAllRubrics(c *gin.Context) {
	projectID := c.Query("projectId")
	var rubrics []models.Rubric
	
	query := database.DB.Preload("Criteria")
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	if result := query.Find(&rubrics); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener rúbricas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rubrics})
}

func GetRubric(c *gin.Context) {
	id := c.Param("id")
	var rubric models.Rubric
	if result := database.DB.Preload("Criteria").First(&rubric, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rúbrica no encontrada"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rubric})
}

func CreateRubric(c *gin.Context) {
	var req CreateRubricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		rubric := models.Rubric{
			ID:          utils.GenerateCUID(),
			Name:        req.Name,
			Description: &req.Description,
			CreatedAt:   time.Now(),
		}
		
		if req.ProjectID != "" {
			rubric.ProjectID = &req.ProjectID
		}

		if err := tx.Create(&rubric).Error; err != nil {
			return err
		}

		for _, crit := range req.Criteria {
			maxScore := 100
			if crit.MaxScore > 0 {
				maxScore = crit.MaxScore
			}
			weight := 1
			if crit.Weight > 0 {
				weight = crit.Weight
			}

			criteria := models.Criteria{
				ID:          utils.GenerateCUID(),
				RubricID:    rubric.ID,
				Name:        crit.Name,
				Description: &crit.Description,
				MaxScore:    maxScore,
				Weight:      weight,
			}
			if err := tx.Create(&criteria).Error; err != nil {
				return err
			}
		}
		
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear rúbrica"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Rúbrica creada exitosamente"})
}

func DeleteRubric(c *gin.Context) {
	id := c.Param("id")
	if result := database.DB.Delete(&models.Rubric{}, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar rúbrica"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Rúbrica eliminada"})
}
