package dto

import (
	"zenful_shopping_backend/internal/store"
)

type IngredientRequest struct {
	Name        string `json:"name"`
	Measurement string `json:"measurement_unit"`
	FoodGroup   string `json:"food_group"`
}

type CreateIngredientRequest struct {
	IngredientID uint                 `json:"ingredient_id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Measurement store.MeasurementUnit `json:"measurement_unit"`
	FoodGroup   store.FoodGroup       `json:"food_group"`
}
