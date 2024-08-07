package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/diana-gemini/drevmass/internal/models"
)

type LessonRepository struct {
	db *pgxpool.Pool
}

func NewLessonRepository(db *pgxpool.Pool) models.LessonRepository {
	return &LessonRepository{db: db}
}

func (h *LessonRepository) CreateLesson(c context.Context, lesson models.Lesson) (int, error) {
	var lessonID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `INSERT INTO lessons(
		lesson_number, name, description, image, video, time, course_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) returning id;`
	err := h.db.QueryRow(c, userQuery, lesson.LessonNumber, lesson.Name, lesson.Description, lesson.Image, lesson.Video, lesson.VideoTime, lesson.CourseID, currentTime).Scan(&lessonID)
	if err != nil {
		return 0, err
	}
	return lessonID, nil
}

func (h *LessonRepository) UpdateLessonByID(c context.Context, id uint, updateLesson models.Lesson) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `UPDATE lessons SET lesson_number=$1, name=$2, description=$3, image=$4, 
				video=$5, time=$6, course_id=$7, updated_at=$8 WHERE id=$9;`

	_, err := h.db.Exec(c, userQuery, updateLesson.LessonNumber, updateLesson.Name, updateLesson.Description, updateLesson.Image, updateLesson.Video, updateLesson.VideoTime, updateLesson.CourseID, currentTime, id)
	if err != nil {
		return err
	}
	return nil
}

func (h *LessonRepository) GetLessonByID(c context.Context, id uint) (models.Lesson, error) {
	lesson := models.Lesson{}

	query := `SELECT id, name, description, image, video, time, course_id FROM lessons WHERE id=$1;`
	row := h.db.QueryRow(c, query, id)
	err := row.Scan(&lesson.ID, &lesson.Name, &lesson.Description, &lesson.Image, &lesson.Video, &lesson.Time, &lesson.CourseID)

	if err != nil {
		return lesson, err
	}

	return lesson, nil
}

func (h *LessonRepository) DeleteLessonByID(c context.Context, id uint) error {
	userQuery := `DELETE FROM lessons WHERE id=$1;`
	_, err := h.db.Exec(c, userQuery, id)
	if err != nil {
		return err
	}
	return nil
}

func (h *LessonRepository) GetCourseByID(c context.Context, id uint) (models.Course, error) {
	course := models.Course{}

	query := `SELECT id, name, description, image, points FROM courses WHERE id=$1;`
	row := h.db.QueryRow(c, query, id)
	err := row.Scan(&course.ID, &course.Name, &course.Description, &course.Image, &course.Points)

	if err != nil {
		return course, err
	}

	return course, nil
}

func (h *LessonRepository) CreateLessonProduct(c context.Context, lessonID uint, productID int) (int, error) {
	var id int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `INSERT INTO lesson_products(
		lesson_id, product_id, created_at)
		VALUES ($1, $2, $3) returning id;`
	err := h.db.QueryRow(c, userQuery, lessonID, productID, currentTime).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *LessonRepository) GetLessonProductsByID(c context.Context, lessonID uint) ([]int, error) {
	var IDs []int
	query := `SELECT product_id FROM lesson_products WHERE lesson_id = $1`
	rows, err := h.db.Query(c, query, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		IDs = append(IDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return IDs, nil
}

func (h *LessonRepository) DeleteLessonProductsByID(c context.Context, id uint) error {
	query := `DELETE FROM lesson_products WHERE lesson_id=$1;`
	_, err := h.db.Exec(c, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (h *LessonRepository) IsLessonNumberExist(c context.Context, courseID uint, lessonNumber int) (int, error) {
	var id int

	query := `SELECT id FROM lessons WHERE course_id=$1 AND lesson_number=$2;`
	row := h.db.QueryRow(c, query, courseID, lessonNumber)
	err := row.Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
}
