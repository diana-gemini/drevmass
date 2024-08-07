package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/controller/auth"
	"github.com/diana-gemini/drevmass/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	UserRepository models.UserRepository
}

// Get Profile godoc
// @Summary Get Profile
// @Security ApiKeyAuth
// @Tags user-profile-controller
// @ID get-user-profile
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile [get]
func (uc *UserController) GetProfile(c *gin.Context) {
	userID := c.GetUint("userID")

	user, _ := uc.UserRepository.GetUserByID(c, userID)
	if user.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find user")
		return
	}

	err := uc.UserRepository.CheckExpiredUserPoints(c, user.ID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to check user expired points")
		return
	}

	sum, err := uc.UserRepository.GetSumPointsByUserID(c, userID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to sum user points")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"mobile":   user.Mobile,
		"points":   sum,
	})
}

// Logout godoc
// @Summary Logout
// @Security ApiKeyAuth
// @Tags user-profile-controller
// @ID logout
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/logout [post]
func (uc *UserController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Logout": "Logout successful",
	})
}

// Change Password godoc
// @Summary Change Password
// @Security ApiKeyAuth
// @Tags user-profile-controller
// @ID change-password
// @Accept multipart/form-data
// @Produce json
// @Param current formData string true "current"
// @Param new formData string true "new"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/changepassword [post]
func (uc *UserController) ChangePassword(c *gin.Context) {
	var pass models.UserChangePassword
	if err := c.ShouldBind(&pass); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse request parameters")
		return
	}

	userID := c.GetUint("userID")

	user, _ := uc.UserRepository.GetUserByID(c, userID)
	if user.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find user")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass.Current)); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to verify current password")
		return
	}

	if !auth.CheckPassword(pass.New) {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to check new password")
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(pass.New), 10)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to hash password")
		return
	}

	updatePassword := string(hashPassword)

	err = uc.UserRepository.UpdateUserPassword(c, user, updatePassword)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to update user password")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"password": "Password successfully changed",
	})
}
