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

func (s *ingredientService) Create(ingredient *dto.CreateIngredientRequest) error {
	return s.storage.DB.Transaction(func(tx *gorm.DB) error {
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

		var ingredient_to_create *store.Ingredient
		if ingredient.IngredientID == 0 {
			ingredient_to_create = &store.Ingredient{
				Name: ingredient.Name,
			}
			if err := ingredientsRepo.Create(ingredient_to_create); err != nil {
				return err
			}
		} else {
			result, err := ingredientsRepo.GetByID(ingredient.IngredientID)
			if err != nil {
				return err
			}
			ingredient_to_create = result
		}

		if err := ingredientsRepo.CreateMeasurementUnitAssociation(ingredient_to_create, &measurement); err != nil {
			return err
		}
		if err := ingredientsRepo.CreateFoodGroupAssociation(ingredient_to_create, &food_group); err != nil {
			return err
		}
		return nil
	})
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
