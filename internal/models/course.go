package models

import "context"

type Course struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	Image       string `json:"image" form:"image"`
	Points      int    `json:"points" form:"points"`
	LessonCount int    `json:"lessonCount"`
	WatchCount  int64  `json:"watchCount"`
	Time        int    `json:"time"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type Bookmark struct {
	ID         uint   `json:"id"`
	UserID     uint   `json:"userID"`
	CourseID   uint   `json:"courseID"`
	CourseName string `json:"courseName"`
	Lessons    []Lesson
}

type CourseRepository interface {
	CreateCourse(c context.Context, course Course) (int, error)
	UpdateCourseByID(c context.Context, id uint, updateCourse Course) error
	GetCourseByID(c context.Context, id uint) (Course, error)
	DeleteCourseByID(c context.Context, id uint) error
}
