package service

import (
	"context"
	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/store"
)

type Service struct {
	Recipes interface {
		Create(context.Context, *store.Recipe) error
	}
	Ingredients interface {
		Create(*dto.CreateIngredientRequest) error
		GetByID(uint) (*store.Ingredient, error)
		GetAll() ([]store.Ingredient, error)
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
}

func NewService(store *store.Storage) Service {
	openAIService := NewOpenAIService()
	return Service{
		Recipes: &recipeService{store},
		Users:   &userService{store},
		Ingredients: &ingredientService{store},
		MeasurmentUnits: &MeasurementUnitService{store},
		FoodGroups: &FoodGroupService{store},
		ChatGPT: openAIService,
	}
}
