package information

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
)

type CompanyController struct {
	CompanyRepository models.CompanyRepository
}

// Create Company Information godoc
// @Summary Create Company Information
// @Security ApiKeyAuth
// @Tags admin-information-company-controller
// @ID create-company-information
// @Accept json
// @Produce json
// @Param name formData string true "name"
// @Param description formData string true "description"
// @Param image formData file true "image"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/company [post]
func (h *CompanyController) CreateCompanyInformation(c *gin.Context) {
	var company models.Company

	company.Name = c.PostForm("name")
	company.Description = c.PostForm("description")

	image, err := c.FormFile("image")
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse image param")
		return
	}

	path := "files//company//" + image.Filename
	c.SaveUploadedFile(image, path)
	company.Image = path

	_, err = h.CompanyRepository.CreateCompanyInformation(c, company)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create company information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":        company.Name,
		"Description": company.Description,
		"Image":       company.Image,
	})
}

// Update Company Information godoc
// @Summary Update Company Information
// @Security ApiKeyAuth
// @Tags admin-information-company-controller
// @ID update-company-information
// @Accept json
// @Produce json
// @Param name formData string true "name"
// @Param description formData string true "description"
// @Param image formData file true "image"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/company [put]
func (h *CompanyController) UpdateCompanyInformation(c *gin.Context) {
	var updateCompany models.Company

	updateCompany.Name = c.PostForm("name")
	updateCompany.Description = c.PostForm("description")

	image, err := c.FormFile("image")
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse image param")
		return
	}
	path := "files//company//" + image.Filename
	c.SaveUploadedFile(image, path)
	updateCompany.Image = path

	company, _ := h.CompanyRepository.GetCompanyInformation(c)
	if company.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find company information")
		return
	}

	err = h.CompanyRepository.UpdateCompanyInformation(c, company, updateCompany)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create company information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":        updateCompany.Name,
		"Description": updateCompany.Description,
		"Image":       updateCompany.Image,
	})
}

// Get Company Information godoc
// @Summary Get Company Information
// @Security ApiKeyAuth
// @Tags admin-information-company-controller
// @ID get-company-information
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/company [get]
func (h *CompanyController) GetCompanyInformation(c *gin.Context) {
	company, _ := h.CompanyRepository.GetCompanyInformation(c)
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

// Delete Company Information godoc
// @Summary Delete Company Information
// @Security ApiKeyAuth
// @Tags admin-information-company-controller
// @ID delete-company-information
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/company [delete]
func (h *CompanyController) DeleteCompanyInformation(c *gin.Context) {
	company, _ := h.CompanyRepository.GetCompanyInformation(c)
	if company.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find company information")
		return
	}

	err := h.CompanyRepository.DeleteCompanyInformation(c, company)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to delete company information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Company": "Company information deleted",
	})
}
