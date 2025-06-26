package dto

type PantryIngredientResponse struct {
	IngredientID       uint   `json:"ingredient_id"`
	IngredientName     string `json:"ingredient_name"`
}

type PantryIngredientRequest struct {
	IngredientID uint `json:"ingredient_id"`
}
