package service

import (
	"fmt"
	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/store"

	"gorm.io/gorm"
)

type ingredientService struct {
	storage *store.Storage
}

func (s *ingredientService) Create(ingredient *dto.CreateIngredientRequest) (*dto.IngredientResponse, error) {
	var created_ingredient *store.Ingredient
	err := s.storage.DB.Transaction(func(tx *gorm.DB) error {
		ingredientsRepo := s.storage.Ingredients.(*store.IngredientRepository).WithTransaction(tx)
		measurementUnitsRepo := s.storage.MeasurementUnits.(*store.MeasurementRepository).WithTransaction(tx)
		foodGroupsRepo := s.storage.FoodGroups.(*store.FoodGroupRepository).WithTransaction(tx)

		created_ingredient = &store.Ingredient{
			Name: ingredient.Name,
		}
		if err := ingredientsRepo.Create(created_ingredient); err != nil {
			return err
		}

		for _, new_measurement := range ingredient.Measurements {
			var measurement *store.MeasurementUnit
			if new_measurement.MeasurementUnitID == 0 {
				if new_measurement.Name == "" {
					return fmt.Errorf("measurement name cannot be empty")
				}
				measurement = &store.MeasurementUnit{Name: new_measurement.Name}
				if err := measurementUnitsRepo.Create(measurement); err != nil {
					return err
				}
			} else {
				result, err := measurementUnitsRepo.GetByID(new_measurement.MeasurementUnitID)
				if err != nil {
					return err
				}
				measurement = result
			}
			if err := ingredientsRepo.CreateMeasurementUnitAssociation(created_ingredient, measurement); err != nil {
				return err
			}
		}

		for _, new_food_group := range ingredient.FoodGroups {
			var food_group *store.FoodGroup
			result, err := foodGroupsRepo.GetByName(new_food_group)
			if err != nil && err == gorm.ErrRecordNotFound {
				if new_food_group == "" {
					return fmt.Errorf("food group name cannot be empty")
				}
				food_group = &store.FoodGroup{Name: new_food_group}
				if err := foodGroupsRepo.Create(food_group); err != nil {
					return err
				}
			} else if err != nil {
				return err
			} else {
				food_group = result
			}
			if err := ingredientsRepo.CreateFoodGroupAssociation(created_ingredient, food_group); err != nil {
				return err
			}
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

func (s *ingredientService) AddToShoppingList(ingredient *dto.AddIngredientToShoppingListRequest) (*dto.ShoppingListItemResponse, error) {
	item := &store.ShoppingListItem{
		IngredientID:      ingredient.IngredientID,
		MeasurementUnitID: ingredient.MeasurementUnitID,
		Quantity:          ingredient.Quantity,
		ShoppingListID:    ingredient.ShoppingListID,
	}
	err := s.storage.ShoppingLists.CreateItemAssociation(item)
	if err != nil {
		return nil, err
	}
	ingredient_response := &dto.ShoppingListItemResponse{
		ShoppingListID:     item.ShoppingListID,
		ShoppingListItemID: item.ShoppingListItemID,
		IngredientID:       item.IngredientID,
		Quantity:           item.Quantity,
		MeasurementUnitID:  item.MeasurementUnitID,
	}
	return ingredient_response, nil
}
