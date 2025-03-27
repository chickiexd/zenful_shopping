package service

import (
	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/store"
)

type recipeService struct {
	storage *store.Storage
}

func (s *recipeService) Create(recipe *dto.CreateRecipeRequest) (*store.Recipe, error) {
	return nil, nil
}

func (s *recipeService) GetAll() ([]store.Recipe, error) {
	recipes, err := s.storage.Recipes.GetAll()
	return recipes, err
}
