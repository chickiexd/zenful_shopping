package service

import (
	"context"
	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/chickiexd/zenful_shopping/internal/store"
)

type Service struct {
	Recipes interface {
		GetAll() ([]dto.RecipeResponse, error)
		Create(*dto.CreateRecipeRequest) (*dto.RecipeResponse, error)
		AddToShoppingList(uint) error
		RemoveFromShoppingList(uint) error
	}
	Ingredients interface {
		Create(*dto.CreateIngredientRequest) (*dto.IngredientResponse, error)
		GetByID(uint) (*dto.IngredientResponse, error)
		GetAll() ([]dto.IngredientResponse, error)
		AddToShoppingList(*dto.AddIngredientToShoppingListRequest) (*dto.ShoppingListItemResponse, error)
	}
	Pantry interface {
		GetAll() ([]dto.PantryIngredientResponse, error)
		Create(uint) error
		Delete(uint) error
		DeleteAll() error
	}
	MeasurmentUnits interface {
		Create(*dto.CreateMeasurementUnit) error
		GetAll() ([]store.MeasurementUnit, error)
	}
	FoodGroups interface {
		Create(*dto.CreateFoodGroupRequest) (*store.FoodGroup, error)
		GetAll() ([]store.FoodGroup, error)
	}
	Users interface {
		Create(context.Context, *store.User) error
	}
	ChatGPT interface {
		ParseRecipe(string) (*dto.ParsedRecipe, error)
	}
	ShoppingList interface {
		GetAll() ([]dto.ShoppingListResponse, error)
		RemoveItemFromShoppingList(uint) error
		RemoveAllItemsFromShoppingList(uint) error
	}
	KeepSync interface {
		SyncShoppingLists() error
	}
}

func NewService(store *store.Storage) Service {
	openAIService := NewOpenAIService(store)
	return Service{
		Recipes:         &recipeService{store},
		Users:           &userService{store},
		Ingredients:     &ingredientService{store},
		MeasurmentUnits: &MeasurementUnitService{store},
		FoodGroups:      &FoodGroupService{store},
		ShoppingList:    &ShoppingListService{store},
		Pantry:          &pantryService{store},
		KeepSync:        &keepSyncService{store},
		ChatGPT:         openAIService,
	}
}
