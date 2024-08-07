package main

import (
	"log"

	_ "github.com/diana-gemini/drevmass/docs"
	"github.com/diana-gemini/drevmass/internal/controller"
	"github.com/diana-gemini/drevmass/internal/db"
	"github.com/diana-gemini/drevmass/pkg"
	"github.com/diana-gemini/drevmass/pkg/envs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Drevmass API
// @version 1.0
// @description API Server for Drevmass

// @host 185.100.67.103:3000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func init() {
	envs.LoadEnvVariables()
	// 	db.ConnectDB()
}

func main() {
	app, err := pkg.App()
	if err != nil {
		log.Fatal(err)
	}
	defer app.CloseDBConnection()

	err = db.CreateTable(app.DB)
	if err != nil {
		log.Fatal(err)
	}

	err = db.CreateAdmin(app.DB)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://185.100.67.103:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Serve Swagger documentation
	url := ginSwagger.URL("http://185.100.67.103:3000/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	controller.GetRoute(app, r)

	r.Run("0.0.0.0:3000")

}
