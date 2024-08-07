package repository

import (
	"context"
	"time"

	"github.com/diana-gemini/drevmass/internal/models"
)

func (h *UserRepository) GetPointsByUserID(c context.Context, userID uint) ([]models.Point, error) {
	query := `SELECT description, point, created_at FROM points WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := h.db.Query(c, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	points := []models.Point{}

	for rows.Next() {
		var point models.Point
		if err := rows.Scan(&point.Description, &point.Point, &point.CreatedAt); err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return points, nil
}

func (h *UserRepository) GetAvailablePointsByUserID(c context.Context, userID uint) ([]models.Point, error) {
	query := `SELECT id, description, point, available_point, promocode_id FROM points WHERE user_id = $1 ORDER BY created_at`
	rows, err := h.db.Query(c, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	points := []models.Point{}

	for rows.Next() {
		var point models.Point
		if err := rows.Scan(&point.PointID, &point.Description, &point.Point, &point.AvailablePoint, &point.PromocodeId); err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return points, nil
}

func (h *UserRepository) GetSumPointsByUserID(c context.Context, userID uint) (int, error) {
	var sum int

	query := `SELECT COALESCE(SUM(available_point), 0) FROM points where user_id = $1`
	row := h.db.QueryRow(c, query, userID)
	err := row.Scan(&sum)

	if err != nil {
		return sum, err
	}

	return sum, nil
}

func (h *UserRepository) CheckExpiredUserPoints(c context.Context, userID uint) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	query := `SELECT id, point, expired_at, is_available, available_point FROM points WHERE user_id = $1 AND is_available = $2 AND expired_at < $3`
	rows, err := h.db.Query(c, query, userID, true, currentTime)
	if err != nil {
		return err
	}

	defer rows.Close()

	expiredPoints := []models.Point{}

	for rows.Next() {
		var point models.Point
		if err := rows.Scan(&point.PointID, &point.Point, &point.ExpiredAt, &point.IsAvailable, &point.AvailablePoint); err != nil {
			return err
		}
		expiredPoints = append(expiredPoints, point)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for _, point := range expiredPoints {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		userQuery := `UPDATE points SET is_available=$1, available_point=$2, updated_at=$3 WHERE id=$4;`
		_, err := h.db.Exec(c, userQuery, false, 0, currentTime, point.PointID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *UserRepository) GetExpiredPointByUserID(c context.Context, userID uint) (models.Point, error) {
	point := models.Point{}

	query := `SELECT id, expired_at, is_available, available_point FROM points WHERE user_id = $1 AND is_available = $2 ORDER BY expired_at DESC`
	row := h.db.QueryRow(c, query, userID, "t")
	err := row.Scan(&point.PointID, &point.ExpiredAt, &point.IsAvailable, &point.AvailablePoint)

	if err != nil {
		return point, err
	}

	return point, nil
}

func (h *UserRepository) UpdateAvailablePointsZero(c context.Context, pointId uint, userId uint) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `UPDATE points SET is_available=$1, available_point=$2, updated_at=$3,  is_used = true WHERE id=$4 and user_id=$5;`
	_, err := h.db.Exec(c, userQuery, false, 0, currentTime, pointId, userId)
	if err != nil {
		return err
	}
	return nil

}

func (h *UserRepository) UpdateAvailablePoints(c context.Context, pointId uint, userId uint, bonus float64) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `UPDATE points SET is_available=$1, available_point= available_point-$2, updated_at=$3, is_used = true WHERE id=$4 and user_id=$5;`
	_, err := h.db.Exec(c, userQuery, true, bonus, currentTime, pointId, userId)
	if err != nil {
		return err
	}
	return nil

}
