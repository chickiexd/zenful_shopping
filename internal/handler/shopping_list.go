package handler

import (
	// "log"
	"net/http"
	"zenful_shopping_backend/internal/dto"
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

func (h *ShoppingListHandler) RemoveItemFromShoppingList(w http.ResponseWriter, r *http.Request) {
	var item_id dto.ShoppingListItemID
	if err := utils.ReadJSON(w, r, &item_id); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	err := h.service.ShoppingList.RemoveItemFromShoppingList(item_id.ItemID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, "Item removed from shopping list")
}

func (h *ShoppingListHandler) RemoveAllItemsFromShoppingList(w http.ResponseWriter, r *http.Request) {
	var shopping_list_id dto.ShoppingListID
	if err := utils.ReadJSON(w, r, &shopping_list_id); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	err := h.service.ShoppingList.RemoveAllItemsFromShoppingList(shopping_list_id.ShoppingListID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, "All items removed from shopping list")
}
