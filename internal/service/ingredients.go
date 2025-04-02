package service

import (
	"fmt"
	"log"
	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/store"

	"gorm.io/gorm"
)

type ingredientService struct {
	storage *store.Storage
}

func (s *ingredientService) Create(ingredient *dto.CreateIngredientRequest) (*dto.IngredientResponse, error) {
	log.Println("service create ingredient_id")
	log.Println(ingredient)
	var created_ingredient *store.Ingredient
	err := s.storage.DB.Transaction(func(tx *gorm.DB) error {
		ingredientsRepo := s.storage.Ingredients.(*store.IngredientRepository).WithTransaction(tx)
		measurementUnitsRepo := s.storage.MeasurementUnits.(*store.MeasurementRepository).WithTransaction(tx)
		foodGroupsRepo := s.storage.FoodGroups.(*store.FoodGroupRepository).WithTransaction(tx)

		var measurement *store.MeasurementUnit
		if ingredient.Measurement.MeasurementUnitID == 0 {
			if ingredient.Measurement.Name == "" {
				return fmt.Errorf("measurement name cannot be empty")
			}
			measurement = &store.MeasurementUnit{Name: ingredient.Measurement.Name}
			if err := measurementUnitsRepo.Create(measurement); err != nil {
				return err
			}
		} else {
			result, err := measurementUnitsRepo.GetByID(ingredient.Measurement.MeasurementUnitID)
			if err != nil {
				return err
			}
			measurement = result
		}

		var food_group *store.FoodGroup
		if ingredient.FoodGroup.FoodGroupID == 0 {
			if ingredient.FoodGroup.Name == "" {
				return fmt.Errorf("food group name cannot be empty")
			}
			food_group = &store.FoodGroup{Name: ingredient.FoodGroup.Name}
			if err := foodGroupsRepo.Create(food_group); err != nil {
				return err
			}
		} else {
			result, err := foodGroupsRepo.GetByID(ingredient.FoodGroup.FoodGroupID)
			if err != nil {
				return err
			}
			food_group = result
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

		if err := ingredientsRepo.CreateMeasurementUnitAssociation(created_ingredient, measurement); err != nil {
			return err
		}
		if err := ingredientsRepo.CreateFoodGroupAssociation(created_ingredient, food_group); err != nil {
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

	created_ingredient_response := convert_store_to_dto_ingredient(created_ingredient)

	return created_ingredient_response, nil
}

func (s *ingredientService) GetAll() ([]dto.IngredientResponse, error) {
	ingredients, err := s.storage.Ingredients.GetAll()
	if err != nil {
		return nil, err
	}
	ingredients_response := make([]dto.IngredientResponse, len(ingredients))
	for i, ingredient := range ingredients {
		ingredients_response[i] = *convert_store_to_dto_ingredient(&ingredient)
	}
	return ingredients_response, err
}

func (s *ingredientService) GetByID(id uint) (*dto.IngredientResponse, error) {
	ingredient, err := s.storage.Ingredients.GetByID(id)
	if err != nil {
		return nil, err
	}
	ingredient_response := convert_store_to_dto_ingredient(ingredient)
	return ingredient_response, nil
}

func convert_store_to_dto_ingredient(ingredient *store.Ingredient) *dto.IngredientResponse {
	measurements := make([]dto.MeasurementUnitResponse, len(ingredient.MeasurementUnits))
	for i, measurement := range ingredient.MeasurementUnits {
		measurements[i] = dto.MeasurementUnitResponse{
			MeasurementUnitID: measurement.MeasurementUnitID,
			Name:              measurement.Name,
		}
	}
	food_groups := make([]dto.FoodGroupResponse, len(ingredient.FoodGroups))
	for i, food_group := range ingredient.FoodGroups {
		food_groups[i] = dto.FoodGroupResponse{
			FoodGroupID: food_group.FoodGroupID,
			Name:        food_group.Name,
		}
	}
	ingredient_response := &dto.IngredientResponse{
		IngredientID: ingredient.IngredientID,
		Name:         ingredient.Name,
		Description:  ingredient.Description,
		Measurements: measurements,
		FoodGroups:   food_groups,
	}
	return ingredient_response
}
