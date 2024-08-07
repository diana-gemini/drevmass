package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
	"github.com/sirupsen/logrus"
)

// Get Promocode godoc
// @Summary Get Promocode
// @Security ApiKeyAuth
// @Tags user-profile-promocode-controller
// @ID get-user-promocode
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/promocode [get]
func (uc *UserController) GetPromocode(c *gin.Context) {
	userID := c.GetUint("userID")

	promocode, err := uc.UserRepository.GetPromocodeByUserID(c, userID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find promocode")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"promocode":       promocode.Promocode,
		"promocode count": promocode.Count,
	})
}

// Copy Promocode godoc
// @Summary Copy Promocode
// @Security ApiKeyAuth
// @Tags user-profile-promocode-controller
// @ID copy-user-promocode
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/promocode/copy [post]
func (uc *UserController) CopyPromocode(c *gin.Context) {
	userID := c.GetUint("userID")
	var count int

	promocode, err := uc.UserRepository.GetPromocodeByUserID(c, userID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find promocode")
		return
	}

	if promocode.Count < 2 {
		count = promocode.Count + 1
		updatePromocode := models.Promocode{
			Count: count,
		}

		err = uc.UserRepository.UpdatePromocodeByUserID(c, userID, updatePromocode)
		if err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, "failed to update promocode count")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"promocode count": count,
	})
}

// activate Promocode
// @Summary Activate Promocode
// @Security ApiKeyAuth
// @Tags user-profile-promocode-controller
// @ID activate-user-promocode
// @Param promocode formData string true "promocode"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /profile/promocode/activate [post]
func (uc *UserController) ActivatePromocode(c *gin.Context) {
	userID := c.GetUint("userID")
	var count int
	promo := c.PostForm("promocode")

	points, err := uc.UserRepository.GetPointsByUserID(c, userID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find points by user id")
		return
	}

	for _, point := range points {
		if point.Description == "Начисление за промокод друга" {
			models.NewErrorResponse(c, http.StatusBadRequest, "you already used the promocode")
			return
		}
	}

	promocode, err := uc.UserRepository.GetPromocodeByPromocode(c, promo)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find promocode")
		return
	}

	if userID == promocode.UserID {
		models.NewErrorResponse(c, http.StatusInternalServerError, "you cannot to use your promo")
		return
	}

	if promocode.Count < 2 {
		count = promocode.Count + 1
		updatePromocode := models.Promocode{
			Count: count,
		}

		err = uc.UserRepository.UpdatePromocodeByUserID(c, promocode.UserID, updatePromocode)
		if err != nil {
			logrus.Print(err.Error())
			models.NewErrorResponse(c, http.StatusInternalServerError, "failed to update promocode count")
			return
		}
	} else {
		models.NewErrorResponse(c, http.StatusInternalServerError, "the promocode was already used")
		return
	}

	point := models.Point{
		Description: "Начисление за промокод друга",
		Point:       2500,
		IsAvailable: true,
		PromocodeId: promocode.ID,
	}

	_, err = uc.UserRepository.CreatePointForPromo(c, uint(userID), point)
	if err != nil {
		logrus.Print(err.Error())
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create point for promocode")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"promocode count": count,
	})
}
