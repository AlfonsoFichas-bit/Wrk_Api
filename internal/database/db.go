package database

import (
	"log"
	"os"

	"Wrk_Api/internal/models"
	"Wrk_Api/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "test.db"
	}

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to database")

	// Auto Migrate
	err = DB.AutoMigrate(
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
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")

	// Seed Admin
	seedAdmin()
}

func seedAdmin() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
		admin := models.User{
			ID:       utils.GenerateCUID(),
			Name:     "Admin System",
			Email:    "admin@workflo.ws",
			Password: string(hashedPassword),
			Role:     "ADMIN",
			Active:   true,
		}
		if err := DB.Create(&admin).Error; err != nil {
			log.Println("Error seeding admin:", err)
		} else {
			log.Println("Default admin created: admin@workflo.ws / admin123")
		}
	}
}
