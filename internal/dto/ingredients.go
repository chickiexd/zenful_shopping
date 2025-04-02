package dto

type CreateIngredientRequest struct {
	IngredientID uint                  `json:"ingredient_id"`
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	Measurement  CreateMeasurementUnit `json:"measurement_unit"`
	FoodGroup    CreateFoodGroup       `json:"food_group"`
}

type IngredientResponse struct {
	IngredientID uint              `json:"ingredient_id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Measurements []MeasurementUnitResponse `json:"measurements"`
	FoodGroups   []FoodGroupResponse       `json:"food_groups"`
}

type FoodGroupResponse struct {
	FoodGroupID uint   `json:"food_group_id"`
	Name        string `json:"name"`
}	

