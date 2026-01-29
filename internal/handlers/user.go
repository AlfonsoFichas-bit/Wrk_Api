package handlers

import (
	"net/http"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"

	"github.com/gin-gonic/gin"
)

// Helper to check if user is admin
func isAdmin(c *gin.Context) bool {
	role, exists := c.Get("role")
	if !exists {
		return false
	}
	return role == "ADMIN" // Adjust based on your role definitions
}

func GetAllUsers(c *gin.Context) {
	// Optional: Restriction to admins or project members?
	// The original API seems to allow getting users, possibly for assignment.
	
	var users []models.User
	if result := database.DB.Find(&users); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
		return
	}

	// Filter sensitive data
	var response []gin.H
	for _, u := range users {
		response = append(response, gin.H{
			"id":    u.ID,
			"name":  u.Name,
			"email": u.Email,
			"role":  u.Role,
			"avatar": u.Avatar,
		})
	}

	c.JSON(http.StatusOK, response)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if result := database.DB.First(&user, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
		"avatar": user.Avatar,
	})
}
