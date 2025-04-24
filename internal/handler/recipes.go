package handler

import (
	"encoding/json"
	"net/http"
	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/service"
	"zenful_shopping_backend/utils"
)

type RecipeHandler struct {
	service *service.Service
}

func (h *RecipeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.service.Recipes.GetAll()
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err.Error())
	}
	utils.WriteJSON(w, http.StatusOK, recipes)
}

func (h *RecipeHandler) AddToShoppingList(w http.ResponseWriter, r *http.Request) {
	var recipe_id dto.AddRecipeToShoppingListRequest
	err := utils.ReadJSON(w, r, &recipe_id)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err.Error())
	}
	err = h.service.Recipes.AddToShoppingList(recipe_id.RecipeID)
	utils.WriteJSON(w, http.StatusOK, recipe_id)
}

func (h *RecipeHandler) RemoveFromShoppingList(w http.ResponseWriter, r *http.Request) {
	var recipe_id dto.AddRecipeToShoppingListRequest
	err := utils.ReadJSON(w, r, &recipe_id)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err.Error())
	}
	err = h.service.Recipes.RemoveFromShoppingList(recipe_id.RecipeID)
	utils.WriteJSON(w, http.StatusOK, recipe_id)
}

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
