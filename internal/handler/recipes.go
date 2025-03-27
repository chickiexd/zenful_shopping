package handler

import (
	"net/http"
	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/service"
	"zenful_shopping_backend/utils"
)

type RecipeHandler struct {
	service *service.Service
}


func (h *RecipeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.service.Recipes.GetAll();
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err.Error())
	}
	utils.WriteJSON(w, http.StatusOK, recipes)
}

func (h *RecipeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var recipe dto.CreateRecipeRequest
	if err := utils.ReadJSON(w, r, &recipe); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	created_recipe, err := h.service.Recipes.Create(&recipe)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, created_recipe)
}
