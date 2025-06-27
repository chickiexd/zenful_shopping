package service

import (
	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/chickiexd/zenful_shopping/internal/logger"
	"github.com/chickiexd/zenful_shopping/internal/store"
	"gorm.io/gorm"
)

type pantryService struct {
	storage *store.Storage
}

func (s *pantryService) GetAll() ([]dto.PantryIngredientResponse, error) {
	pantryIngredients, err := s.storage.Pantry.GetAll()
	if err != nil {
		logger.Log.Errorw("Error fetching pantry ingredients", "error", err)
		return nil, err
	}
	var pantryResponses []dto.PantryIngredientResponse
	for _, pantryIngredient := range pantryIngredients {
		response := dto.PantryIngredientResponse{
			IngredientID:   pantryIngredient.IngredientID,
			IngredientName: pantryIngredient.Ingredient.Name,
		}
		pantryResponses = append(pantryResponses, response)
	}
	return pantryResponses, nil
}

func (s *pantryService) Create(ingredient_id uint) error {
	pantry_ingredient := &store.PantryIngredient{
		IngredientID: ingredient_id,
	}
	existing_entry, err := s.storage.Pantry.GetByIngredientID(ingredient_id)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Log.Errorw("Error checking existing ingredient in pantry", "error", err, "ingredient_id", ingredient_id)
		return err
	}
	if existing_entry != nil {
		logger.Log.Warnw("Ingredient already exists in pantry", "ingredient_id", ingredient_id)
		return nil
	}
	_, err = s.storage.Pantry.Create(pantry_ingredient)
	if err != nil {
		logger.Log.Errorw("Error adding ingredient to pantry", "error", err, "ingredient_id", ingredient_id)
		return err
	}
	return nil
}

func (s *pantryService) Delete(ingredient_id uint) error {
	err := s.storage.Pantry.DeleteByIngredientID(ingredient_id)
	if err != nil {
		logger.Log.Errorw("Error removing ingredient from pantry", "error", err, "ingredient_id", ingredient_id)
		return err
	}
	return nil
}

func (s *pantryService) DeleteAll() error {
	err := s.storage.Pantry.DeleteAll()
	if err != nil {
		logger.Log.Errorw("Error removing all ingredients from pantry", "error", err)
		return err
	}
	return nil
}
