package handler

import (
	"net/http"
	"zenful_shopping_backend/internal/service"
	"zenful_shopping_backend/utils"
)

type SubmittedRecipe struct {
	Text string `json:"text" binding:"required"`
}

type ChatGPTHandler struct {
	service *service.Service
}

func (h *ChatGPTHandler) ParseRecipe(w http.ResponseWriter, r *http.Request) {
	var submitted_recipe SubmittedRecipe
	if err := utils.ReadJSON(w, r, &submitted_recipe); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
	}

	parsed_recipe, err := h.service.ChatGPT.ParseRecipe(submitted_recipe.Text)
	if err != nil {
		return
	}
	utils.WriteJSON(w, http.StatusOK, parsed_recipe)
}

