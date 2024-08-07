package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
	"gopkg.in/gomail.v2"
)

var resetTokens = make(map[string]string)

// PasswordRecover godoc
// @Summary Password Recover
// @Tags auth-controller
// @ID password-recover
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "email"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /passwordrecover [post]
func (ac *AuthController) PasswordRecover(c *gin.Context) {
	var request models.UserEmail

	if err := c.ShouldBind(&request); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse request parameters")
		return
	}

	user, _ := ac.UserRepository.GetUserByEmail(c, request.Email)
	if user.ID == 0 {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find user")
		return
	}

	token := generateToken()
	resetTokens[token] = request.Email

	if err := sendResetEmail(request.Email, token); err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to send reset email")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "We sent instructions for reset your password to " + request.Email,
	})
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func sendResetEmail(email, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "diana-test-project@mail.ru")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password Reset")
	m.SetBody("text/html", fmt.Sprintf("Click the following link to reset your password: <a href=\"http://localhost:3000/resetpassword?token=%s\">Reset Password</a>", token))

	d := gomail.NewDialer("smtp.mail.ru", 587, "diana-test-project@mail.ru", "phwEkPnEPudpnmU0Pvvn")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
