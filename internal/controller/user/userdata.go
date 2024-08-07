package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
)

// Get User Data godoc
// @Summary Get User Data
// @Security ApiKeyAuth
// @Tags user-profile-userdata-controller
// @ID get-user-data
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/userdata [get]
func (uc *UserController) GetUserData(c *gin.Context) {
	userID := c.GetUint("userID")

	user, _ := uc.UserRepository.GetUserByID(c, userID)
	if user.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find user")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":  user.Username,
		"mobile":    user.Mobile,
		"email":     user.Email,
		"birthdate": user.BirthDate,
		"gender":    user.Gender,
		"height":    user.Height,
		"weght":     user.Weight,
		"activity":  user.Activity,
	})
}

// Update User Data godoc
// @Summary Update User Data
// @Security ApiKeyAuth
// @Tags user-profile-userdata-controller
// @ID update-user-data
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "username"
// @Param mobile formData string true "mobile"
// @Param birthdate formData string true "birthdate"
// @Param gender formData string true "gender"
// @Param height formData string true "height"
// @Param weight formData string true "weight"
// @Param activity formData string true "activity"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/userdata [put]
func (uc *UserController) UpdateUserData(c *gin.Context) {
	var updateUser models.User
	if err := c.ShouldBind(&updateUser); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse request parameters")
		return
	}

	userID := c.GetUint("userID")

	user, _ := uc.UserRepository.GetUserByID(c, userID)
	if user.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find user")
		return
	}

	err := uc.UserRepository.UpdateUser(c, user, updateUser)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to update user")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":  updateUser.Username,
		"mobile":    updateUser.Mobile,
		"birthdate": updateUser.BirthDate,
		"gender":    updateUser.Gender,
		"height":    updateUser.Height,
		"weight":    updateUser.Weight,
		"activity":  updateUser.Activity,
	})
}

// Delete Profile godoc
// @Summary Delete Profile
// @Security ApiKeyAuth
// @Tags user-profile-userdata-controller
// @ID delete-profile
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/userdata [delete]
func (uc *UserController) DeleteProfile(c *gin.Context) {
	userID := c.GetUint("userID")

	user, _ := uc.UserRepository.GetUserByID(c, userID)
	if user.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find user")
		return
	}

	err := uc.UserRepository.DeleteUser(c, user)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to update user")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": "Profile deleted",
	})
}
