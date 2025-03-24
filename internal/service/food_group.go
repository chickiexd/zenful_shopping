package service

import (
	// "log"
	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/store"
)

type FoodGroupService struct {
	storage *store.Storage
}

func (s *FoodGroupService) Create(ingredient *dto.CreateFoodGroup) error {
	return s.storage.FoodGroups.Create(nil)
}

func (s *FoodGroupService) GetAll() ([]store.FoodGroup, error) {
	food_groups, err := s.storage.FoodGroups.GetAll()
	return food_groups, err
}

// func (s *FoodGroupService) GetByID(id uint) (*store.Measurement, error) {
// 	ingredient, err := s.storage.FoodGroups.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ingredient, nil
// }
