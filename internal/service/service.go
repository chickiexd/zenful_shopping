package service

import (
	"context"
	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/store"
)

type Service struct {
	Recipes interface {
		GetAll() ([]dto.RecipeResponse, error)
		Create(*dto.CreateRecipeRequest) (*dto.RecipeResponse, error)
	}
	Ingredients interface {
		Create(*dto.CreateIngredientRequest) (*dto.IngredientResponse, error)
		GetByID(uint) (*dto.IngredientResponse, error)
		GetAll() ([]dto.IngredientResponse, error)
	}
	MeasurmentUnits interface {
		Create(*dto.CreateMeasurementUnit) error
		GetAll() ([]store.MeasurementUnit, error)
	}
	FoodGroups interface {
		Create(*dto.CreateFoodGroup) error
		GetAll() ([]store.FoodGroup, error)
	}
	Users interface {
		Create(context.Context, *store.User) error
	}
	ChatGPT interface {
		ParseRecipe(string) (*dto.CreateRecipeRequest, error)
	}
	Images interface {
		Get(string) ([]byte, error)
	}
}

func NewService(store *store.Storage) Service {
	openAIService := NewOpenAIService()
	return Service{
		Recipes: &recipeService{store},
		Users:   &userService{store},
		Ingredients: &ingredientService{store},
		MeasurmentUnits: &MeasurementUnitService{store},
		FoodGroups: &FoodGroupService{store},
		Images: &ImageService{},
		ChatGPT: openAIService,
	}
}
