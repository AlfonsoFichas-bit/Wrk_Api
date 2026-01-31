package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"
	"Wrk_Api/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestEvaluationHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	// User & Project & Criteria
	user := models.User{ID: "eval-u1", Name: "Evaluator", Email: "eval@test.com", Role: "SCRUM_MASTER"}
	student := models.User{ID: "stud-1", Name: "Student", Email: "student@test.com"}
	database.DB.Create(&user)
	database.DB.Create(&student)
	
	project := models.Project{ID: "p1", Name: "Eval Project", OwnerID: user.ID}
	database.DB.Create(&project)
	
	// Create a rubric criterion
	rubric := models.Rubric{ID: "rub1", ProjectID: &project.ID, Name: "Rubric"}
	database.DB.Create(&rubric)
	crit := models.Criteria{ID: "crit1", RubricID: rubric.ID, Name: "Code Quality", MaxScore: 10}
	database.DB.Create(&crit)

	token := generateTestToken(user.ID, user.Email, user.Role)
	authHeader := "Bearer " + token

	var evalID string

	t.Run("CreateGenericEvaluation", func(t *testing.T) {
		body := map[string]interface{}{
			"projectId":   project.ID,
			"evaluatorId": user.ID,
			"feedback":    "Great work",
			"score":       10,
			"criteriaScores": []map[string]interface{}{
				{"criteriaId": crit.ID, "score": 10},
			},
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/evaluations/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		
		// Find ID
		var eval models.Evaluation
		database.DB.First(&eval, "feedback = ?", "Great work")
		evalID = eval.ID
		assert.NotEmpty(t, evalID)
	})

	t.Run("GetEvaluation", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/evaluations/"+evalID, nil)
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Great work", resp["data"]["Feedback"])
	})

	t.Run("UpdateEvaluation", func(t *testing.T) {
		body := map[string]interface{}{
			"feedback": "Updated feedback",
			"score":    9,
			"criteriaScores": []map[string]interface{}{
				{"criteriaId": crit.ID, "score": 9},
			},
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/evaluations/"+evalID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", authHeader)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var eval models.Evaluation
		database.DB.First(&eval, "id = ?", evalID)
		assert.Equal(t, "Updated feedback", *eval.Feedback)
		assert.Equal(t, 9, *eval.Score)
	})
}
