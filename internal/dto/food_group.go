package dto

type CreateFoodGroup struct {
	Name        string `json:"name"`
}

type CreateFoodGroupRequest struct {
	Name           string `json:"name"`
	ShoppingListID uint   `json:"shopping_list_id"`
}
