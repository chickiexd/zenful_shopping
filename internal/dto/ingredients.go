package dto

type CreateIngredientRequest struct {
	IngredientID uint                 `json:"ingredient_id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Measurement CreateMeasurementUnit `json:"measurement_unit"`
	FoodGroup   CreateFoodGroup       `json:"food_group"`
}
