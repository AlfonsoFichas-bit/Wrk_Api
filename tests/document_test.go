package tests

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"
	"Wrk_Api/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDocumentHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	// User & Project
	user := models.User{ID: "u1", Name: "User", Email: "test@docs.com"}
	database.DB.Create(&user)
	project := models.Project{ID: "p1", Name: "Doc Project", OwnerID: user.ID}
	database.DB.Create(&project)

	token := generateTestToken(user.ID, user.Email, user.Role)
	authHeader := "Bearer " + token

	t.Run("UploadDocument", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		
		// Add file field
		part, _ := writer.CreateFormFile("file", "test.txt")
		part.Write([]byte("file content"))
		
		// Add projectId field
		writer.WriteField("projectId", project.ID)
		writer.Close()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/documents/", body)
		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var doc models.Document
		database.DB.Where("project_id = ?", project.ID).First(&doc)
		assert.Equal(t, "test.txt", doc.Name)
	})
}
