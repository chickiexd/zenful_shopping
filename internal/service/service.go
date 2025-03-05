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
