package handler

import (
	"context"
	"zenful_shopping_backend/internal/service"

)

type RecipeHandler struct {
	service *service.Service
}

func (h *RecipeHandler) GetAll(c *context.Context) error {
	// recipes, err := h.service.GetAllRecipes()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, context.H{"error": err.Error()})
	// 	return err
	// }
	// c.JSON(http.StatusOK, recipes)
	return nil
}

// func (h *RecipeHandler) Create(c *context.Context, recipe *CreateRecipeRequest) {
// 	var recipe store.Recipe
// 	if err := c.ShouldBindJSON(&recipe); err != nil {
// 		c.JSON(http.StatusBadRequest, context.H{"error": "Invalid request"})
// 		return
// 	}
//
// 	err := h.service.Recipes.Create(context.Background(), &recipe)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, context.H{"error": "Failed to create recipe"})
// 		return
// 	}
//
// 	c.JSON(http.StatusCreated, context.H{"message": "Recipe created successfully", "recipe": recipe})
// }

// func (h *RecipeHandler) CreateRecipe(c *context.Context) {
// 	fmt.Println("inside CreateRecipeHandler")
// 	for key, values := range c.Request.PostForm {
// 		fmt.Printf("Form field: %s, Value: %s\n", key, values)
// 	}
// 	var recipe dto.CreateRecipeRequest
// 	if err := c.ShouldBind(&recipe); err != nil {
// 		fmt.Printf("Binding error: %v\n", err)
// 		c.JSON(http.StatusBadRequest, context.H{"error": err.Error()})
// 		return
// 	}
// 	if ingredientsJson := c.PostForm("ingredients"); ingredientsJson != "" {
// 		if err := json.Unmarshal([]byte(ingredientsJson), &recipe.Ingredients); err != nil {
// 			c.JSON(http.StatusBadRequest, context.H{"error": "Invalid ingredients format"})
// 			return
// 		}
// 	}
// 	if instructionsJson := c.PostForm("instructions"); instructionsJson != "" {
// 		if err := json.Unmarshal([]byte(instructionsJson), &recipe.Instructions); err != nil {
// 			c.JSON(http.StatusBadRequest, context.H{"error": "Invalid instructions format"})
// 			return
// 		}
// 	}
// 	// save the image
// 	var imagePath string
// 	if recipe.Image != nil {
// 		path, err := saveImage(recipe.Image)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, context.H{"error": "image upload failed"})
// 			return
// 		}
// 		imagePath = path
// 	}
//
// 	if _, err := h.service.CreateRecipe(&recipe, imagePath); err != nil {
// 		c.JSON(http.StatusInternalServerError, context.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, recipe)
// }
