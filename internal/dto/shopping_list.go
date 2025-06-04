package dto

type ShoppingListResponse struct {
	ShoppingListID uint                       `json:"shopping_list_id"`
	Name           string                     `json:"name"`
	Items          []ShoppingListItemResponse `json:"ingredients"`
}

type ShoppingListItemResponse struct {
	ShoppingListItemID uint    `json:"shopping_list_item_id"`
	ShoppingListID     uint    `json:"shopping_list_id"`
	IngredientID       uint    `json:"ingredient_id"`
	Quantity           float64 `json:"quantity"`
	MeasurementUnitID  uint    `json:"measurement_unit_id"`
}

type ShoppingListItemID struct {
	ItemID uint `json:"item_id"`
}

type ShoppingListID struct {
	ShoppingListID uint `json:"shopping_list_id"`
}

type ShoppingListName struct {
	Name string `json:"shopping_list_name"`
}
