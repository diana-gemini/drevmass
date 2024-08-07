package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
)

// Get Company Information godoc
// @Summary Get Company Information
// @Security ApiKeyAuth
// @Tags user-profile-information-controller
// @ID get-company-information-user
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/information/company [get]
func (uc *UserController) GetCompanyInformation(c *gin.Context) {
	company, _ := uc.UserRepository.GetCompanyInformation(c)
	if company.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find company information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":        company.Name,
		"Description": company.Description,
		"Image":       company.Image,
	})
}

// Get App Information godoc
// @Summary Get App Information
// @Security ApiKeyAuth
// @Tags user-profile-information-controller
// @ID get-app-information-user
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/information/app [get]
func (uc *UserController) GetAppInformation(c *gin.Context) {
	app, _ := uc.UserRepository.GetAppInformation(c)
	if app.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find app information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":        app.Name,
		"Image":       app.Image,
		"App Name":    app.AppName,
		"Description": app.Description,
		"Version":     app.Version,
		"ReleaseDate": app.ReleaseDate,
	})
}