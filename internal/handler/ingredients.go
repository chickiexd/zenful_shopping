package handler

import (
	// "log"
	"net/http"

	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/service"
	"zenful_shopping_backend/utils"
)

type IngredientHandler struct {
	service *service.Service
}

func (h *IngredientHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ingredients, err := h.service.Ingredients.GetAll();
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err.Error())
	}
	utils.WriteJSON(w, http.StatusOK, ingredients)
}

func (h *IngredientHandler) Create(w http.ResponseWriter, r *http.Request) {
	var ingredient dto.CreateIngredientRequest
	if err := utils.ReadJSON(w, r, &ingredient); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
	}
	h.service.Ingredients.Create(&ingredient)
}

// func (h *RecipeHandler) Create(c *gin.Context, recipe *CreateRecipeRequest) {
// 	var recipe store.Recipe
// 	if err := c.ShouldBindJSON(&recipe); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}
//
// 	err := h.service.Recipes.Create(context.Background(), &recipe)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe"})
// 		return
// 	}
//
// 	c.JSON(http.StatusCreated, gin.H{"message": "Recipe created successfully", "recipe": recipe})
// }

// func (h *RecipeHandler) CreateRecipe(c *gin.Context) {
// 	fmt.Println("inside CreateRecipeHandler")
// 	for key, values := range c.Request.PostForm {
// 		fmt.Printf("Form field: %s, Value: %s\n", key, values)
// 	}
// 	var recipe dto.CreateRecipeRequest
// 	if err := c.ShouldBind(&recipe); err != nil {
// 		fmt.Printf("Binding error: %v\n", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if ingredientsJson := c.PostForm("ingredients"); ingredientsJson != "" {
// 		if err := json.Unmarshal([]byte(ingredientsJson), &recipe.Ingredients); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredients format"})
// 			return
// 		}
// 	}
// 	if instructionsJson := c.PostForm("instructions"); instructionsJson != "" {
// 		if err := json.Unmarshal([]byte(instructionsJson), &recipe.Instructions); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructions format"})
// 			return
// 		}
// 	}
// 	// save the image
// 	var imagePath string
// 	if recipe.Image != nil {
// 		path, err := saveImage(recipe.Image)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "image upload failed"})
// 			return
// 		}
// 		imagePath = path
// 	}
//
// 	if _, err := h.service.CreateRecipe(&recipe, imagePath); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, recipe)
// }
