package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/diana-gemini/drevmass/internal/models"
)

type BonusRepository struct {
	db *pgxpool.Pool
}

func NewBonusRepository(db *pgxpool.Pool) models.BonusRepository {
	return &BonusRepository{db: db}
}

func (h *BonusRepository) CreateBonusInformation(c context.Context, bonus models.Bonus) (int, error) {
	var bonusID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `INSERT INTO info_bonus(
		name, description, created_at)
		VALUES ($1, $2, $3) returning id;`
	err := h.db.QueryRow(c, userQuery, bonus.Name, bonus.Description, currentTime).Scan(&bonusID)
	if err != nil {
		return 0, err
	}
	return bonusID, nil
}

func (h *BonusRepository) UpdateBonusInformation(c context.Context, bonus models.Bonus, updateBonus models.Bonus) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `UPDATE info_bonus SET name=$1, description=$2, updated_at=$3 WHERE id=$4;`

	_, err := h.db.Exec(c, userQuery, updateBonus.Name, updateBonus.Description, currentTime, bonus.ID)
	if err != nil {
		return err
	}
	return nil
}

func (h *BonusRepository) GetBonusInformation(c context.Context) (models.Bonus, error) {
	bonus := models.Bonus{}

	query := `SELECT id, name, description FROM info_bonus `
	row := h.db.QueryRow(c, query)
	err := row.Scan(&bonus.ID, &bonus.Name, &bonus.Description)

	if err != nil {
		return bonus, err
	}

	return bonus, nil
}

func (h *BonusRepository) DeleteBonusInformation(c context.Context, bonus models.Bonus) error {
	userQuery := `DELETE FROM info_bonus WHERE id=$1;`
	_, err := h.db.Exec(c, userQuery, bonus.ID)
	if err != nil {
		return err
	}
	return nil
}
