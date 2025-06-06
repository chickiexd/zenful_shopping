package handler

import (
	"log"
	"net/http"

	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/chickiexd/zenful_shopping/internal/service"
	"github.com/chickiexd/zenful_shopping/utils"
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
		return
	}
	created_ingredient, err := h.service.Ingredients.Create(&ingredient)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("created_ingredient: ", created_ingredient)
	utils.WriteJSON(w, http.StatusOK, created_ingredient)
}

func (h *IngredientHandler) AddToShoppingList(w http.ResponseWriter, r *http.Request) {
	var ingredient dto.AddIngredientToShoppingListRequest
	if err := utils.ReadJSON(w, r, &ingredient); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.service.Ingredients.AddToShoppingList(&ingredient)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, item)
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
