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
		GetByID(int64) (*dto.IngredientResponse, error)
		GetAll() ([]dto.IngredientResponse, error)
	}
	Users interface {
		Create(context.Context, *store.User) error
	}
}

func NewService(store *store.Storage) Service {
	return Service{
		Recipes: &recipeService{store},
		Users:   &userService{store},
		Ingredients: &ingredientService{store},
	}
}
