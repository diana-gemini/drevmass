package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
)

// Get Lesson godoc
// @Summary Get Lesson
// @Security ApiKeyAuth
// @Tags user-course-lesson-controller
// @ID get-user-lesson
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param lessonID path string true "lessonID"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /course/{id}/lesson/{lessonID} [get]
func (h *UserController) GetLesson(c *gin.Context) {
	userID := c.GetUint("userID")

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
		return
	}

	lessonID, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
		return
	}

	lesson, _ := h.UserRepository.GetLessonByID(c, uint(courseID), uint(lessonID))
	if lesson.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find lesson")
		return
	}

	var isBookmark bool
	bookmark, _ := h.UserRepository.GetUserBookmarkByIDs(c, userID, uint(courseID), uint(lessonID))

	if bookmark != 0 {
		isBookmark = true
	}

	products, err := h.UserRepository.GetLessonProducts(c, lesson.ID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find lesson products")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Lesson Number": lesson.LessonNumber,
		"Image":         lesson.Image,
		"Video":         lesson.Video,
		"Time":          lesson.Time,
		"Is bookmark":   isBookmark,
		"Name":          lesson.Name,
		"Description":   lesson.Description,
		"Products":      products,
	})
}

// Show Video godoc
// @Summary Show Video
// @Security ApiKeyAuth
// @Tags user-course-lesson-controller
// @ID show-video
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param lessonID path string true "lessonID"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /course/{id}/lesson/{lessonID}/video [post]
func (h *UserController) ShowVideo(c *gin.Context) {
	userID := c.GetUint("userID")

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse courseID param")
		return
	}

	lessonID, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse lessonID param")
		return
	}

	id, _ := h.UserRepository.CreateOfLessonsWatch(c, userID, uint(courseID), uint(lessonID))
	if id == 0 {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to create user watch lesson")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Watch lesson id": id,
	})
}

// Set Bookmark godoc
// @Summary Set Bookmark
// @Security ApiKeyAuth
// @Tags user-course-lesson-controller
// @ID set-bookmark
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param lessonID path string true "lessonID"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /course/{id}/lesson/{lessonID}/bookmark [post]
func (h *UserController) SetBookmark(c *gin.Context) {
	userID := c.GetUint("userID")

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
		return
	}

	lessonID, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse lessonID param")
		return
	}
	var isBookmark bool
	var bookmarkID uint

	bookmark, _ := h.UserRepository.GetUserBookmarkByIDs(c, userID, uint(courseID), uint(lessonID))
	if bookmark == 0 {
		bookmarkID, err = h.UserRepository.CreateUserBookmark(c, userID, uint(courseID), uint(lessonID))
		if err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, "failed to create bookmark")
			return
		}
	} else {
		err := h.UserRepository.DeleteUserBookmark(c, userID, uint(courseID), uint(lessonID))
		if err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, "failed to delete bookmark")
			return
		}
	}

	if bookmarkID != 0 {
		isBookmark = true
	}
	c.JSON(http.StatusOK, gin.H{
		"Is bookmark": isBookmark,
	})
}
