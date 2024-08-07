package models

import "context"

type Lesson struct {
	ID           uint   `json:"id"`
	LessonNumber int    `json:"lessonNumber" form:"lessonNumber"`
	Name         string `json:"name" form:"name"`
	Description  string `json:"description" form:"description"`
	Image        string `json:"image" form:"image"`
	Video        string `json:"video" form:"video"`
	VideoTime    string `json:"videoTime" form:"videoTime"`
	Time         int    `json:"time" form:"time"`
	CourseID     int    `json:"courseID" form:"courseID"`
	BookmarkID   uint   `json:"bookmarkID"`
	IsWatched    bool   `json:"isWatched"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type LessonRepository interface {
	CreateLesson(c context.Context, lesson Lesson) (int, error)
	UpdateLessonByID(c context.Context, id uint, updateLesson Lesson) error
	GetLessonByID(c context.Context, id uint) (Lesson, error)
	DeleteLessonByID(c context.Context, id uint) error
	GetCourseByID(c context.Context, id uint) (Course, error)
	CreateLessonProduct(c context.Context, lessonID uint, productID int) (int, error)
	GetLessonProductsByID(c context.Context, lessonID uint) ([]int, error)
	DeleteLessonProductsByID(c context.Context, id uint) error
	IsLessonNumberExist(c context.Context, courseID uint, lessonNumber int) (int, error)
}
