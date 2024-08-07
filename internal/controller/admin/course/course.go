package course

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
)

type CourseController struct {
	CourseRepository models.CourseRepository
}

// Create Course godoc
// @Summary Create Course
// @Security ApiKeyAuth
// @Tags admin-course-controller
// @ID create-course
// @Accept json
// @Produce json
// @Param name formData string true "name"
// @Param description formData string true "description"
// @Param image formData file true "image"
// @Param points formData string true "points"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/course [post]
func (h *CourseController) CreateCourse(c *gin.Context) {
	var course models.Course

	course.Name = c.PostForm("name")
	course.Description = c.PostForm("description")

	image, err := c.FormFile("image")
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse image param")
		return
	}

	path := "files//course//" + image.Filename
	c.SaveUploadedFile(image, path)
	course.Image = path

	points, err := strconv.Atoi(c.PostForm("points"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
		return
	}
	course.Points = points

	_, err = h.CourseRepository.CreateCourse(c, course)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create course")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":        course.Name,
		"Description": course.Description,
		"Image":       course.Image,
		"Points":      course.Points,
	})
}

// Update Course godoc
// @Summary Update Course
// @Security ApiKeyAuth
// @Tags admin-course-controller
// @ID update-course
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param name formData string true "name"
// @Param description formData string true "description"
// @Param image formData file true "image"
// @Param points formData string true "points"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/course/{id} [put]
func (h *CourseController) UpdateCourse(c *gin.Context) {
	var updateCourse models.Course

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
		return
	}

	updateCourse.Name = c.PostForm("name")
	updateCourse.Description = c.PostForm("description")

	image, err := c.FormFile("image")
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse image param")
		return
	}
	path := "files//course//" + image.Filename
	c.SaveUploadedFile(image, path)
	updateCourse.Image = path

	points, err := strconv.Atoi(c.PostForm("points"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
		return
	}
	updateCourse.Points = points

	course, _ := h.CourseRepository.GetCourseByID(c, uint(id))
	if course.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find course")
		return
	}

	err = h.CourseRepository.UpdateCourseByID(c, course.ID, updateCourse)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to update course")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":        updateCourse.Name,
		"Description": updateCourse.Description,
		"Image":       updateCourse.Image,
		"Points":      updateCourse.Points,
	})
}

// Get Course godoc
// @Summary Get Course
// @Security ApiKeyAuth
// @Tags admin-course-controller
// @ID get-course
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/course/{id} [get]
func (h *CourseController) GetCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
	}

	course, _ := h.CourseRepository.GetCourseByID(c, uint(id))
	if course.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find course")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":        course.Name,
		"Description": course.Description,
		"Image":       course.Image,
		"Points":      course.Points,
	})
}

// Delete Course godoc
// @Summary Delete Course
// @Security ApiKeyAuth
// @Tags admin-course-controller
// @ID delete-course
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/course/{id} [delete]
func (h *CourseController) DeleteCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
	}
	course, _ := h.CourseRepository.GetCourseByID(c, uint(id))
	if course.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find course")
		return
	}

	err = h.CourseRepository.DeleteCourseByID(c, course.ID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to delete course")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Course": "Course deleted",
	})
}
