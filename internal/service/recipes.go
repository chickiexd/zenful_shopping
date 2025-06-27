package service

import (
	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/chickiexd/zenful_shopping/internal/store"
	"github.com/chickiexd/zenful_shopping/utils"
	"log"
	"os"

	"gorm.io/gorm"
)

type recipeService struct {
	storage *store.Storage
}

func (s *recipeService) Create(req *dto.CreateRecipeRequest) (*dto.RecipeResponse, error) {
	var created_recipe *store.Recipe

	image_path, err := utils.SaveImageLocally(req.Image, req.ImageHeader)
	if err != nil {
		return nil, err
	}

	err = s.storage.DB.Transaction(func(tx *gorm.DB) error {
		recipeRepo := s.storage.Recipes.(*store.RecipeRepository).WithTransaction(tx)
		instructionRepo := s.storage.Instructions.(*store.InstructionRepository).WithTransaction(tx)

		recipe := &store.Recipe{
			Title:       req.Recipe.Title,
			Description: req.Recipe.Description,
			CookTime:    req.Recipe.CookTime,
			Servings:    req.Recipe.Servings,
			Public:      req.Recipe.Public,
			ImagePath:   image_path,
			MealTypeID:  uint(req.Recipe.MealType),
		}

		if err := recipeRepo.Create(recipe); err != nil {
			return err
		}
		created_recipe = recipe

		for _, instr := range req.Instructions {
			instur_req := &store.Instruction{
				StepNumber:  instr.Numbering,
				Description: instr.Content,
				RecipeID:    recipe.RecipeID,
			}
			if err := instructionRepo.Create(instur_req); err != nil {
				return err
			}
		}

		for _, ing := range req.Ingredients {
			recipe_ingredient := &store.RecipeIngredient{
				RecipeID:          recipe.RecipeID,
				IngredientID:      ing.IngredientID,
				MeasurementUnitID: ing.MeasurementUnitID,
				Quantity:          ing.Quantity,
			}
			if err := recipeRepo.CreateIngredientAssociation(recipe_ingredient); err != nil {
				return err
			}
		}

		if err := tx.Preload("MealType").Preload("Instructions").Preload("RecipeIngredients.Ingredient").Preload("RecipeIngredients.MeasurementUnit").First(&created_recipe, created_recipe.RecipeID).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		os.Remove(image_path)
		return nil, err
	}

	instructions := make([]dto.InstructionResponse, len(created_recipe.Instructions))
	for i, instr := range created_recipe.Instructions {
		instructions[i] = dto.InstructionResponse{
			Description: instr.Description,
			StepNumber:  instr.StepNumber,
		}
	}
	ingredients := make([]dto.RecipeIngredientResponse, len(created_recipe.RecipeIngredients))
	for i, ing := range created_recipe.RecipeIngredients {
		ingredients[i] = dto.RecipeIngredientResponse{
			IngredientID:      ing.IngredientID,
			Quantity:          ing.Quantity,
			MeasurementUnitID: ing.MeasurementUnitID,
		}
	}

	recipe_response := &dto.RecipeResponse{
		RecipeID:     created_recipe.RecipeID,
		Title:        created_recipe.Title,
		Description:  created_recipe.Description,
		Public:       created_recipe.Public,
		CookTime:     created_recipe.CookTime,
		Servings:     created_recipe.Servings,
		ImagePath:    created_recipe.ImagePath,
		MealType:     created_recipe.MealTypeID,
		CreatedAt:    created_recipe.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    created_recipe.UpdatedAt.Format("2006-01-02 15:04:05"),
		Ingredients:  ingredients,
		Instructions: instructions,
	}
	return recipe_response, nil
}

func (s *recipeService) GetAll() ([]dto.RecipeResponse, error) {
	recipes, err := s.storage.Recipes.GetAll()
	if err != nil {
		return nil, err
	}

	recipes_response := make([]dto.RecipeResponse, len(recipes))
	for i, recipe := range recipes {
		instructions := make([]dto.InstructionResponse, len(recipe.Instructions))
		for j, instr := range recipe.Instructions {
			instructions[j] = dto.InstructionResponse{
				InstructionID: instr.InstructionID,
				Description:   instr.Description,
				StepNumber:    instr.StepNumber,
			}
		}
		ingredients := make([]dto.RecipeIngredientResponse, len(recipe.RecipeIngredients))
		for j, ing := range recipe.RecipeIngredients {
			ingredients[j] = dto.RecipeIngredientResponse{
				IngredientID:      ing.IngredientID,
				Quantity:          ing.Quantity,
				MeasurementUnitID: ing.MeasurementUnitID,
			}
		}
		recipes_response[i] = dto.RecipeResponse{
			RecipeID:     recipe.RecipeID,
			Title:        recipe.Title,
			Description:  recipe.Description,
			Public:       recipe.Public,
			CookTime:     recipe.CookTime,
			Servings:     recipe.Servings,
			ImagePath:    recipe.ImagePath,
			MealType:     recipe.MealTypeID,
			CreatedAt:    recipe.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    recipe.UpdatedAt.Format("2006-01-02 15:04:05"),
			Ingredients:  ingredients,
			Instructions: instructions,
		}
	}
	return recipes_response, err
}

func (s *recipeService) AddToShoppingList(recipe_id uint) error {
	err := s.storage.DB.Transaction(func(tx *gorm.DB) error {
		shoppingListRepo := s.storage.ShoppingLists.WithTransaction(tx)

		recipe, err := s.storage.Recipes.GetByID(recipe_id)
		if err != nil {
			return err
		}
		for _, ing := range recipe.RecipeIngredients {
			pantry_ingredient, err := s.storage.Pantry.GetByIngredientID(ing.IngredientID)
			if pantry_ingredient != nil {
				continue
			}
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			food_groups, err := s.storage.Ingredients.GetFoodGroupsByID(ing.IngredientID)
			if err != nil {
				return err
			}
			for _, food_group := range food_groups {
				shopping_list, err := shoppingListRepo.GetByFoodGroupID(food_group.FoodGroupID)
				if err != nil {
					return err
				}
				if len(shopping_list) != 1 {
					log.Println("0 or more than 1 shoppinlist found for: ", food_group.Name)
				} else {
					items, err := shoppingListRepo.GetItemsByIngredientID(ing.IngredientID)
					if err != nil {
						return err
					}
					if len(items) > 1 {
						log.Println("more than 1 ingredient entry found in shoppinglist for: ", ing.IngredientID)
						return err
					}
					if len(items) == 0 {
						item := &store.ShoppingListItem{
							IngredientID:      ing.IngredientID,
							MeasurementUnitID: ing.MeasurementUnitID,
							Quantity:          ing.Quantity,
							ShoppingListID:    shopping_list[0].ShoppingListID,
						}
						if err := shoppingListRepo.CreateItemAssociation(item); err != nil {
							return err
						}
					} else {
						item := items[0]
						// TODO check for different measurement units
						item.Quantity += ing.Quantity
						if err := shoppingListRepo.UpdateItemAssociation(&item); err != nil {
							return err
						}
					}
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Println("Error adding to shopping list: ", err)
		return err
	}
	return nil
}

func (s *recipeService) RemoveFromShoppingList(recipe_id uint) error {
	err := s.storage.DB.Transaction(func(tx *gorm.DB) error {
		shoppingListRepo := s.storage.ShoppingLists.WithTransaction(tx)

		recipe, err := s.storage.Recipes.GetByID(recipe_id)
		if err != nil {
			return err
		}
		for _, ing := range recipe.RecipeIngredients {
			items, err := shoppingListRepo.GetItemsByIngredientID(ing.IngredientID)
			if err != nil {
				return err
			}
			if len(items) > 1 {
				log.Println("more than 1 ingredient entry found in shoppinglist for: ", ing.IngredientID)
				return err
			}
			if len(items) == 0 {
				log.Println("0 items found in shoppinglist for: ", ing.IngredientID)
				continue
			} else {
				item := items[0]
				item.Quantity -= ing.Quantity
				if item.Quantity <= 0 {
					item.Quantity = 0
					if err := shoppingListRepo.DeleteItemAssociation(&item); err != nil {
						return err
					}
				} else {
					if err := shoppingListRepo.UpdateItemAssociation(&item); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Println("Error removing from shopping list: ", err)
		return err
	}
	return nil
}
