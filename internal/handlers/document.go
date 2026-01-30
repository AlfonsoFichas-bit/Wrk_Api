package handlers

import (
	"fmt"
	"net/http"
	"time"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"
	"Wrk_Api/internal/utils"

	"github.com/gin-gonic/gin"
)

type UploadDocumentRequest struct {
	ProjectID string `form:"projectId" binding:"required"`
	Name      string `form:"name" binding:"required"`
}

func GetProjectDocuments(c *gin.Context) {
	projectID := c.Param("projectId")
	var docs []models.Document
	// Get latest versions (where parent_id is null usually, or handle versioning logic)
	// For simplicity, returning all
	if result := database.DB.Where("project_id = ?", projectID).Find(&docs); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener documentos"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": docs})
}

func UploadDocument(c *gin.Context) {
	// Multipart form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	projectID := c.PostForm("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID required"})
		return
	}
	
	// Mock upload (save to /tmp or just DB record)
	// In real app: Save to S3 or disk
	url := fmt.Sprintf("/uploads/%s_%s", utils.GenerateCUID(), file.Filename)
	size := int(file.Size / 1024)

	doc := models.Document{
		ID:         utils.GenerateCUID(),
		ProjectID:  projectID,
		Name:       file.Filename,
		URL:        url,
		Type:       "FILE", // Detect MIME type in real app
		Size:       &size,
		Version:    1,
		UploadedAt: time.Now(),
	}

	if result := database.DB.Create(&doc).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving document metadata"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": doc})
}

func DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	if result := database.DB.Delete(&models.Document{}, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting document"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted"})
}
