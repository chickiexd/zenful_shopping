package handler

import (
	"log"
	"net/http"
	"github.com/chickiexd/zenful_shopping/internal/service"
	"github.com/chickiexd/zenful_shopping/utils"
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
		return
	}

	parsed_recipe, err := h.service.ChatGPT.ParseRecipe(submitted_recipe.Text)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("berfore WREITETIE Parsed recipe: %+v", parsed_recipe)
	err = utils.WriteJSON(w, http.StatusOK, parsed_recipe)
	if err != nil {
		log.Printf("error writing JSON response: %v", err)
	}

}

