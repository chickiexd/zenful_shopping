package service

import (
	"fmt"
	"log"
	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/store"
)

type ingredientService struct {
	storage *store.Storage
}

func (s *ingredientService) Create(ingredient *dto.CreateIngredientRequest) error {
	log.Println("service create ingredient_id")
	return s.storage.Ingredients.Create(nil)
}

func (s *ingredientService) GetAll() ([]dto.IngredientResponse, error) {
	ingredients, err := s.storage.Ingredients.GetAll()
	if err != nil {
		return nil, err
	}
	fmt.Println(ingredients)
	// type IngredientResponse struct {
	// 	IngredientID uint                    `json:"ingredient_id"`
	// 	Name         string                  `json:"name"`
	// 	Measurements []IngredientMeasurement `json:"measurements"`
	// 	FoodGroup    IngredientFoodGroup     `json:"food_group"`
	// }
	// get food group
	// get measurements

	return nil, nil
}

func (s *ingredientService) GetByID(id int64) (*dto.IngredientResponse, error) {
	ingredients, err := s.storage.Ingredients.GetByID(id)
	log.Println(ingredients, err)
	return nil, nil
}
