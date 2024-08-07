package repository

import (
	"context"
	"time"

	"github.com/diana-gemini/drevmass/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

// CreatePointForSignup implements models.UserRepository.
func (*UserRepository) CreatePointForSignup(c context.Context, userID uint, point models.Point) (int, error) {
	panic("unimplemented")
}

func NewUserRepository(db *pgxpool.Pool) models.UserRepository {
	return &UserRepository{db: db}
}

func (h *UserRepository) CreateUser(c context.Context, user models.SignupUser) (int, error) {
	var userID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `INSERT INTO users(
		email, username, mobile, password, role_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) returning id;`
	err := h.db.QueryRow(c, userQuery, user.Email, user.Username, user.Mobile, user.Password, 2, currentTime).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (h *UserRepository) GetUserByEmail(c context.Context, email string) (models.User, error) {
	user := models.User{}

	query := `SELECT id, email, username, mobile, password, role_id, created_at FROM users where email=$1`
	row := h.db.QueryRow(c, query, email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Mobile, &user.Password, &user.RoleID, &user.CreatedAt)

	if err != nil {
		return user, err
	}
	return user, nil
}

func (h *UserRepository) UpdateUserPassword(c context.Context, user models.User, password string) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `UPDATE users SET password=$1, updated_at=$2 WHERE id=$3;`
	_, err := h.db.Exec(c, userQuery, password, currentTime, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (h *UserRepository) GetUserByID(c context.Context, userID uint) (models.User, error) {
	user := models.User{}

	query := `SELECT id, email, username, mobile, notifications, role_id, 
			birthdate, gender, height, weight, activity FROM users WHERE id=$1`
	row := h.db.QueryRow(c, query, userID)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Mobile, &user.Notifications, &user.RoleID,
		&user.BirthDate, &user.Gender, &user.Height, &user.Weight, &user.Activity)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (h *UserRepository) GetUserInfoByID(c context.Context, userID uint) (models.User, error) {
	user := models.User{}

	query := `SELECT id, email, username, mobile FROM users WHERE id=$1`
	row := h.db.QueryRow(c, query, userID)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Mobile)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (h *UserRepository) UpdateUser(c context.Context, user models.User, updateUser models.User) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `UPDATE users SET username=$1, mobile=$2, birthdate=$3, gender=$4, height=$5, weight=$6, activity=$7, updated_at=$8 WHERE id=$9;`
	_, err := h.db.Exec(c, userQuery, updateUser.Username, updateUser.Mobile, updateUser.BirthDate, updateUser.Gender, updateUser.Height, updateUser.Weight, updateUser.Activity, currentTime, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (h *UserRepository) DeleteUser(c context.Context, user models.User) error {
	userQuery := `DELETE FROM users WHERE id=$1;`
	_, err := h.db.Exec(c, userQuery, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (h *UserRepository) UpdateUserNotifications(c context.Context, user models.User, notifications string) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `UPDATE users SET notifications=$1, updated_at=$2 WHERE id=$3;`
	_, err := h.db.Exec(c, userQuery, notifications, currentTime, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (h *UserRepository) GetCompanyInformation(c context.Context) (models.Company, error) {
	company := models.Company{}

	query := `SELECT id, name, description, image FROM company`
	row := h.db.QueryRow(c, query)
	err := row.Scan(&company.ID, &company.Name, &company.Description, &company.Image)

	if err != nil {
		return company, err
	}

	return company, nil
}

func (h *UserRepository) CreatePoint(c context.Context, userID uint, point models.Point) (int, error) {
	var pointID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	expiredData := time.Now().AddDate(0, 0, 14).Format("2006-01-02 15:04:05")
	userQuery := `INSERT INTO points(
		description, point, created_at, expired_at, user_id, is_available, available_point)
		VALUES ($1, $2, $3, $4, $5, $6, $7) returning id;`
	err := h.db.QueryRow(c, userQuery, point.Description, point.Point, currentTime, expiredData, userID, point.IsAvailable, point.Point).Scan(&pointID)
	if err != nil {
		return 0, err
	}
	return pointID, nil
}

func (h *UserRepository) CreatePointForPromo(c context.Context, userID uint, point models.Point) (int, error) {
	var pointID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	expiredData := time.Now().AddDate(0, 0, 14).Format("2006-01-02 15:04:05")
	userQuery := `INSERT INTO points(
		description, point, created_at, expired_at, user_id, is_available, available_point,promocode_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7,$8) returning id;`
	err := h.db.QueryRow(c, userQuery, point.Description, point.Point, currentTime, expiredData, userID, point.IsAvailable, point.Point, point.PromocodeId).Scan(&pointID)
	if err != nil {
		return 0, err
	}
	return pointID, nil
}

func (h *UserRepository) CreatePointUsedPromo(c context.Context, userID uint, point models.Point) error {
	var pointID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	expiredData := time.Now().AddDate(0, 0, 14).Format("2006-01-02 15:04:05")
	userQuery := `INSERT INTO points(
		description, point, created_at, expired_at, user_id, is_available, available_point)
		VALUES ($1, $2, $3, $4, $5, $6, $7) returning id;`
	err := h.db.QueryRow(c, userQuery, point.Description, point.Point, currentTime, expiredData, userID, point.IsAvailable, point.Point).Scan(&pointID)
	if err != nil {
		return err
	}
	return nil
}

func (h *UserRepository) GetBonusInformation(c context.Context) (models.Bonus, error) {
	bonus := models.Bonus{}

	query := `SELECT id, name, description FROM bonus`
	row := h.db.QueryRow(c, query)
	err := row.Scan(&bonus.ID, &bonus.Name, &bonus.Description)

	if err != nil {
		return bonus, err
	}

	return bonus, nil
}

func (h *UserRepository) GetAppInformation(c context.Context) (models.App, error) {
	app := models.App{}

	query := `SELECT id, name, image, app_name, description, version, release_date FROM app `
	row := h.db.QueryRow(c, query)
	err := row.Scan(&app.ID, &app.Name, &app.Image, &app.AppName, &app.Description, &app.Version, &app.ReleaseDate)

	if err != nil {
		return app, err
	}

	return app, nil
}

func (h *UserRepository) GetUniquePromocode(c context.Context, promocode string) (int, error) {
	var id int

	query := `SELECT id FROM promocodes WHERE promocode=$1`
	row := h.db.QueryRow(c, query, promocode)
	err := row.Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
}

func (h *UserRepository) CreatePromocode(c context.Context, userID uint, promocode string) (int, error) {
	var id int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `INSERT INTO promocodes(
		promocode, user_id, count, created_at)
		VALUES ($1, $2, $3, $4) returning id;`
	err := h.db.QueryRow(c, userQuery, promocode, userID, 0, currentTime).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *UserRepository) GetPromocodeByUserID(c context.Context, userID uint) (models.Promocode, error) {
	promocode := models.Promocode{}

	query := `SELECT id, promocode, user_id, count FROM promocodes WHERE user_id=$1`
	row := h.db.QueryRow(c, query, userID)
	err := row.Scan(&promocode.ID, &promocode.Promocode, &promocode.UserID, &promocode.Count)

	if err != nil {
		return promocode, err
	}

	return promocode, nil
}

func (h *UserRepository) GetPromocodeById(c context.Context, id uint) (models.Promocode, error) {
	promocode := models.Promocode{}

	query := `SELECT id, promocode, user_id, count FROM promocodes WHERE id=$1`
	row := h.db.QueryRow(c, query, id)
	err := row.Scan(&promocode.ID, &promocode.Promocode, &promocode.UserID, &promocode.Count)

	if err != nil {
		return promocode, err
	}

	return promocode, nil
}

func (h *UserRepository) GetPromocodeByPromocode(c context.Context, promo string) (models.Promocode, error) {
	promocode := models.Promocode{}

	query := `SELECT id, promocode, user_id, count FROM promocodes WHERE promocode=$1`
	row := h.db.QueryRow(c, query, promo)
	err := row.Scan(&promocode.ID, &promocode.Promocode, &promocode.UserID, &promocode.Count)
	//logrus.Print(query)
	if err != nil {
		return promocode, err
	}

	return promocode, nil
}

func (h *UserRepository) UpdatePromocodeByUserID(c context.Context, userID uint, updatePromocode models.Promocode) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `UPDATE promocodes SET count=$1, updated_at=$2 WHERE user_id=$3;`
	_, err := h.db.Exec(c, userQuery, updatePromocode.Count, currentTime, userID)
	if err != nil {
		return err
	}
	return nil
}

func (h *UserRepository) GetContactInformation(c context.Context) (models.Contact, error) {
	contact := models.Contact{}

	query := `SELECT id, name FROM info_contact`
	row := h.db.QueryRow(c, query)
	err := row.Scan(&contact.ID, &contact.Name)

	if err != nil {
		return contact, err
	}

	return contact, nil
}

func (h *UserRepository) GetAllCourses(c context.Context) ([]models.Course, error) {
	query := `SELECT c.id, c.name, c.image, c.points, 
					COALESCE(COUNT(l.id),0) AS lesson_count, 
					COALESCE(ROUND(EXTRACT(epoch FROM SUM(l.time))/60),0) AS course_time
			FROM courses c
			LEFT JOIN lessons l ON c.id = l.course_id
			GROUP BY c.id, c.name, c.image, c.points`

	rows, err := h.db.Query(c, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	courses := []models.Course{}

	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Image, &course.Points, &course.LessonCount, &course.Time); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

func (h *UserRepository) GetCourseByID(c context.Context, id, userID uint) (models.Course, error) {
	course := models.Course{}

	query := `WITH watch_count AS (
				SELECT course_id, COUNT(is_watched) AS count FROM user_watched_lessons
				WHERE is_watched = $1 AND user_id = $2
				GROUP BY course_id)
	
			SELECT c.id, c.name, c.description, c.image, c.points, 
                COALESCE(COUNT(l.id), 0) AS lesson_count, 
                COALESCE(ROUND(EXTRACT(epoch FROM SUM(l.time)) / 60), 0) AS course_time,
				COALESCE(w.count,0) AS watch_count
			FROM courses c
			LEFT JOIN lessons l ON c.id = l.course_id
			LEFT JOIN watch_count w ON c.id = w.course_id
			WHERE c.id = $3
			GROUP BY c.id, c.name, c.description, c.image, c.points, w.count;`

	row := h.db.QueryRow(c, query, "t", userID, id)
	err := row.Scan(&course.ID, &course.Name, &course.Description, &course.Image, &course.Points,
		&course.LessonCount, &course.Time, &course.WatchCount)

	if err != nil {
		return course, err
	}

	return course, nil
}

func (h *UserRepository) GetLessonsByCourseID(c context.Context, id, userID uint) ([]models.Lesson, error) {
	query := `WITH bookmark AS (
				SELECT id, lesson_id FROM bookmarks 
				WHERE user_id = $1),
			watch AS (
				SELECT lesson_id, is_watched FROM user_watched_lessons
				WHERE user_id = $1)

			SELECT l.id, l.lesson_number, l.name, l.description, l.image, l.video, 
					COALESCE(ROUND(EXTRACT(epoch FROM l.time)/60),0) AS lesson_time,
					COALESCE(b.id,0) AS bookmark,
					COALESCE(w.is_watched, false) AS is_watched
			FROM lessons l
			LEFT JOIN bookmark b ON l.id = b.lesson_id
			LEFT JOIN watch w ON l.id = w.lesson_id
			WHERE l.course_id = $2 
			ORDER BY l.lesson_number;`

	rows, err := h.db.Query(c, query, userID, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	lessons := []models.Lesson{}
	for rows.Next() {
		var lesson models.Lesson
		if err := rows.Scan(&lesson.ID, &lesson.LessonNumber, &lesson.Name, &lesson.Description, &lesson.Image,
			&lesson.Video, &lesson.Time, &lesson.BookmarkID, &lesson.IsWatched); err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	if err := rows.Err(); err != nil {
		return nil, err

	}
	return lessons, nil
}

func (h *UserRepository) GetLessonByID(c context.Context, courseID, lessonID uint) (models.Lesson, error) {
	lesson := models.Lesson{}
	query := `SELECT id, lesson_number, name, description, image, video, 
				COALESCE(ROUND(EXTRACT(epoch FROM time)/60),0) AS lesson_time,
				course_id
			FROM lessons 
			WHERE id=$1 AND course_id=$2`

	row := h.db.QueryRow(c, query, lessonID, courseID)
	err := row.Scan(&lesson.ID, &lesson.LessonNumber, &lesson.Name, &lesson.Description, &lesson.Image, &lesson.Video, &lesson.Time, &lesson.CourseID)
	if err != nil {
		return lesson, err
	}
	return lesson, nil
}

func (h *UserRepository) GetUserBookmarkByIDs(c context.Context, userID, courseID, lessonID uint) (int, error) {
	var id int

	query := `SELECT id FROM bookmarks WHERE user_id=$1 AND course_id=$2 AND lesson_id=$3`
	row := h.db.QueryRow(c, query, userID, courseID, lessonID)
	err := row.Scan(&id)
	if err != nil {

		return id, err
	}
	return id, nil
}

func (h *UserRepository) CreateUserBookmark(c context.Context, userID, courseID, lessonID uint) (uint, error) {
	var id uint
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO bookmarks (user_id, course_id, lesson_id, created_at) VALUES ($1, $2, $3, $4) returning id;`
	err := h.db.QueryRow(c, query, userID, courseID, lessonID, currentTime).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *UserRepository) DeleteUserBookmark(c context.Context, userID, courseID, lessonID uint) error {
	query := `DELETE FROM bookmarks WHERE user_id=$1 AND course_id=$2 AND lesson_id=$3`
	_, err := h.db.Exec(c, query, userID, courseID, lessonID)
	if err != nil {
		return err
	}
	return nil
}

func (h *UserRepository) GetUserBookmarks(c context.Context, userID uint) ([]models.Bookmark, error) {
	query := `WITH lesson_data AS (
				SELECT id, lesson_number, COALESCE(ROUND(EXTRACT(epoch FROM time)/60),0) AS time, name, image, video
				FROM lessons
			)
			SELECT b.course_id, c.name, 
				   array_agg(json_build_object(	'id', l.id, 
				   								'lessonNumber', l.lesson_number,
												'time', l.time,
												'name', l.name,
												'image', l.image,
												'video', l.video )
												ORDER BY l.lesson_number) AS lessons
			FROM bookmarks b 
			LEFT JOIN courses c ON b.course_id = c.id
			LEFT JOIN lesson_data l ON b.lesson_id = l.id 
			WHERE b.user_id = $1
			GROUP BY b.course_id, c.name;`

	rows, err := h.db.Query(c, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookmarks := []models.Bookmark{}
	for rows.Next() {
		var bookmark models.Bookmark
		if err := rows.Scan(&bookmark.CourseID, &bookmark.CourseName, &bookmark.Lessons); err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bookmark)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (h *UserRepository) GetLessonProducts(c context.Context, id uint) ([]models.Product, error) {
	query := `SELECT l.product_id, p.name, p.image, p.price
			FROM lesson_products l
			LEFT JOIN products p ON l.product_id=p.id
			WHERE l.lesson_id=$1;`
	rows, err := h.db.Query(c, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Image, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (h *UserRepository) CreateUserStartCourse(c context.Context, userID, courseID uint) (int, error) {
	var id int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO user_courses (is_started, is_finished, user_id, course_id, created_at) 
			VALUES ($1, $2, $3, $4, $5) returning id;`
	err := h.db.QueryRow(c, query, "t", "f", userID, courseID, currentTime).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *UserRepository) GetUserStartCourse(c context.Context, userID, courseID uint) (int, error) {
	var id int
	query := `SELECT id FROM user_courses WHERE user_id=$1 AND course_id=$2 AND is_started=$3;`
	row := h.db.QueryRow(c, query, userID, courseID, "t")
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *UserRepository) CreateOfLessonsWatch(c context.Context, userID, courseID, lessonID uint) (int, error) {
	var id int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO user_watched_lessons (user_id, course_id, lesson_id, is_watched, created_at) 
			VALUES ($1, $2, $3, $4, $5) returning id;`
	err := h.db.QueryRow(c, query, userID, courseID, lessonID, "t", currentTime).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *UserRepository) GetCountOfLessonsWatch(c context.Context, id, userID uint) (models.Course, error) {
	course := models.Course{}

	query := `WITH watch_count AS (
				SELECT course_id, COUNT(is_watched) AS count FROM user_watched_lessons
				WHERE is_watched = $1 AND user_id = $2
				GROUP BY course_id)
	
			SELECT l.course_id, COALESCE(COUNT(l.id), 0) AS lesson_count, 
				COALESCE(w.count,0) AS watch_count
			FROM lessons l 
			LEFT JOIN watch_count w ON l.course_id = w.course_id
			WHERE l.course_id = $3
			GROUP BY l.course_id, watch_count;`

	row := h.db.QueryRow(c, query, "t", userID, id)
	err := row.Scan(&course.ID, &course.LessonCount, &course.WatchCount)

	if err != nil {
		return course, err
	}

	return course, nil
}

func (h *UserRepository) UpdateUserFinishCourse(c context.Context, userID, courseID uint) error {
	query := `UPDATE user_courses SET is_finished=$1 WHERE user_id=$2 AND course_id=$3;`
	_, err := h.db.Exec(c, query, "t", userID, courseID)
	if err != nil {
		return err
	}
	return nil
}
