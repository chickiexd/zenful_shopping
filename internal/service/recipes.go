package service

import (
	"context"

	"zenful_shopping_backend/internal/store"
)

type recipeService struct {
	storage store.Storage
}

func (s *recipeService) Create(ctx context.Context, recipe *store.Recipe) error {
	return s.storage.Recipes.Create(ctx, recipe)
}
