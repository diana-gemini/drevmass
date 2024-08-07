package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/diana-gemini/drevmass/internal/models"
)

type CourseRepository struct {
	db *pgxpool.Pool
}

func NewCourseRepository(db *pgxpool.Pool) models.CourseRepository {
	return &CourseRepository{db: db}
}

func (h *CourseRepository) CreateCourse(c context.Context, course models.Course) (int, error) {
	var courseID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `INSERT INTO courses(
		name, description, image, points, created_at)
		VALUES ($1, $2, $3, $4, $5) returning id;`
	err := h.db.QueryRow(c, userQuery, course.Name, course.Description, course.Image, course.Points, currentTime).Scan(&courseID)
	if err != nil {
		return 0, err
	}
	return courseID, nil
}

func (h *CourseRepository) UpdateCourseByID(c context.Context, id uint, updateCourse models.Course) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `UPDATE courses SET name=$1, description=$2, image=$3, points=$4, updated_at=$5 WHERE id=$6;`

	_, err := h.db.Exec(c, userQuery, updateCourse.Name, updateCourse.Description, updateCourse.Image, updateCourse.Points, currentTime, id)
	if err != nil {
		return err
	}
	return nil
}

func (h *CourseRepository) GetCourseByID(c context.Context, id uint) (models.Course, error) {
	course := models.Course{}

	query := `SELECT id, name, description, image, points FROM courses WHERE id=$1;`
	row := h.db.QueryRow(c, query, id)
	err := row.Scan(&course.ID, &course.Name, &course.Description, &course.Image, &course.Points)

	if err != nil {
		return course, err
	}

	return course, nil
}

func (h *CourseRepository) DeleteCourseByID(c context.Context, id uint) error {
	userQuery := `DELETE FROM courses WHERE id=$1;`
	_, err := h.db.Exec(c, userQuery, id)
	if err != nil {
		return err
	}
	return nil
}
