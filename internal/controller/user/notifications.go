package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
)

// Get User Notifications godoc
// @Summary Get User Notifications
// @Security ApiKeyAuth
// @Tags user-profile-notifications-controller
// @ID get-user-notifications
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/notifications [get]
func (uc *UserController) GetNotifications(c *gin.Context) {
	userID := c.GetUint("userID")

	user, _ := uc.UserRepository.GetUserByID(c, userID)
	if user.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find user")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": user.Notifications,
	})
}

// On User Notifications godoc
// @Summary On User Notifications
// @Security ApiKeyAuth
// @Tags user-profile-notifications-controller
// @ID on-user-notifications
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/notifications/on [post]
func (uc *UserController) OnNotifications(c *gin.Context) {
	userID := c.GetUint("userID")

	user, _ := uc.UserRepository.GetUserByID(c, userID)
	if user.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find user")
		return
	}

	err := uc.UserRepository.UpdateUserNotifications(c, user, "t")
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to update user notifications")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": "true",
	})
}

// Off User Notifications godoc
// @Summary Off User Notifications
// @Security ApiKeyAuth
// @Tags user-profile-notifications-controller
// @ID off-user-notifications
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/notifications/off [post]
func (uc *UserController) OffNotifications(c *gin.Context) {
	userID := c.GetUint("userID")

	user, _ := uc.UserRepository.GetUserByID(c, userID)
	if user.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find user")
		return
	}

	err := uc.UserRepository.UpdateUserNotifications(c, user, "f")
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to update user notifications")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": "false",
	})
}
