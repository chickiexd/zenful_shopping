package handler

import (
	"net/http"

	"github.com/chickiexd/zenful_shopping/internal/logger"
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
	err = utils.WriteJSON(w, http.StatusOK, parsed_recipe)
	if err != nil {
		logger.Log.Errorw("Error during response", "error", err)
	}

}

