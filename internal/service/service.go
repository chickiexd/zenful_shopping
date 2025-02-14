package service

import (
	"context"
	"zenful_shopping_backend/internal/store"
)

type Service struct {
	Recipes interface {
		Create(context.Context, *store.Recipe) error
	}

	Users interface {
		Create(context.Context, *store.User) error
	}
}

func NewService(store store.Storage) Service {
	return Service{
		Recipes: &recipeService{store},
		Users:   &userService{store},
	}
}
