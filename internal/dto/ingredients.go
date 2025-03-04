package dto

type IngredientRequest struct {
	Name        string `json:"name"`
	Measurement string `json:"measurement_unit"`
	FoodGroup   string `json:"food_group"`
}

type CreateIngredientRequest struct {
	Ingredient IngredientRequest `json:"ingredient" form:"ingredient"`
}

type IngredientResponse struct {
	IngredientID uint   `json:"ingredient_id"`
	Name         string `json:"name"`
	Measurements []struct {
		MeasurementID uint   `json:"measurement_id"`
		Name          string `json:"name"`
	} `json:"measurements"`
	FoodGroup struct {
		FoodGroupID uint   `json:"food_group_id"`
		Name        string `json:"name"`
	} `json:"food_group"`
}

