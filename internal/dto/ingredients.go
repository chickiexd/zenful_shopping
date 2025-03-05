package dto

type IngredientRequest struct {
	Name        string `json:"name"`
	Measurement string `json:"measurement_unit"`
	FoodGroup   string `json:"food_group"`
}

type CreateIngredientRequest struct {
	Ingredient IngredientRequest `json:"ingredient" form:"ingredient"`
}
