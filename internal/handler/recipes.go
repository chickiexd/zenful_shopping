package handler

import (
	"encoding/json"
	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/chickiexd/zenful_shopping/internal/service"
	"github.com/chickiexd/zenful_shopping/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type RecipeHandler struct {
	service *service.Service
}

// GetAll godoc
//
//	@Summary		Get all recipes
//	@Description	Get all recipes from the database
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		dto.RecipeResponse
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/recipes [get]
func (h *RecipeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.service.Recipes.GetAll()
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err.Error())
	}
	utils.WriteJSON(w, http.StatusOK, recipes)
}

// AddToShoppingList godoc
//
//	@Summary		Add a recipe to the shopping list
//	@Description	Add a recipe to the shopping list by its ID
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Recipe ID"
//	@Success		200	{integer}	uint	"Recipe ID"
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/recipes/add/{id} [post]
func (h *RecipeHandler) AddToShoppingList(w http.ResponseWriter, r *http.Request) {
	recipeIDStr := chi.URLParam(r, "id")

	parsedID, err := strconv.ParseUint(recipeIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}
	recipeID := uint(parsedID)

	if err := h.service.Recipes.AddToShoppingList(recipeID); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, recipeID)
}

// RemoveFromShoppingList godoc
//
//	@Summary		Remove a recipe from the shopping list
//	@Description	Remove a recipe from the shopping list by its ID
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.AddRecipeToShoppingListRequest	true	"Recipe ID"
//	@Success		200		{object}	dto.AddRecipeToShoppingListRequest
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/recipes/remove [post]
func (h *RecipeHandler) RemoveFromShoppingList(w http.ResponseWriter, r *http.Request) {
	var recipe_id dto.AddRecipeToShoppingListRequest
	err := utils.ReadJSON(w, r, &recipe_id)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err.Error())
	}
	err = h.service.Recipes.RemoveFromShoppingList(recipe_id.RecipeID)
	utils.WriteJSON(w, http.StatusOK, recipe_id)
}

// Create godoc
//
//	@Summary		Create a new recipe
//	@Description	Accepts multipart/form-data with JSON fields and an image to create a new recipe
//	@Tags			recipes
//	@Accept			multipart/form-data
//	@Produce		plain
//	@Param			recipe			formData	string	true	"Recipe JSON string"
//	@Param			ingredients		formData	string	true	"Ingredients JSON array"
//	@Param			instructions	formData	string	true	"Instructions JSON array"
//	@Param			image			formData	file	true	"Recipe image file"
//	@Success		201				{string}	string	"Recipe created successfully"
//	@Failure		400				{string}	string	"Bad request"
//	@Failure		404				{string}	string	"Not found"
//	@Failure		500				{string}	string	"Internal server error"
//	@Router			/recipes [post]
func (h *RecipeHandler) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB max memory
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	var req dto.CreateRecipeRequest

	if err := json.Unmarshal([]byte(r.FormValue("recipe")), &req.Recipe); err != nil {
		http.Error(w, "Invalid recipe field", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal([]byte(r.FormValue("ingredients")), &req.Ingredients); err != nil {
		http.Error(w, "Invalid ingredients field", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal([]byte(r.FormValue("instructions")), &req.Instructions); err != nil {
		http.Error(w, "Invalid instructions field", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error reading image", http.StatusBadRequest)
		return
	}
	defer file.Close()
	req.Image = file
	req.ImageHeader = fileHeader

	h.service.Recipes.Create(&req)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Recipe created successfully"))
}
