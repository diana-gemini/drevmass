package information

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
)

type BonusController struct {
	BonusRepository models.BonusRepository
}

// Create Bonus Information godoc
// @Summary Create Bonus Information
// @Security ApiKeyAuth
// @Tags admin-information-bonus-controller
// @ID create-bonus-information
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "name"
// @Param description formData string true "description"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/bonus [post]
func (h *BonusController) CreateBonusInformation(c *gin.Context) {
	var bonus models.Bonus

	if err := c.ShouldBind(&bonus); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse request parameters")
		return
	}

	fmt.Printf("bonus - %v \n", bonus)

	_, err := h.BonusRepository.CreateBonusInformation(c, bonus)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create bonus information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":        bonus.Name,
		"Description": bonus.Description,
	})
}

// Update Bonus Information godoc
// @Summary Update Bonus Information
// @Security ApiKeyAuth
// @Tags admin-information-bonus-controller
// @ID update-bonus-information
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "name"
// @Param description formData string true "description"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/bonus [put]
func (h *BonusController) UpdateBonusInformation(c *gin.Context) {
	var updateBonus models.Bonus

	if err := c.ShouldBind(&updateBonus); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse request parameters")
		return
	}

	bonus, _ := h.BonusRepository.GetBonusInformation(c)
	if bonus.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find bonus information")
		return
	}

	err := h.BonusRepository.UpdateBonusInformation(c, bonus, updateBonus)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create bonus information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":        updateBonus.Name,
		"Description": updateBonus.Description,
	})
}

// Get Bonus Information godoc
// @Summary Get Bonus Information
// @Security ApiKeyAuth
// @Tags admin-information-bonus-controller
// @ID get-bonus-information
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/bonus [get]
func (h *BonusController) GetBonusInformation(c *gin.Context) {
	bonus, _ := h.BonusRepository.GetBonusInformation(c)
	if bonus.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find bonus information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":        bonus.Name,
		"Description": bonus.Description,
	})
}

// Delete Bonus Information godoc
// @Summary Delete Bonus Information
// @Security ApiKeyAuth
// @Tags admin-information-bonus-controller
// @ID delete-bonus-information
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/information/bonus [delete]
func (h *BonusController) DeleteBonusInformation(c *gin.Context) {
	bonus, _ := h.BonusRepository.GetBonusInformation(c)
	if bonus.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find bonus information")
		return
	}

	err := h.BonusRepository.DeleteBonusInformation(c, bonus)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to delete bonus information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Bonus": "Bonus information deleted",
	})
}
