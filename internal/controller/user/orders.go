package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
	"github.com/sirupsen/logrus"
)

// @Summary make order
// @Security ApiKeyAuth
// @Tags admin-order-controller
// @Accept json
// @Produce json
// @Success 200 {integer} integer
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /order [post]
func (h *UserController) CreateOrder(c *gin.Context) {
	userId := c.GetUint("userID")
	if userId == 0 {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find user")
		return
	}

	totalAmount, err := h.UserRepository.GetTotalAmout(c, userId)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find total amount")
		return
	}

	products, err := h.UserRepository.GetAllFromCart(c, userId)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to find products")
		return
	}

	productsStr := ""
	for _, v := range products {
		str := fmt.Sprintf("%s %d  %f", v.Name, v.Quantity, v.Price)
		productsStr = productsStr + "\n" + str
	}

	err = h.UserRepository.CheckExpiredUserPoints(c, userId)

	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to check bonus")
		return
	}

	bonuses, err := h.UserRepository.GetAvailablePointsByUserID(c, userId)
	if err != nil {
		logrus.Print(err.Error())
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to get bonus")
		return
	}

	percent30 := totalAmount * 30 / 100

	sum := 0
	for _, v := range bonuses {
		sum += v.AvailablePoint
	}

	bonus := 0.0
	if percent30 >= float64(sum) {
		bonus = float64(sum)
	} else {
		bonus = percent30
	}

	overall := totalAmount - bonus

	id, err := h.UserRepository.CreateOrder(c, int(userId), totalAmount, bonus, overall, products)
	if err != nil {
		logrus.Print(err.Error())
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to create the order")
		return
	}
	user, err := h.UserRepository.GetUserInfoByID(c, userId)
	if err != nil {
		logrus.Print(err.Error())
		models.NewErrorResponse(c, http.StatusBadRequest, "failed find client")
		return
	}

	userStr := fmt.Sprintf("%s %s  %s", user.Username, user.Email, user.Mobile)

	err = SendDealFunc(productsStr, userStr, int(totalAmount), bonus, overall)
	if err != nil {
		logrus.Print(err.Error())
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to send to crm")
		return
	}

	for _, v := range bonuses {
		if bonus > 0 && v.Description == "Начисление за промокод друга" && v.PromocodeId > 0 {
			promo, err := h.UserRepository.GetPromocodeById(c, uint(v.PromocodeId))
			if err != nil {
				logrus.Print(err.Error())
				models.NewErrorResponse(c, http.StatusBadRequest, "failed get user of promo")
				return
			}
			logrus.Print(promo, promo.UserID)
			point := models.Point{
				Description: "Начисление за воспользовние твоего промокода",
				Point:       500,
				IsAvailable: true,
			}

			err = h.UserRepository.CreatePointUsedPromo(c, promo.UserID, point)
			if err != nil {
				logrus.Print(err.Error())
				models.NewErrorResponse(c, http.StatusBadRequest, "failed get user of promo")
				return
			}
		}
		if bonus > 0 && float64(v.AvailablePoint) <= bonus {
			// update AvailablePoint = 0
			err := h.UserRepository.UpdateAvailablePointsZero(c, v.PointID, userId)
			if err != nil {
				logrus.Print(err.Error())
				models.NewErrorResponse(c, http.StatusBadRequest, "failed to update bonus")
				return
			}
			bonus -= float64(v.AvailablePoint)
		} else if bonus > 0 && float64(v.AvailablePoint) > bonus {
			// update AvailablePoint = AvailablePoint - bonus

			err := h.UserRepository.UpdateAvailablePoints(c, v.PointID, userId, bonus)
			if err != nil {
				logrus.Print(err.Error())
				models.NewErrorResponse(c, http.StatusBadRequest, "failed to update bonus")
				return
			}
			bonus = 0
		}
	}

	if err := h.UserRepository.DeleteFromCart(c, userId); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to clear the cart")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary all orders
// @Security ApiKeyAuth
// @Tags admin-order-controller
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/order [get]
func (h *UserController) GetAllOrders(c *gin.Context) {
	orders, err := h.UserRepository.GetAllOrders(c)
	if err != nil {
		logrus.Print(err.Error())
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to get all orders")
		return
	}

	orderMap := make(map[int]*models.OrderWithProducts)
	for _, order := range orders {
		if _, found := orderMap[order.Id]; !found {

			orderMap[order.Id] = &models.OrderWithProducts{
				Id:          order.Id,
				TotalAmount: order.TotalAmount,
				Bonus:       order.Bonus,
				Overall:     order.Overall,
				CreatedAt:   order.CreatedAt,
				UserId:      order.UserId,
			}

		}

		product := models.GetProductsFromCart{
			Id:       order.ProductId,
			Name:     order.Name,
			Image:    order.Image,
			Price:    order.Price,
			Quantity: order.Quantity,
		}
		orderMap[order.Id].Products = append(orderMap[order.Id].Products, product)

	}

	// Convert map to slice
	result := make([]models.OrderWithProducts, 0, len(orderMap))
	for _, v := range orderMap {
		result = append(result, *v)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"orders": result,
	})

}
