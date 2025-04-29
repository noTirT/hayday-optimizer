package models

import (
	"time"

	"github.com/google/uuid"
)

type Ingredient struct {
	ProductID   uuid.UUID
	ProductName string
	Amount      int
}

type HayDayGood struct {
	ID             uuid.UUID
	Name           string
	RequiredLevel  int
	MaxPrice       int
	ProductionTime time.Duration
	GainedXP       int
	Ingredients    []Ingredient
	Source         string
	RawIngredients string
}

type RawTableRow struct {
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Price    int    `json:"price"`
	TimeStr  string `json:"timeStr"`
	XP       int    `json:"xp"`
	RawNeeds string `json:"rawNeeds"`
	Source   string `json:"source"`
}

type HayDayGoodList []HayDayGood
