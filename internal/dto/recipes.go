package dto

import "mime/multipart"

type RecipeResponse struct {
	RecipeID     uint                       `json:"recipe_id"`
	Title        string                     `json:"title"`
	Description  string                     `json:"description"`
	Public       bool                       `json:"public"`
	CookTime     int                        `json:"cook_time"`
	Servings     int                        `json:"servings"`
	ImagePath    string                     `json:"image_path"`
	MealType     uint                       `json:"meal_type"`
	Ingredients  []RecipeIngredientResponse `json:"ingredients"`
	Instructions []InstructionResponse      `json:"instructions"`
	CreatedAt    string                     `json:"created_at"`
	UpdatedAt    string                     `json:"updated_at"`
}

type RecipeIngredientResponse struct {
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

type AddRecipeToShoppingListRequest struct {
	RecipeID uint `json:"recipe_id"`
}

// returned from the openai api
type ParsedRecipeInformation struct {
	Recipe       recipeRequest           `json:"recipe" form:"recipe"`
	Ingredients  []IngredientInformation `json:"ingredients"`
	Instructions []instructionRequest    `json:"instructions"`
}

type IngredientInformation struct {
	Name            string  `json:"name"`
	Quantity        float64 `json:"quantity"`
	MeasurementUnit string  `json:"measurement_unit"`
}

type ParsedIngredientInformation struct {
	IngredientName        string   `json:"ingredient_name"`
	IngredientDescription string   `json:"ingredient_description"`
	MeasurementUnits      []string `json:"measurement_units"`
	FoodGroups            []string `json:"food_groups"`
}

type ParsedMultipleIngredientInformation struct {
	IngredientName        string   `json:"ingredient_name"`
	IngredientDescription string   `json:"ingredient_description"`
	MeasurementUnits      []string `json:"measurement_units"`
	FoodGroups            []string `json:"food_groups"`
	Quantity              float64  `json:"quantity"`
	MeasurementUnitID     uint     `json:"measurement_unit_id"`
}

// returned to the client
type ParsedRecipe struct {
	Recipe         recipeRequest              `json:"recipe" form:"recipe"`
	Ingredients    []RecipeIngredientResponse `json:"ingredients"`
	Instructions   []instructionRequest       `json:"instructions"`
	NewIngredients []ParsedIngredient         `json:"new_ingredients"`
	NewFoodGrops   []ParsedFoodGroup          `json:"new_food_groups"`
}

type ParsedIngredient struct {
	Name                 string            `json:"name"`
	Quantity             float64           `json:"quantity"`
	MeasurementUnitNames []string          `json:"measurement_unit_names"`
	MeasurementUnitID    uint              `json:"measurement_unit_id"`
	ParsedFoodGroups     []ParsedFoodGroup `json:"food_groups"`
}

type ParsedFoodGroup struct {
	Name          string `json:"name"`
	ShoppingLists string `json:"shopping_list"`
}
