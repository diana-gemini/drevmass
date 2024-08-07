package controller

import (
	"github.com/diana-gemini/drevmass/internal/controller/admin/course"
	"github.com/diana-gemini/drevmass/internal/controller/admin/information"
	"github.com/diana-gemini/drevmass/internal/controller/auth"
	"github.com/diana-gemini/drevmass/internal/controller/middleware"
	"github.com/diana-gemini/drevmass/internal/controller/user"
	"github.com/diana-gemini/drevmass/internal/repository"
	"github.com/diana-gemini/drevmass/pkg"
	"github.com/gin-gonic/gin"
)

func GetRoute(app pkg.Application, r *gin.Engine) {
	db := app.DB

	loginController := &auth.AuthController{
		UserRepository: repository.NewUserRepository(db),
	}

	userController := &user.UserController{
		UserRepository: repository.NewUserRepository(db),
	}

	productController := &ProductController{
		ProductRepository: repository.NewProductRepository(db),
	}
	/*
		cartController := &CartController{
			CartRepository: repository.NewCartRepository(db),
		}*/

	companyController := &information.CompanyController{
		CompanyRepository: repository.NewCompanyRepository(db),
	}

	bonusController := &information.BonusController{
		BonusRepository: repository.NewBonusRepository(db),
	}

	appController := &information.AppController{
		AppRepository: repository.NewAppRepository(db),
	}

	contactController := &information.ContactController{
		ContactRepository: repository.NewContactRepository(db),
	}

	courseController := &course.CourseController{
		CourseRepository: repository.NewCourseRepository(db),
	}

	lessonController := &course.LessonController{
		LessonRepository: repository.NewLessonRepository(db),
	}

	r.POST("/signup", loginController.Signup)
	r.POST("/login", loginController.Login)
	r.POST("/passwordrecover", loginController.PasswordRecover)
	r.POST("/passwordreset", loginController.PasswordReset)

	profile := r.Group("/profile", middleware.RequireAuth)
	{
		profile.GET("/", userController.GetProfile)
		profile.POST("/logout", userController.Logout)
		profile.POST("/changepassword", userController.ChangePassword)

		points := profile.Group("/points")
		{
			points.GET("/", userController.GetPoints)
			points.GET("/about", userController.AboutBonus)
		}

		promocode := profile.Group("/promocode")
		{
			promocode.GET("/", userController.GetPromocode)
			promocode.POST("/copy", userController.CopyPromocode)
			promocode.POST("/activate", userController.ActivatePromocode)
		}

		userdata := profile.Group("/userdata")
		{
			userdata.GET("/", userController.GetUserData)
			userdata.PUT("/", userController.UpdateUserData)
			userdata.DELETE("/", userController.DeleteProfile)
		}

		notifications := profile.Group("/notifications")
		{
			notifications.GET("/", userController.GetNotifications)
			notifications.POST("/on", userController.OnNotifications)
			notifications.POST("/off", userController.OffNotifications)
		}

		contact := profile.Group("/contact")
		{
			contact.GET("/call", userController.GetCall)
			contact.POST("/support", userController.SendToSupport)
			contact.GET("/whatsapp", userController.GetWhatsapp)
		}

		information := profile.Group("/information")
		{
			information.GET("/company", userController.GetCompanyInformation)
			information.GET("/app", userController.GetAppInformation)
		}
	}

	courses := r.Group("/course", middleware.RequireAuth)
	{
		courses.GET("/", userController.GetAllCourses)
		courses.GET("/:id", userController.GetCourse)
		courses.POST("/:id/startcourse", userController.StartCourse)

		lesson := courses.Group("/:id/lesson/:lessonID")
		{
			lesson.GET("/", userController.GetLesson)
			lesson.POST("/video", userController.ShowVideo)
			lesson.POST("/bookmark", userController.SetBookmark)
		}

		bookmarks := courses.Group("/bookmarks")
		{
			bookmarks.GET("/", userController.GetUserBookmarks)
		}
	}

	products := r.Group("/products", middleware.RequireAuth)
	{
		products.GET("/", productController.getAlProdacts)
		products.GET("/:id", productController.getProdactById)
	}

	cart := r.Group("/cart", middleware.RequireAuth)
	{
		cart.POST("/:id", userController.AddToCart)
		cart.PUT("/minus/:id", userController.Minus)
		cart.PUT("/plus/:id", userController.Plus)
		cart.DELETE("/", userController.DeleteAllFromCart)
		cart.GET("/", userController.GetAllFromCart)

	}

	admin := r.Group("/admin", middleware.RequireAuth, middleware.IsAdmin())
	{
		information := admin.Group("/information")
		{
			company := information.Group("/company")
			{
				company.POST("/", companyController.CreateCompanyInformation)
				company.GET("/", companyController.GetCompanyInformation)
				company.PUT("/", companyController.UpdateCompanyInformation)
				company.DELETE("/", companyController.DeleteCompanyInformation)
			}

			app := information.Group("/app")
			{
				app.POST("/", appController.CreateAppInformation)
				app.GET("/", appController.GetAppInformation)
				app.PUT("/", appController.UpdateAppInformation)
				app.DELETE("/", appController.DeleteAppInformation)
			}

			bonus := information.Group("/bonus")
			{
				bonus.POST("/", bonusController.CreateBonusInformation)
				bonus.GET("/", bonusController.GetBonusInformation)
				bonus.PUT("/", bonusController.UpdateBonusInformation)
				bonus.DELETE("/", bonusController.DeleteBonusInformation)
			}

			contact := information.Group("/contact")
			{
				contact.POST("/", contactController.CreateContactInformation)
				contact.GET("/", contactController.GetContactInformation)
				contact.PUT("/", contactController.UpdateContactInformation)
				contact.DELETE("/", contactController.DeleteContactInformation)
			}
		}

		course := admin.Group("/course")
		{
			course.POST("/", courseController.CreateCourse)
			course.GET("/:id", courseController.GetCourse)
			course.PUT("/:id", courseController.UpdateCourse)
			course.DELETE("/:id", courseController.DeleteCourse)

			lesson := course.Group(":id/lesson")
			{
				lesson.POST("/", lessonController.CreateLesson)
				lesson.GET("/:lessonID", lessonController.GetLesson)
				lesson.PUT("/:lessonID", lessonController.UpdateLesson)
				lesson.DELETE("/:lessonID", lessonController.DeleteLesson)

				products := lesson.Group(":lessonID/products")
				{
					products.POST("/", lessonController.CreateLessonProducts)
					products.GET("/", lessonController.GetLessonProducts)
					products.PUT("/", lessonController.UpdateLessonProducts)
					products.DELETE("/", lessonController.DeleteLessonProducts)
				}
			}
		}

		products := admin.Group("/products")
		{
			products.POST("/", productController.createProduct)
			products.DELETE("/:id", productController.deleteProduct)
			products.PUT("/:id", productController.updateProduct)
		}

		orders := admin.Group("/order")
		{
			orders.GET("/", userController.GetAllOrders)
		}

	}

	order := r.Group("/order", middleware.RequireAuth)
	{
		order.POST("/", userController.CreateOrder)
	}

	crm := r.Group("/crm", middleware.RequireAuth)
	{
		crm.POST("/", user.SendDeal)
	}

}
