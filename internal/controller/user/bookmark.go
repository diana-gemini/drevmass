package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
)

// Get User Bookmarks godoc
// @Summary Get User Bookmarks
// @Security ApiKeyAuth
// @Tags user-course-controller
// @ID get-user-bookmarks
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /course/bookmarks [get]
func (h *UserController) GetUserBookmarks(c *gin.Context) {
	userID := c.GetUint("userID")

	bookmarks, err := h.UserRepository.GetUserBookmarks(c, userID)
	fmt.Printf("err - %v \n", err)
	if err != nil {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find user bookmarks")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Bookmarks": bookmarks,
	})

}
