package dto

type CreateIngredientRequest struct {
	IngredientID uint                    `json:"ingredient_id"`
	Name         string                  `json:"name"`
	Description  string                  `json:"description"`
	Measurements []CreateMeasurementUnit `json:"measurement_units"`
	FoodGroups   []string                `json:"food_groups"`
}

type IngredientResponse struct {
	IngredientID uint                      `json:"ingredient_id"`
	Name         string                    `json:"name"`
	Description  string                    `json:"description"`
	Measurements []MeasurementUnitResponse `json:"measurements"`
	FoodGroups   []FoodGroupResponse       `json:"food_groups"`
}

type FoodGroupResponse struct {
	FoodGroupID uint   `json:"food_group_id"`
	Name        string `json:"name"`
}

type AddIngredientToShoppingListRequest struct {
	IngredientID      uint    `json:"ingredient_id"`
	Quantity          float64 `json:"quantity"`
	MeasurementUnitID uint    `json:"measurement_unit_id"`
	ShoppingListID    uint    `json:"shopping_list_id"`
}
