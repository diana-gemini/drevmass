package models

import "context"

type App struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" form:"name"`
	Image       string `json:"image" form:"image"`
	AppName     string `json:"appName" form:"appName"`
	Description string `json:"description" form:"description"`
	Version     string `json:"version" form:"version"`
	ReleaseDate string `json:"releaseDate" form:"releaseDate"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type AppRepository interface {
	CreateAppInformation(c context.Context, app App) (int, error)
	UpdateAppInformation(c context.Context, app App, updateApp App) error
	GetAppInformation(c context.Context) (App, error)
	DeleteAppInformation(c context.Context, app App) error
}
