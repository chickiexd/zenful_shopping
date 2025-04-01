package service

import (
	"os"
	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/store"
	"zenful_shopping_backend/utils"

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
	ingredients := make([]dto.IngredientResponse, len(created_recipe.RecipeIngredients))
	for i, ing := range created_recipe.RecipeIngredients {
		ingredients[i] = dto.IngredientResponse{
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

func (s *recipeService) GetAll() ([]store.Recipe, error) {
	recipes, err := s.storage.Recipes.GetAll()
	return recipes, err
}
