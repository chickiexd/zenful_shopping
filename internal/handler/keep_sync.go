package handler

import (
	"net/http"

	"github.com/chickiexd/zenful_shopping/internal/errors"
	"github.com/chickiexd/zenful_shopping/internal/service"
	"github.com/chickiexd/zenful_shopping/utils"
)

type KeepSyncHandler struct {
	service *service.Service
}

func (h *KeepSyncHandler) SyncShoppingLists(w http.ResponseWriter, r *http.Request) {
	if err := h.service.KeepSync.SyncShoppingLists(r.Context()); err != nil {
		errors.InternalServerError(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "Shopping lists synced successfully")
}
