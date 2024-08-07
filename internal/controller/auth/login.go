package auth

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/diana-gemini/drevmass/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary Login
// @Tags auth-controller
// @ID login
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "email"
// @Param password formData string true "password"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var request models.LoginUser

	if err := c.ShouldBind(&request); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse request parameters")
		return
	}

	request.Email = strings.ToLower(request.Email)

	user, _ := ac.UserRepository.GetUserByEmail(c, request.Email)
	if user.ID == 0 {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find user")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to verify email or password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(),
		"role": user.RoleID,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to create token")
		return
	}

	err = ac.UserRepository.CheckExpiredUserPoints(c, user.ID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to check user expired points")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Token": tokenString,
	})
}
