package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/chickiexd/zenful_shopping/internal/store"

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

		if strings.TrimSpace(ingredient.Name) == "" {
			return validationErrorf("ingredient name must not be empty")
		}
		if _, err := ingredientsRepo.GetByName(ingredient.Name); err == nil {
			return validationErrorf("an ingredient named %q already exists", ingredient.Name)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		created_ingredient = &store.Ingredient{
			Name:        strings.TrimSpace(ingredient.Name),
			Description: ingredient.Description,
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

		for _, food_group_name := range ingredient.FoodGroups {
			food_group, err := foodGroupsRepo.GetByName(food_group_name)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// creating it here would leave it without a shopping list, so the
				// ingredient would never be added to one
				return validationErrorf("food group %q does not exist yet; create it with a shopping list first", food_group_name)
			}
			if err != nil {
				return err
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
	item, err := addItemToShoppingList(s.storage.ShoppingLists.WithTransaction(s.storage.DB), &store.ShoppingListItem{
		IngredientID:      ingredient.IngredientID,
		MeasurementUnitID: ingredient.MeasurementUnitID,
		Quantity:          ingredient.Quantity,
		ShoppingListID:    ingredient.ShoppingListID,
	})
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

// AddAlias remembers req.Alias as another name for an existing ingredient.
// Aliases equal to the ingredient's own name are skipped since they already match.
func (s *ingredientService) AddAlias(req *dto.CreateIngredientAliasRequest) error {
	ingredient, err := s.storage.Ingredients.GetByID(req.IngredientID)
	if err != nil {
		return err
	}
	alias := normalizeName(req.Alias)
	if alias == "" || alias == normalizeName(ingredient.Name) {
		return nil
	}
	return s.storage.Aliases.Upsert(&store.IngredientAlias{Alias: alias, IngredientID: ingredient.IngredientID})
}
