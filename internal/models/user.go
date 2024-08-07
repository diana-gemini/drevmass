package models

import (
	"context"
	"time"
)

type User struct {
	ID            uint      `json:"id"`
	Email         string    `json:"email" form:"email"`
	Username      string    `json:"username" form:"username"`
	Mobile        string    `json:"mobile" form:"mobile"`
	Password      string    `json:"password" form:"password"`
	BirthDate     string    `json:"birthdate" form:"birthdate"`
	Gender        string    `json:"gender" form:"gender"`
	Height        string    `json:"height" form:"height"`
	Weight        string    `json:"weight" form:"weight"`
	Activity      string    `json:"activity" form:"activity"`
	Notifications bool      `json:"notifications" form:"notifications"`
	RoleID        uint      `json:"roleID" `
	CreatedAt     string    `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type SignupUser struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Mobile   string `json:"mobile" form:"mobile"`
	Password string `json:"password" form:"password"`
}

type LoginUser struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type UserEmail struct {
	Email string `json:"email" form:"email"`
}

type UserPassword struct {
	Password string `json:"password" form:"password"`
}

type UserChangePassword struct {
	Current string `json:"current" form:"current"`
	New     string `json:"new" form:"new"`
}

type AuthUser struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	RoleID uint   `json:"roleID"`
}

type Promocode struct {
	ID        uint   `json:"id"`
	Promocode string `json:"promocode"`
	UserID    uint   `json:"userID"`
	Count     int    `json:"count"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type UserRepository interface {
	CreateUser(c context.Context, user SignupUser) (int, error)
	GetUserByEmail(c context.Context, email string) (User, error)
	GetUserByID(c context.Context, userID uint) (User, error)
	UpdateUserPassword(c context.Context, user User, password string) error
	GetAvailablePointsByUserID(c context.Context, userID uint) ([]Point, error)
	GetPointsByUserID(c context.Context, userID uint) ([]Point, error)
	GetSumPointsByUserID(c context.Context, userID uint) (int, error)
	UpdateUser(c context.Context, user User, updateUser User) error
	DeleteUser(c context.Context, user User) error
	UpdateUserNotifications(c context.Context, user User, notifications string) error
	GetCompanyInformation(c context.Context) (Company, error)
	CreatePoint(c context.Context, userID uint, point Point) (int, error)
	CreatePointForSignup(c context.Context, userID uint, point Point) (int, error)
	CreatePointForPromo(c context.Context, userID uint, point Point) (int, error)
	CheckExpiredUserPoints(c context.Context, userID uint) error
	GetExpiredPointByUserID(c context.Context, userID uint) (Point, error)
	GetBonusInformation(c context.Context) (Bonus, error)
	GetAppInformation(c context.Context) (App, error)
	GetUniquePromocode(c context.Context, promocode string) (int, error)
	GetPromocodeByPromocode(c context.Context, promo string) (Promocode, error)
	GetPromocodeById(c context.Context, id uint) (Promocode, error)
	CreatePromocode(c context.Context, userID uint, promocode string) (int, error)
	CreatePointUsedPromo(c context.Context, userID uint, point Point) error
	GetPromocodeByUserID(c context.Context, userID uint) (Promocode, error)
	UpdatePromocodeByUserID(c context.Context, userID uint, updatePromocode Promocode) error
	GetContactInformation(c context.Context) (Contact, error)
	GetAllCourses(c context.Context) ([]Course, error)
	GetCourseByID(c context.Context, id, userID uint) (Course, error)
	GetLessonsByCourseID(c context.Context, id, userID uint) ([]Lesson, error)
	GetLessonByID(c context.Context, courseID, lessonID uint) (Lesson, error)
	GetUserBookmarkByIDs(c context.Context, userID, courseID, lessonID uint) (int, error)
	CreateUserBookmark(c context.Context, userID, courseID, lessonID uint) (uint, error)
	DeleteUserBookmark(c context.Context, userID, courseID, lessonID uint) error
	GetLessonProducts(c context.Context, id uint) ([]Product, error)
	GetUserBookmarks(c context.Context, userID uint) ([]Bookmark, error)
	CreateUserStartCourse(c context.Context, userID, courseID uint) (int, error)
	GetUserStartCourse(c context.Context, userID, courseID uint) (int, error)
	CreateOfLessonsWatch(c context.Context, userID, courseID, lessonID uint) (int, error)
	GetCountOfLessonsWatch(c context.Context, id, userID uint) (Course, error)
	UpdateUserFinishCourse(c context.Context, userID, courseID uint) error
	Add(c context.Context, userId uint, productId int) (int, error)
	Minus(c context.Context, userId uint, productId int) error
	Plus(c context.Context, userId uint, productId int) error
	DeleteFromCart(c context.Context, userId uint) error
	GetAllFromCart(c context.Context, userId uint) ([]GetProductsFromCart, error)
	GetTotalAmout(c context.Context, userId uint) (float64, error)
	GetUserInfoByID(c context.Context, userID uint) (User, error)
	CreateOrder(c context.Context, userId int, totalAmount float64, bonus float64, overall float64, products []GetProductsFromCart) (int, error)
	GetAllOrdersForUser(c context.Context, userId int) ([]OrderForUser, error)
	GetAllOrders(c context.Context) ([]OrderForAdmin, error)
	GetOrderById(c context.Context, id int) ([]OrderForAdmin, error)

	UpdateAvailablePointsZero(c context.Context, pointId uint, userId uint) error
	UpdateAvailablePoints(c context.Context, pointId uint, userId uint, bonus float64) error
}
