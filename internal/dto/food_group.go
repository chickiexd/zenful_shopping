package dto

type CreateFoodGroup struct {
	FoodGroupID uint   `json:"food_group_id"`
	Name string `json:"name"`
	Description string `json:"description"`
}
