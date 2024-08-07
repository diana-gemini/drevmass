package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
)

// Get All Courses godoc
// @Summary Get All Courses
// @Security ApiKeyAuth
// @Tags user-course-controller
// @ID get-all-courses
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /course [get]
func (h *UserController) GetAllCourses(c *gin.Context) {
	courses, err := h.UserRepository.GetAllCourses(c)
	if err != nil {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find courses")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Courses": courses,
	})
}

// Get Course godoc
// @Summary Get Course
// @Security ApiKeyAuth
// @Tags user-course-controller
// @ID get-user-course
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /course/{id} [get]
func (h *UserController) GetCourse(c *gin.Context) {
	userID := c.GetUint("userID")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
		return
	}

	course, _ := h.UserRepository.GetCourseByID(c, uint(id), userID)
	if course.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find course")
		return
	}

	lessons, err := h.UserRepository.GetLessonsByCourseID(c, course.ID, userID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find lessons of course")
		return
	}

	isCourseStarted := false
	startCourseID, _ := h.UserRepository.GetUserStartCourse(c, userID, course.ID)
	if startCourseID != 0 {
		isCourseStarted = true
	}

	count, err := h.UserRepository.GetCountOfLessonsWatch(c, course.ID, userID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find count of watch")
		return
	}

	isCourseFinished := false
	if count.LessonCount == int(count.WatchCount) {
		err := h.UserRepository.UpdateUserFinishCourse(c, userID, course.ID)
		if err != nil {
			models.NewErrorResponse(c, http.StatusNotFound, "failed to update user course finish")
			return
		}

		point := models.Point{
			Description: "Начисление за прохождение видеокурса",
			Point:       course.Points,
			IsAvailable: true,
		}

		_, err = h.UserRepository.CreatePoint(c, uint(userID), point)
		if err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create points for finish course")
			return
		}

		isCourseFinished = true
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":            course.Name,
		"Lesson count":    course.LessonCount,
		"Time":            course.Time,
		"Description":     course.Description,
		"Image":           course.Image,
		"Points":          course.Points,
		"Lessons":         lessons,
		"Course started":  isCourseStarted,
		"Course finished": isCourseFinished,
		"Watch count":     course.WatchCount,
	})
}

// Start Course godoc
// @Summary Start Course
// @Security ApiKeyAuth
// @Tags user-course-controller
// @ID start-course
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /course/{id}/startcourse [post]
func (h *UserController) StartCourse(c *gin.Context) {
	userID := c.GetUint("userID")

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
		return
	}

	id, _ := h.UserRepository.CreateUserStartCourse(c, userID, uint(courseID))
	if id == 0 {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to create start course")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Start course id": id,
	})

}
