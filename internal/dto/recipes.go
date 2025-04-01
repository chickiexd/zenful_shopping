package dto

import "mime/multipart"

type RecipeResponse struct {
	RecipeID     uint                  `json:"recipe_id"`
	Title        string                `json:"title"`
	Description  string                `json:"description"`
	Public       bool                  `json:"public"`
	CookTime     int                   `json:"cook_time"`
	Servings     int                   `json:"servings"`
	ImagePath    string                `json:"image_path"`
	MealType     uint                  `json:"meal_type"`
	Ingredients  []IngredientResponse  `json:"ingredients"`
	Instructions []InstructionResponse `json:"instructions"`
	CreatedAt    string                `json:"created_at"`
	UpdatedAt    string                `json:"updated_at"`
}

type IngredientResponse struct {
	IngredientID      uint    `json:"ingredient_id"`
	Quantity          float64 `json:"quantity"`
	MeasurementUnitID uint    `json:"measurement_unit_id"`
}

type InstructionResponse struct {
	InstructionID uint   `json:"instruction_id"`
	StepNumber    int    `json:"step_number"`
	Description   string `json:"description"`
}

type MeasurementUnitResponse struct {
	MeasurementUnitID uint   `json:"measurement_unit_id"`
	Name              string `json:"name"`
}

type CreateRecipeRequest struct {
	Recipe       recipeRequest         `json:"recipe" form:"recipe"`
	Ingredients  []ingredientRequest   `json:"ingredients"`
	Instructions []instructionRequest  `json:"instructions"`
	Image        multipart.File        `form:"-"`
	ImageHeader  *multipart.FileHeader `form:"-"`
}

type ingredientRequest struct {
	IngredientID      uint    `json:"ingredient_id"`
	Quantity          float64 `json:"quantity"`
	MeasurementUnitID uint    `json:"measurement_unit_id"`
}

type instructionRequest struct {
	Content   string `json:"content"`
	Numbering int    `json:"numbering"`
}

type recipeRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CookTime    int    `json:"cook_time"`
	Servings    int    `json:"servings"`
	Public      bool   `json:"public"`
	MealType    int    `json:"meal_type"`
}
