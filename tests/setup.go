package tests

import (
	"log"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() {
	var err error
	// Use in-memory SQLite for tests
	database.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	// Drop tables to ensure clean state
	database.DB.Migrator().DropTable(
		&models.User{},
		&models.Project{},
		&models.ProjectMember{},
		&models.Sprint{},
		&models.UserStory{},
		&models.Task{},
		&models.Rubric{},
		&models.Criteria{},
		&models.Evaluation{},
		&models.EvaluationCriteria{},
		&models.Chat{},
		&models.ChatParticipant{},
		&models.Message{},
		&models.Notification{},
		&models.RetrospectiveItem{},
		&models.Document{},
	)

	err = database.DB.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.ProjectMember{},
		&models.Sprint{},
		&models.UserStory{},
		&models.Task{},
		&models.Rubric{},
		&models.Criteria{},
		&models.Evaluation{},
		&models.EvaluationCriteria{},
		&models.Chat{},
		&models.ChatParticipant{},
		&models.Message{},
		&models.Notification{},
		&models.RetrospectiveItem{},
		&models.Document{},
	)
	if err != nil {
		log.Fatal("Failed to migrate test database:", err)
	}
}
