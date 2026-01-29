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
	}
}
