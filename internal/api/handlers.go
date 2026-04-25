package api

import (
	"net/http"

	"github.com/Carlos-Marrugo/pigbank-user-service/internal/models"
	"github.com/Carlos-Marrugo/pigbank-user-service/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	res, err := service.RegisterHandler(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": res})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials format"})
		return
	}

	token, err := service.LoginHandler(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.Param("user_id")

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
		return
	}

	err := service.UpdateUserProfile(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully for " + userID})
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := c.Param("user_id")
	authenticatedUUID := c.GetString("user_uuid")

	if userID != authenticatedUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	var req models.AvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	url, err := service.UploadAvatar(c.Request.Context(), userID, req.AvatarBase64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatar_url": url})
}
