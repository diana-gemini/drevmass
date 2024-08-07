package auth

import (
	"errors"
	"math/rand"
	"net/http"
	"strings"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	UserRepository models.UserRepository
}

// Signup godoc
// @Summary Signup
// @Tags auth-controller
// @ID create-account
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "username"
// @Param email formData string true "email"
// @Param mobile formData string true "mobile"
// @Param password formData string true "password"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /signup [post]
func (ac *AuthController) Signup(c *gin.Context) {
	var request models.SignupUser

	if err := c.ShouldBind(&request); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse request parameters")
		return
	}

	err := IsEmailExist(request.Email)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, _ := ac.UserRepository.GetUserByEmail(c, request.Email)
	if user.ID > 0 {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to register, email is already exist")
		return
	}

	if !CheckPassword(request.Password) {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to check password")
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)

	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to hash password")
		return
	}

	request.Password = string(hashPassword)
	request.Email = strings.ToLower(request.Email)

	userID, err := ac.UserRepository.CreateUser(c, request)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create user")
		return
	}

	point := models.Point{
		Description: "Начисление за регистрацию",
		Point:       500,
		IsAvailable: true,
	}

	_, err = ac.UserRepository.CreatePoint(c, uint(userID), point)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create point for signup")
		return
	}

	promocode := ac.generateUniquePromocode(c)

	_, err = ac.UserRepository.CreatePromocode(c, uint(userID), promocode)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create promocode")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"User": request,
	})
}

func CheckPassword(password string) bool {
	if password == "" {
		return false
	}
	if len(password) < 8 {
		return false
	}
	if len(password) > 50 {
		return false
	}
	return true
}

var (
	verifier = emailverifier.NewVerifier().
		EnableSMTPCheck().
		EnableDomainSuggest().
		DisableCatchAllCheck()
)

func IsEmailExist(email string) error {
	res, err := verifier.Verify(email)
	if err != nil {
		return errors.New("failed to verify email address")
	}

	// fmt.Println("email validation result", res)
	// fmt.Println("Email:", res.Email, "\nReachable:", res.Reachable, "\nSyntax:",
	// res.Syntax, "\nSMTP:", res.SMTP, "\nGravatar:", res.Gravatar, "\nSuggestion:",
	// res.Suggestion, "\nDisposable:", res.Disposable, "\nRoleAccount:", res.RoleAccount,
	// "\nFree:", res.Free, "\nHasMxRecords:", res.HasMxRecords)

	if !res.Syntax.Valid {
		return errors.New("email address syntax is invalid")
	}
	if res.Disposable {
		return errors.New("do not accept disposable email addresses")
	}
	if res.Suggestion != "" {
		return errors.New("email address is not reachable")
	}
	if res.Reachable == "no" {
		return errors.New("email address is not reachable")
	}
	if !res.HasMxRecords {
		return errors.New("domain entered not properly setup to recieve emails, MX record not found")
	}
	if !res.SMTP.Deliverable {
		return errors.New("email address is invalid, the email server is not to accept emails")
	}
	return nil
}

func (ac *AuthController) generateUniquePromocode(c *gin.Context) string {
	lenPromocode := 8
	for {
		promocode := generateRandomPromocode(lenPromocode)
		id, _ := ac.UserRepository.GetUniquePromocode(c, promocode)
		if id == 0 {
			return promocode
		}
	}
}

func generateRandomPromocode(n int) string {
	const symbols = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}
