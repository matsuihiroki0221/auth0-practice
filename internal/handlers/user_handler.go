package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserHandler handles GET requests to fetch user information
func GetUsersHandler(c *gin.Context) {
	// Implement your logic to fetch user information here
	c.JSON(http.StatusOK, gin.H{
		"message": "GET request received",
	})
}

// CreateUserHandler handles POST requests to create a new user
func CreateUserHandler(c *gin.Context) {
	// Implement your logic to create a new user here
	c.JSON(http.StatusCreated, gin.H{
		"message": "POST request received",
	})
}
