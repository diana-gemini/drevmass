package course

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/diana-gemini/drevmass/internal/models"
)

type LessonController struct {
	LessonRepository models.LessonRepository
}

// Create Lesson godoc
// @Summary Create Lesson
// @Security ApiKeyAuth
// @Tags admin-course-lesson-controller
// @ID create-lesson
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param lessonNumber formData string true "lessonNumber"
// @Param name formData string true "name"
// @Param description formData string true "description"
// @Param image formData file false "image"
// @Param video formData string true "video"
// @Param time formData string true "time in the format 'X min Y sec'"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/course/{id}/lesson [post]
func (h *LessonController) CreateLesson(c *gin.Context) {
	var lesson models.Lesson

	number, err := strconv.Atoi(c.PostForm("lessonNumber"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse lesson number param")
		return
	}
	lesson.LessonNumber = number

	lesson.Name = c.PostForm("name")
	lesson.Description = c.PostForm("description")
	lesson.Video = c.PostForm("video")

	image, err := c.FormFile("image")
	if err == nil {
		path := "files//course//" + image.Filename
		c.SaveUploadedFile(image, path)
		lesson.Image = path
	} else {
		imageURL, err := getYouTubeImage(lesson.Video)
		if err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, "failed to get image url from youtube")
			return
		}
		lesson.Image = imageURL
	}

	lesson.VideoTime = c.PostForm("time")

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse courseID param")
		return
	}

	course, _ := h.LessonRepository.GetCourseByID(c, uint(courseID))
	if course.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find course")
		return
	}

	lesson.CourseID = courseID

	lessonNumber, _ := h.LessonRepository.IsLessonNumberExist(c, uint(courseID), number)
	if lessonNumber != 0 {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create lesson, lesson number for this course is already exist")
		return
	}

	_, err = h.LessonRepository.CreateLesson(c, lesson)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create lesson")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Lesson Number": lesson.LessonNumber,
		"Name":          lesson.Name,
		"Description":   lesson.Description,
		"Image":         lesson.Image,
		"Video":         lesson.Video,
		"Video Time":    lesson.VideoTime,
		"CourseID":      lesson.CourseID,
	})
}

// Update Lesson godoc
// @Summary Update Lesson
// @Security ApiKeyAuth
// @Tags admin-course-lesson-controller
// @ID update-lesson
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "id"
// @Param lessonID path string true "lessonID"
// @Param lessonNumber formData string true "lessonNumber"
// @Param name formData string true "name"
// @Param description formData string true "description"
// @Param image formData file false "image"
// @Param video formData string true "video"
// @Param time formData string true "time in the format 'X min Y sec'"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/course/{id}/lesson/{lessonID} [put]
func (h *LessonController) UpdateLesson(c *gin.Context) {
	var updateLesson models.Lesson

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse courseID param")
		return
	}

	lessonID, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
		return
	}

	number, err := strconv.Atoi(c.PostForm("lessonNumber"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse lesson number param")
		return
	}
	updateLesson.LessonNumber = number

	updateLesson.Name = c.PostForm("name")
	updateLesson.Description = c.PostForm("description")

	updateLesson.Video = c.PostForm("video")

	image, err := c.FormFile("image")
	if err == nil {
		path := "files//course//" + image.Filename
		c.SaveUploadedFile(image, path)
		updateLesson.Image = path
	} else {
		imageURL, err := getYouTubeImage(updateLesson.Video)
		if err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, "failed to get image url from youtube")
			return
		}
		updateLesson.Image = imageURL
	}

	updateLesson.VideoTime = c.PostForm("time")

	course, _ := h.LessonRepository.GetCourseByID(c, uint(courseID))
	if course.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find course")
		return
	}

	updateLesson.CourseID = courseID

	lesson, _ := h.LessonRepository.GetLessonByID(c, uint(lessonID))
	if lesson.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find lesson")
		return
	}

	err = h.LessonRepository.UpdateLessonByID(c, lesson.ID, updateLesson)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to update lesson")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":        updateLesson.Name,
		"Description": updateLesson.Description,
		"Image":       updateLesson.Image,
		"Video":       updateLesson.Video,
		"Time":        updateLesson.Time,
		"CourseID":    updateLesson.CourseID,
	})
}

// Get Lesson godoc
// @Summary Get Lesson
// @Security ApiKeyAuth
// @Tags admin-course-lesson-controller
// @ID get-lesson
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param lessonID path string true "lessonID"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/course/{id}/lesson/{lessonID} [get]
func (h *LessonController) GetLesson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
	}

	lesson, _ := h.LessonRepository.GetLessonByID(c, uint(id))
	if lesson.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find lesson")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":        lesson.Name,
		"Description": lesson.Description,
		"Image":       lesson.Image,
		"Video":       lesson.Video,
		"Time":        lesson.Time,
		"CourseID":    lesson.CourseID,
	})
}

// Delete Lesson godoc
// @Summary Delete Lesson
// @Security ApiKeyAuth
// @Tags admin-course-lesson-controller
// @ID delete-lesson
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param lessonID path string true "lessonID"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/course/{id}/lesson/{lessonID} [delete]
func (h *LessonController) DeleteLesson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
	}
	lesson, _ := h.LessonRepository.GetLessonByID(c, uint(id))
	if lesson.ID == 0 {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find lesson")
		return
	}

	err = h.LessonRepository.DeleteLessonByID(c, lesson.ID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to delete lesson")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Lesson": "Lesson deleted",
	})
}

// Create Lesson Products godoc
// @Summary Create Lesson Products
// @Security ApiKeyAuth
// @Tags admin-course-lesson-controller
// @ID create-lesson-products
// @Accept multipart/form-data
// @Produce json
// @Param lessonID path string true "lessonID"
// @Param products formData []string true "products"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/course/{id}/lesson/{lessonID}/products [post]
func (h *LessonController) CreateLessonProducts(c *gin.Context) {
	lessonID, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse lessonID param")
		return
	}

	products := c.PostFormArray("products")

	var productIDs []string

	if len(products) > 0 {
		productIDs = strings.Split(products[0], ",")
		for _, productID := range productIDs {
			id, err := strconv.Atoi(string(productID))
			if err != nil {
				models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse productID in productIDs")
				return
			}
			_, err = h.LessonRepository.CreateLessonProduct(c, uint(lessonID), id)
			if err != nil {
				models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create lesson product")
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Lesson ID":       lessonID,
		"Lesson products": productIDs,
	})
}

// Get Lesson Products godoc
// @Summary Get Lesson Products
// @Security ApiKeyAuth
// @Tags admin-course-lesson-controller
// @ID get-lesson-products
// @Accept json
// @Produce json
// @Param lessonID path string true "lessonID"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/course/{id}/lesson/{lessonID}/products [get]
func (h *LessonController) GetLessonProducts(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
	}

	products, err := h.LessonRepository.GetLessonProductsByID(c, uint(id))
	if err != nil {
		models.NewErrorResponse(c, http.StatusNotFound, "failed to find lesson products")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Products": products,
	})
}

// Update Lesson Products godoc
// @Summary Update Lesson Products
// @Security ApiKeyAuth
// @Tags admin-course-lesson-controller
// @ID update-lesson-products
// @Accept multipart/form-data
// @Produce json
// @Param lessonID path string true "lessonID"
// @Param products formData []string true "products"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/course/{id}/lesson/{lessonID}/products [put]
func (h *LessonController) UpdateLessonProducts(c *gin.Context) {
	lessonID, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse lessonID param")
		return
	}

	products := c.PostFormArray("products")

	err = h.LessonRepository.DeleteLessonProductsByID(c, uint(lessonID))
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to delete lesson products")
		return
	}

	var productIDs []string

	if len(products) > 0 {
		productIDs = strings.Split(products[0], ",")
		for _, productID := range productIDs {
			id, err := strconv.Atoi(string(productID))
			if err != nil {
				models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse productID in productIDs")
				return
			}
			_, err = h.LessonRepository.CreateLessonProduct(c, uint(lessonID), id)
			if err != nil {
				models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create lesson product")
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Lesson ID":       lessonID,
		"Lesson products": productIDs,
	})
}

// Delete Lesson Products godoc
// @Summary Delete Lesson Products
// @Security ApiKeyAuth
// @Tags admin-course-lesson-controller
// @ID delete-lesson-products
// @Accept json
// @Produce json
// @Param lessonID path string true "lessonID"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /admin/course/{id}/lesson/{lessonID}/products [delete]
func (h *LessonController) DeleteLessonProducts(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "failed to parse id param")
	}

	err = h.LessonRepository.DeleteLessonProductsByID(c, uint(id))
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to delete lesson products")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Lesson products": "Lesson products deleted",
	})
}

func getYouTubeImage(video string) (string, error) {
	videoID := extractYouTubeVideoID(video)
	if videoID == "" {
		return "", errors.New("failed to parse video url")
	}

	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/default.jpg", videoID)
	return thumbnailURL, nil
}

func extractYouTubeVideoID(video string) string {
	u, err := url.Parse(video)
	if err != nil {
		return ""
	}

	var videoID string
	switch u.Host {
	case "youtu.be":
		videoID = strings.TrimLeft(u.Path, "/")
	case "www.youtube.com", "youtube.com":
		query, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			return ""
		}
		videoID = query.Get("v")
	default:
		return ""
	}

	return videoID
}
