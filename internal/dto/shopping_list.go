package dto

type ShoppingListResponse struct {
	ShoppingListID uint                       `json:"shopping_list_id"`
	Name           string                     `json:"name"`
	Items          []ShoppingListItemResponse `json:"ingredients"`
}

type ShoppingListItemResponse struct {
	ShoppingListItemID uint    `json:"shopping_list_item_id"`
	IngredientID       uint    `json:"ingredient_id"`
	Quantity           float64 `json:"quantity"`
	MeasurementUnitID  uint    `json:"measurement_unit_id"`
}
