package handler

import (
	"log"
	"net/http"

	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/service"
	"zenful_shopping_backend/utils"
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
