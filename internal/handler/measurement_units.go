package handler

import (
	"log"
	"net/http"

	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/chickiexd/zenful_shopping/internal/service"
	"github.com/chickiexd/zenful_shopping/utils"
)

type MeasurementUnitHandler struct {
	service *service.Service
}

func (h *MeasurementUnitHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	log.Println("get all")
	measurements, err := h.service.MeasurmentUnits.GetAll();
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err.Error())
	}
	utils.WriteJSON(w, http.StatusOK, measurements)
}

func (h *MeasurementUnitHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("create")
	var measurement dto.CreateMeasurementUnit
	if err := utils.ReadJSON(w, r, &measurement); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
	}
	h.service.MeasurmentUnits.Create(&measurement)
}
