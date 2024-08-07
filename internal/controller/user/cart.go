package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
	"github.com/sirupsen/logrus"
)

// @Summary Create product
// @Security ApiKeyAuth
// @Tags user-cart-controller
// @Accept json
// @Produce json
// @Param id path string true "product id"
// @Success 200 {integer} integer
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /cart/{id} [post]
func (h *UserController) AddToCart(c *gin.Context) {

	userId := c.GetUint("userID")
	if userId == 0 {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find user")
		return
	}

	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	id, err := h.UserRepository.Add(c, userId, productId)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to add to cart")
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

// @Summary decrease the quantity by one
// @Security ApiKeyAuth
// @Tags user-cart-controller
// @Accept json
// @Produce json
// @Param id path string true "product id"
// @Success 200 {object} models.StatusResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /cart/minus/{id} [put]
func (h *UserController) Minus(c *gin.Context) {
	userId := c.GetUint("userID")
	if userId == 0 {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find user")
		return
	}

	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.UserRepository.Minus(c, userId, productId); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to decrease the quantity of product by one")
		return
	}
	c.JSON(http.StatusOK, models.StatusResponse{
		Status: "ok",
	})
}

// @Summary increase the quantity by one
// @Security ApiKeyAuth
// @Tags user-cart-controller
// @Accept json
// @Produce json
// @Param id path string true "product id"
// @Success 200 {object} models.StatusResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /cart/plus/{id} [put]
func (h *UserController) Plus(c *gin.Context) {
	userId := c.GetUint("userID")
	if userId == 0 {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find user")
		return
	}

	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.UserRepository.Plus(c, userId, productId); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to increase the quantity of product by one")
		return
	}
	c.JSON(http.StatusOK, models.StatusResponse{
		Status: "ok",
	})
}

// @Summary All in cart
// @Security ApiKeyAuth
// @Tags user-cart-controller
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /cart [get]
func (h *UserController) GetAllFromCart(c *gin.Context) {
	userId := c.GetUint("userID")
	if userId == 0 {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find user")
		return
	}

	products, err := h.UserRepository.GetAllFromCart(c, userId)
	if err != nil {
		logrus.Print(err)
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to get products from cart")
		return
	}
	totalAmount, err := h.UserRepository.GetTotalAmout(c, userId)
	if err != nil {
		logrus.Print(err)
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to get total amount from cart")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"products": products,
		"total":    totalAmount,
	})
}

// @Summary clear the cart
// @Security ApiKeyAuth
// @Tags user-cart-controller
// @Accept json
// @Produce json
// @Success 200 {object} models.StatusResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /cart [delete]
func (h *UserController) DeleteAllFromCart(c *gin.Context) {
	userId := c.GetUint("userID")
	if userId == 0 {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find user")
		return
	}

	if err := h.UserRepository.DeleteFromCart(c, userId); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to clear the cart")
		return
	}
	c.JSON(http.StatusOK, models.StatusResponse{
		Status: "ok",
	})

}
