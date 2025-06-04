package handler

import (
	"log"
	"net/http"

	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/service"
	"zenful_shopping_backend/utils"
)

type FoodGroupHandler struct {
	service *service.Service
}

func (h *FoodGroupHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	food_groups, err := h.service.FoodGroups.GetAll()
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err.Error())
	}
	utils.WriteJSON(w, http.StatusOK, food_groups)
}

func (h *FoodGroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("create")
	var food_group dto.CreateFoodGroupRequest
	if err := utils.ReadJSON(w, r, &food_group); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
	}
	log.Println("food_group: ", food_group)
	created_food_group, err := h.service.FoodGroups.Create(&food_group)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, created_food_group)
}
