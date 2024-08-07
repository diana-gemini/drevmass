package models

import "time"

type Point struct {
	PointID        uint      `json:"pointID"`
	Description    string    `json:"description"`
	Point          int       `json:"point"`
	CreatedAt      time.Time `json:"createdAt"`
	ExpiredAt      time.Time `json:"expiredAt"`
	UserID         uint      `json:"userID"`
	IsAvailable    bool      `json:"isAvailable"`
	AvailablePoint int       `json:"availablePoint"`
	IsUsed         bool      `json:"isUsed"`
	PromocodeId    uint      `json:"promocode_id"`
}
