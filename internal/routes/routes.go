package routes

import (
	"Wrk_Api/internal/handlers"
	"Wrk_Api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Global Middleware
	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")

	// Auth Routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Protected Routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Users
		protected.GET("/users", handlers.GetAllUsers)
		protected.GET("/users/:id", handlers.GetUser)
		protected.POST("/users/", handlers.CreateUser)
		protected.PUT("/users/:id", handlers.UpdateUser)
		protected.DELETE("/users/:id", handlers.DeleteUser)

		// Projects
		projects := protected.Group("/projects")
		{
			projects.GET("/", handlers.GetAllProjects)
			projects.GET("/:id", handlers.GetProject)
			projects.POST("/", handlers.CreateProject)
			projects.PUT("/:id", handlers.UpdateProject)
			projects.DELETE("/:id", handlers.DeleteProject)

			// Project Members
			projects.POST("/:id/members", handlers.AddProjectMember)
			projects.DELETE("/:id/members/:userId", handlers.RemoveProjectMember)
		}

		// Sprints
		sprints := protected.Group("/sprints")
		{
			sprints.GET("/", handlers.GetAllSprints)
			sprints.GET("/:id", handlers.GetSprint)
			sprints.POST("/", handlers.CreateSprint)
			sprints.PUT("/:id", handlers.UpdateSprint)
			sprints.DELETE("/:id", handlers.DeleteSprint)
			
			// Sprint Actions
			sprints.POST("/:id/add-story", handlers.AddStoryToSprint)
		}

		// User Stories
		userStories := protected.Group("/user-stories")
		{
			userStories.GET("/", handlers.GetAllUserStories)
			userStories.GET("/:id", handlers.GetUserStory)
			userStories.POST("/", handlers.CreateUserStory)
			userStories.PUT("/:id", handlers.UpdateUserStory)
			userStories.DELETE("/:id", handlers.DeleteUserStory)
		}

		// Tasks
		tasks := protected.Group("/tasks")
		{
			tasks.GET("/", handlers.GetAllTasks)
			tasks.GET("/:id", handlers.GetTask)
			tasks.POST("/", handlers.CreateTask)
			tasks.PUT("/:id", handlers.UpdateTask)
			tasks.DELETE("/:id", handlers.DeleteTask)

			// Task Actions
			tasks.POST("/:id/evaluate", handlers.EvaluateTask)
		}

		// Chat
		chat := protected.Group("/chat")
		{
			// Project Chat
			chat.GET("/:projectId/messages", handlers.GetProjectMessages)
			chat.POST("/:projectId/messages", handlers.SendProjectMessage)

			// Direct Chat
			chat.GET("/user/:userId/all", handlers.GetDirectChats)
			chat.POST("/direct", handlers.CreateOrGetDirectChat)
			chat.GET("/conversation/:chatId/messages", handlers.GetConversationMessages)
			chat.POST("/conversation/:chatId/messages", handlers.SendConversationMessage)
		}

		// Notifications
		notifications := protected.Group("/notifications")
		{
			notifications.GET("/", handlers.GetNotifications)
			notifications.PUT("/:id/read", handlers.MarkNotificationRead)
		}

		// Rubrics
		rubrics := protected.Group("/rubrics")
		{
			rubrics.GET("/", handlers.GetAllRubrics)
			rubrics.GET("/:id", handlers.GetRubric)
			rubrics.POST("/", handlers.CreateRubric)
			rubrics.DELETE("/:id", handlers.DeleteRubric)
		}

		// Evaluations (Module)
		evaluations := protected.Group("/evaluations")
		{
			evaluations.GET("/:id", handlers.GetEvaluation)
			evaluations.POST("/", handlers.CreateEvaluation)
			evaluations.PUT("/:id", handlers.UpdateEvaluation)
			
			evaluations.GET("/task/:taskId", handlers.GetTaskEvaluations)
			evaluations.GET("/sprint/:sprintId", handlers.GetSprintEvaluations)
			evaluations.GET("/project/:projectId/general", handlers.GetProjectEvaluations)
			evaluations.GET("/student/:studentId", handlers.GetStudentEvaluations)
		}

		// Retrospectives
		retrospectives := protected.Group("/retrospectives")
		{
			retrospectives.GET("/:sprintId", handlers.GetSprintRetrospective)
			retrospectives.POST("/", handlers.CreateRetrospectiveItem)
			retrospectives.DELETE("/:id", handlers.DeleteRetrospectiveItem)
		}

		// Documents
		documents := protected.Group("/documents")
		{
			documents.GET("/:projectId", handlers.GetProjectDocuments)
			documents.POST("/", handlers.UploadDocument)
			documents.DELETE("/:id", handlers.DeleteDocument)
		}

		// Metrics
		metrics := protected.Group("/metrics")
		{
			metrics.GET("/sprints/:sprintId/burndown", handlers.GetSprintBurndown)
			metrics.GET("/projects/:projectId/velocity", handlers.GetProjectVelocity)
			metrics.GET("/projects/:projectId/contribution", handlers.GetProjectContribution)
			metrics.GET("/export/projects/:projectId", handlers.ExportProjectCSV)
		}
	}
}
