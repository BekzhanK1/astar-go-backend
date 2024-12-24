package handlers

import (
	"api-gateway/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserServiceClient
}

func NewUserHandler(userService *services.UserServiceClient) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) UserProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	profile, err := h.userService.GetProfile(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get profile: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"first_name": profile.User.FirstName,
		"last_name":  profile.User.LastName,
		"email":      profile.User.Email,
	})
}
