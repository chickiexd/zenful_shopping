package dto

import "mime/multipart"

type CreateRecipeRequest struct {
	Recipe       Recipe        `json:"recipe" form:"recipe"`
	Ingredients  []Ingredient  `json:"ingredients"`
	Instructions []Instruction `json:"instructions"`
	Image *multipart.FileHeader `form:"image"`
}

type Ingredient struct {
	Name            string  `json:"name"`
	Quantity        float64 `json:"quantity"`
	MeasurementUnit string  `json:"measurement_unit"`
}

type Instruction struct {
	Content   string `json:"content"`
	Numbering int    `json:"numbering"`
}

type RecipeResponse struct {
	RecipeID     uint                         `json:"recipe_id"`
	Recipe       Recipe                       `json:"recipe" form:"recipe"`
	Ingredients  []*RecipeIngredientResponse  `json:"ingredients"`
	Instructions []*RecipeInstructionResponse `json:"instructions"`
	ImagePath    string                       `json:"image_data"`
}

type Recipe struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CookTime    int    `json:"cook_time"`
	Servings    int    `json:"servings"`
}

type RecipeIngredientResponse struct {
	IngredientID       uint    `json:"ingredient_id"`
	RecipeIngredientID uint    `json:"recipe_ingredient_id"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	Quantity           float64 `json:"quantity"`
	MeasurementUnit    string  `json:"measurement_unit"`
}

type RecipeInstructionResponse struct {
	InstructionID       uint   `json:"instruction_id"`
	RecipeInstructionID uint   `json:"recipe_instruction_id"`
	Content             string `json:"content"`
	Numbering           uint   `json:"numbering"`
}
