package models

import (
	"context"
)

type Bonus struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type BonusRepository interface {
	CreateBonusInformation(c context.Context, bonus Bonus) (int, error)
	UpdateBonusInformation(c context.Context, bonus Bonus, updateBonus Bonus) error
	GetBonusInformation(c context.Context) (Bonus, error)
	DeleteBonusInformation(c context.Context, bonus Bonus) error
}