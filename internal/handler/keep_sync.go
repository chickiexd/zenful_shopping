package handler

import (
	// "log"
	"net/http"
	"github.com/chickiexd/zenful_shopping/internal/service"
	"github.com/chickiexd/zenful_shopping/utils"
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
