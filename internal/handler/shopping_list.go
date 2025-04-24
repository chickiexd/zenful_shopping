package handler


import (
	"net/http"
	"zenful_shopping_backend/internal/service"
	"zenful_shopping_backend/utils"
)

type ShoppingListHandler struct {
	service *service.Service
}

func (h *ShoppingListHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	shopping_lists, err := h.service.ShoppingList.GetAll()
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, shopping_lists)
}
