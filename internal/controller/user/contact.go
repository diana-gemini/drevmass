package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
)

// Get Call godoc
// @Summary Get Call
// @Security ApiKeyAuth
// @Tags user-profile-contact-controller
// @ID get-call
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/contact/call [get]
func (uc *UserController) GetCall(c *gin.Context) {
	contact, _ := uc.UserRepository.GetContactInformation(c)
	if contact.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find contact information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Mobile": contact.Name,
	})
}

// Send To Support godoc
// @Summary Send To Support
// @Security ApiKeyAuth
// @Tags user-profile-contact-controller
// @ID send-to-support
// @Accept multipart/form-data
// @Produce json
// @Param message formData string true "message"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/contact/support [post]
func (uc *UserController) SendToSupport(c *gin.Context) {
	message := c.PostForm("message")

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

// Get Whatsapp godoc
// @Summary Get Whatsapp
// @Security ApiKeyAuth
// @Tags user-profile-contact-controller
// @ID get-whatsapp
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/contact/whatsapp [get]
func (uc *UserController) GetWhatsapp(c *gin.Context) {
	contact, _ := uc.UserRepository.GetContactInformation(c)
	if contact.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find contact information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Whatsapp": contact.Name,
	})
}
