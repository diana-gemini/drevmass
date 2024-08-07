package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/diana-gemini/drevmass/internal/models"
)

type AppRepository struct {
	db *pgxpool.Pool
}

func NewAppRepository(db *pgxpool.Pool) models.AppRepository {
	return &AppRepository{db: db}
}

func (h *AppRepository) CreateAppInformation(c context.Context, app models.App) (int, error) {
	var appID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `INSERT INTO info_app(
		name, image, app_name, description, version, release_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) returning id;`
	err := h.db.QueryRow(c, userQuery, app.Name, app.Image, app.AppName, app.Description, app.Version, app.ReleaseDate, currentTime).Scan(&appID)
	if err != nil {
		return 0, err
	}
	return appID, nil
}

func (h *AppRepository) UpdateAppInformation(c context.Context, app models.App, updateApp models.App) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `UPDATE info_app SET name=$1, image=$2, app_name=$3, description=$4, version=$5, release_date=$6, updated_at=$7 WHERE id=$8;`

	_, err := h.db.Exec(c, userQuery, updateApp.Name, updateApp.Image, updateApp.AppName, updateApp.Description, updateApp.Version, updateApp.ReleaseDate, currentTime, app.ID)
	if err != nil {
		return err
	}
	return nil
}

func (h *AppRepository) GetAppInformation(c context.Context) (models.App, error) {
	app := models.App{}

	query := `SELECT id, name, image, app_name, description, version, release_date FROM info_app `
	row := h.db.QueryRow(c, query)
	err := row.Scan(&app.ID, &app.Name, &app.Image, &app.AppName, &app.Description, &app.Version, &app.ReleaseDate)

	if err != nil {
		return app, err
	}

	return app, nil
}

func (h *AppRepository) DeleteAppInformation(c context.Context, app models.App) error {
	userQuery := `DELETE FROM info_app WHERE id=$1;`
	_, err := h.db.Exec(c, userQuery, app.ID)
	if err != nil {
		return err
	}
	return nil
}
