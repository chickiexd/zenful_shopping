package service

import (
	// "log"
	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/store"

	"gorm.io/gorm"
)

type ingredientService struct {
	storage *store.Storage
}

func (s *ingredientService) Create(ingredient *dto.CreateIngredientRequest) (*store.Ingredient, error) {
	var created_ingredient *store.Ingredient
	err := s.storage.DB.Transaction(func(tx *gorm.DB) error {
		ingredientsRepo := s.storage.Ingredients.(*store.IngredientRepository).WithTransaction(tx)
		measurementUnitsRepo := s.storage.MeasurementUnits.(*store.MeasurementRepository).WithTransaction(tx)
		foodGroupsRepo := s.storage.FoodGroups.(*store.FoodGroupRepository).WithTransaction(tx)

		var measurement store.MeasurementUnit
		if ingredient.Measurement.MeasurementUnitID == 0 {
			measurement = store.MeasurementUnit{Name: ingredient.Measurement.Name}
			if err := measurementUnitsRepo.Create(&measurement); err != nil {
				return err
			}
		} else {
			measurement = ingredient.Measurement
		}

		var food_group store.FoodGroup
		if ingredient.FoodGroup.FoodGroupID == 0 {
			food_group = store.FoodGroup{Name: ingredient.FoodGroup.Name}
			if err := foodGroupsRepo.Create(&food_group); err != nil {
				return err
			}
		} else {
			food_group = ingredient.FoodGroup
		}

		if ingredient.IngredientID == 0 {
			created_ingredient = &store.Ingredient{
				Name: ingredient.Name,
			}
			if err := ingredientsRepo.Create(created_ingredient); err != nil {
				return err
			}
		} else {
			result, err := ingredientsRepo.GetByID(ingredient.IngredientID)
			if err != nil {
				return err
			}
			created_ingredient = result
		}

		if err := ingredientsRepo.CreateMeasurementUnitAssociation(created_ingredient, &measurement); err != nil {
			return err
		}
		if err := ingredientsRepo.CreateFoodGroupAssociation(created_ingredient, &food_group); err != nil {
			return err
		}

		if err := tx.Preload("MeasurementUnits").Preload("FoodGroups").First(&created_ingredient, created_ingredient.IngredientID).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return created_ingredient, nil
}

func (s *ingredientService) GetAll() ([]store.Ingredient, error) {
	ingredients, err := s.storage.Ingredients.GetAllTest()
	return ingredients, err
}

func (s *ingredientService) GetByID(id uint) (*store.Ingredient, error) {
	ingredient, err := s.storage.Ingredients.GetByID(id)
	if err != nil {
		return nil, err
	}
	return ingredient, nil
}
