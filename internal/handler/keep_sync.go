package handler

import (
	// "log"
	"net/http"
	"zenful_shopping_backend/internal/service"
	"zenful_shopping_backend/utils"
)

type KeepSyncHandler struct {
	service *service.Service
}

func (h *KeepSyncHandler) SyncShoppingLists(w http.ResponseWriter, r *http.Request) {

	err := h.service.KeepSync.SyncShoppingLists()
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, "Shopping lists synced successfully")
}
