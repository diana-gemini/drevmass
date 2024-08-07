package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// PasswordReset godoc
// @Summary Password Reset
// @Tags auth-controller
// @ID password-reset
// @Accept multipart/form-data
// @Produce json
// @Param token query string true "token"
// @Param password formData string true "password"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object}  models.ErrorResponse
// @Router /passwordreset [post]
func (ac *AuthController) PasswordReset(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find token")
		return
	}

	email, ok := resetTokens[token]
	if !ok {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to give email with token")
		return
	}

	var password models.UserPassword
	if err := c.ShouldBind(&password); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse request parameters")
		return
	}

	if !CheckPassword(password.Password) {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to check password")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password.Password), bcrypt.DefaultCost)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to generate password hash")
		return
	}

	delete(resetTokens, token)

	user, _ := ac.UserRepository.GetUserByEmail(c, email)
	if user.ID == 0 {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find email")
		return
	}

	updatePassword := string(hashedPassword)

	err = ac.UserRepository.UpdateUserPassword(c, user, updatePassword)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to update user password")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"password": "Password successfully changed",
	})
}
