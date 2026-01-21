package handlers

import (
	"net/http"

	"go-crud-app/internal/database"
	"go-crud-app/internal/middleware"
	"go-crud-app/internal/models"

	"github.com/gin-gonic/gin"
)

// UpdateUserRequest represents the user update request payload
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// GetCurrentUser returns the currently authenticated user
func GetCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// GetAllUsers returns all registered users except the current user
func GetAllUsers(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var users []models.User
	// Exclude the current user from the list
	if err := database.DB.Where("id != ?", userID).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
		})
		return
	}

	// Convert to response format
	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"users": userResponses,
		"count": len(userResponses),
	})
}

// GetUserByID returns a specific user by ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateUser updates the current user's information
func UpdateUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Get the ID from URL parameter
	id := c.Param("id")

	// Parse ID and check if user is updating their own profile
	var targetUserID uint
	if _, err := c.Params.Get("id"); err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Convert string ID to uint
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}
	targetUserID = user.ID

	// Users can only update their own profile
	if targetUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only update your own profile",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	// Update fields if provided
	updates := make(map[string]interface{})
	if req.Username != "" {
		if !usernameRegex.MatchString(req.Username) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Username must be 3-50 characters and contain only letters, numbers, and underscores",
			})
			return
		}
		updates["username"] = req.Username
	}
	if req.Email != "" {
		if !emailRegex.MatchString(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid email format",
			})
			return
		}
		updates["email"] = req.Email
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No fields to update",
		})
		return
	}

	// Update user
	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	// Fetch updated user
	database.DB.First(&user, userID)

	c.JSON(http.StatusOK, user.ToResponse())
}

// DeleteUser deletes the current user's account
func DeleteUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Get the ID from URL parameter
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Users can only delete their own profile
	if user.ID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only delete your own profile",
		})
		return
	}

	// Soft delete user
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
