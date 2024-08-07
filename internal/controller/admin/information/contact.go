package information

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
)

type ContactController struct {
	ContactRepository models.ContactRepository
}

// Create Contact Information godoc
// @Summary Create Contact Information
// @Security ApiKeyAuth
// @Tags admin-information-contact-controller
// @ID create-contact-information
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "name"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/contact [post]
func (h *ContactController) CreateContactInformation(c *gin.Context) {
	var contact models.Contact

	contact.Name = c.PostForm("name")

	_, err := h.ContactRepository.CreateContactInformation(c, contact)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create contact information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name": contact.Name,
	})
}

// Update Contact Information godoc
// @Summary Update Contact Information
// @Security ApiKeyAuth
// @Tags admin-information-contact-controller
// @ID update-contact-information
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "name"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/contact [put]
func (h *ContactController) UpdateContactInformation(c *gin.Context) {
	var updateContact models.Contact

	updateContact.Name = c.PostForm("name")

	contact, _ := h.ContactRepository.GetContactInformation(c)
	if contact.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find contact information")
		return
	}

	err := h.ContactRepository.UpdateContactInformation(c, contact, updateContact)
	fmt.Printf("err - %v \n", err)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to update contact information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name": updateContact.Name,
	})
}

// Get Contact Information godoc
// @Summary Get Contact Information
// @Security ApiKeyAuth
// @Tags admin-information-contact-controller
// @ID get-contact-information
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/contact [get]
func (h *ContactController) GetContactInformation(c *gin.Context) {
	contact, _ := h.ContactRepository.GetContactInformation(c)
	if contact.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find contact information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name": contact.Name,
	})
}

// Delete Contact Information godoc
// @Summary Delete Contact Information
// @Security ApiKeyAuth
// @Tags admin-information-contact-controller
// @ID delete-contact-information
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/contact [delete]
func (h *ContactController) DeleteContactInformation(c *gin.Context) {
	contact, _ := h.ContactRepository.GetContactInformation(c)
	if contact.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find contact information")
		return
	}

	err := h.ContactRepository.DeleteContactInformation(c, contact)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to delete contact information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Contact": "Contact information deleted",
	})
}
