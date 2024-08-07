package information

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
)

type AppController struct {
	AppRepository models.AppRepository
}

// Create App Information godoc
// @Summary Create App Information
// @Security ApiKeyAuth
// @Tags admin-information-app-controller
// @ID create-app-information
// @Accept json
// @Produce json
// @Param name formData string true "name"
// @Param image formData file true "image"
// @Param appName formData string true "appName"
// @Param description formData string true "description"
// @Param version formData string true "version"
// @Param releaseDate formData string true "releaseDate"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/app [post]
func (h *AppController) CreateAppInformation(c *gin.Context) {
	var app models.App

	app.Name = c.PostForm("name")

	image, err := c.FormFile("image")
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse image param")
		return
	}

	path := "files//company//" + image.Filename
	c.SaveUploadedFile(image, path)
	app.Image = path

	app.AppName = c.PostForm("appName")
	app.Description = c.PostForm("description")
	app.Version = c.PostForm("version")
	app.ReleaseDate = c.PostForm("releaseDate")

	_, err = h.AppRepository.CreateAppInformation(c, app)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create app information")
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

// Update App Information godoc
// @Summary Update App Information
// @Security ApiKeyAuth
// @Tags admin-information-app-controller
// @ID update-app-information
// @Accept json
// @Produce json
// @Param name formData string true "name"
// @Param image formData file true "image"
// @Param appName formData string true "appName"
// @Param description formData string true "description"
// @Param version formData string true "version"
// @Param releaseDate formData string true "releaseDate"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/app [put]
func (h *AppController) UpdateAppInformation(c *gin.Context) {
	var updateApp models.App

	updateApp.Name = c.PostForm("name")

	image, err := c.FormFile("image")
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse image param")
		return
	}
	path := "files//company//" + image.Filename
	c.SaveUploadedFile(image, path)
	updateApp.Image = path

	updateApp.AppName = c.PostForm("appName")
	updateApp.Description = c.PostForm("description")
	updateApp.Version = c.PostForm("version")
	updateApp.ReleaseDate = c.PostForm("releaseDate")

	app, _ := h.AppRepository.GetAppInformation(c)
	if app.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find app information")
		return
	}

	err = h.AppRepository.UpdateAppInformation(c, app, updateApp)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to update app information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":        updateApp.Name,
		"Image":       updateApp.Image,
		"App Name":    updateApp.AppName,
		"Description": updateApp.Description,
		"Version":     updateApp.Version,
		"ReleaseDate": updateApp.ReleaseDate,
	})
}

// Get App Information godoc
// @Summary Get App Information
// @Security ApiKeyAuth
// @Tags admin-information-app-controller
// @ID get-app-information
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/app [get]
func (h *AppController) GetAppInformation(c *gin.Context) {
	app, _ := h.AppRepository.GetAppInformation(c)
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

// Delete App Information godoc
// @Summary Delete App Information
// @Security ApiKeyAuth
// @Tags admin-information-app-controller
// @ID delete-app-information
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/app [delete]
func (h *AppController) DeleteAppInformation(c *gin.Context) {
	app, _ := h.AppRepository.GetAppInformation(c)
	if app.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find app information")
		return
	}

	err := h.AppRepository.DeleteAppInformation(c, app)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to delete app information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"App": "App information deleted",
	})
}
