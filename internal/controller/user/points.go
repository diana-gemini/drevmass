package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
)

// Get Points godoc
// @Summary Get Points
// @Security ApiKeyAuth
// @Tags user-profile-points-controller
// @ID get-user-points
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/points [get]
func (uc *UserController) GetPoints(c *gin.Context) {
	userID := c.GetUint("userID")

	user, _ := uc.UserRepository.GetUserByID(c, userID)
	if user.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find user")
		return
	}

	sum, err := uc.UserRepository.GetSumPointsByUserID(c, user.ID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to sum user points")
		return
	}

	expiredPoint, err := uc.UserRepository.GetExpiredPointByUserID(c, user.ID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to get user expired point")
		return
	}

	points, err := uc.UserRepository.GetPointsByUserID(c, user.ID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to get user points")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sum":           sum,
		"expired point": expiredPoint.AvailablePoint,
		"expired data":  expiredPoint.ExpiredAt,
		"points":        points,
	})
}

// Get About Bonus godoc
// @Summary Get About Bonus
// @Security ApiKeyAuth
// @Tags user-profile-points-controller
// @ID get-about-bonus
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/points/about [get]
func (uc *UserController) AboutBonus(c *gin.Context) {
	bonus, _ := uc.UserRepository.GetBonusInformation(c)
	if bonus.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find bonus information")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":        bonus.Name,
		"Description": bonus.Description,
	})
}
