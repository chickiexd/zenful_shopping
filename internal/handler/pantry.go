package handler

import (
	"net/http"

	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/chickiexd/zenful_shopping/internal/errors"
	"github.com/chickiexd/zenful_shopping/internal/service"
	"github.com/chickiexd/zenful_shopping/utils"
	"gorm.io/gorm"
)

type PantryHandler struct {
	service *service.Service
}

// GetAll godoc
//
//	@Summary		Get all pantry_ingredients
//	@Description	Get all pantry_ingredients from the database
//	@Tags			pantry_ingredients
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		dto.PantryIngredientResponse
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/pantry_ingredients [get]
func (h *PantryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	pantry_ingredients, err := h.service.Pantry.GetAll()
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, pantry_ingredients)
}

// Add godoc
//
//	@Summary		Add ingredient to pantry
//	@Description	Add an ingredient to the pantry by its ID
//	@Tags			pantry_ingredients
//	@Accept			json
//	@Produce		json
//	@Param			ingredient_id	body		dto.PantryIngredientRequest	true	"Ingredient ID to add"
//	@Success		200				{object}	uint
//	@Failure		400				{object}	error
//	@Failure		404				{object}	error
//	@Failure		500				{object}	error
//	@Router			/pantry_ingredients/add [post]
func (h *PantryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request dto.PantryIngredientRequest
	if err := utils.ReadJSON(w, r, &request); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Pantry.Create(request.IngredientID); err != nil {
		errors.InternalServerError(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, request.IngredientID)
}

// Remove godoc
//
//	@Summary		Remove a ingredient from the pantry
//	@Description	Remove a ingredient from the pantry by its ID
//	@Tags			pantry_ingredients
//	@Accept			json
//	@Produce		json
//	@Param			ingredient_id	body		dto.PantryIngredientRequest	true	"Ingredient ID to remove"
//	@Success		200				{object}	uint
//	@Failure		400				{object}	error
//	@Failure		404				{object}	error
//	@Failure		500				{object}	error
//	@Router			/pantry_ingredients/remove [post]
func (h *PantryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var request dto.PantryIngredientRequest
	if err := utils.ReadJSON(w, r, &request); err != nil {
		errors.BadRequest(w, r, err)
		return
	}
	if err := h.service.Pantry.Delete(request.IngredientID); err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.NotFound(w, r)
			return
		}
		errors.InternalServerError(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, request.IngredientID)
}

// RemoveAll godoc
//
//	@Summary		Remove all ingredients from the pantry
//	@Description	Remove all ingredients from the pantry
//	@Tags			pantry_ingredients
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/pantry_ingredients/remove_all [post]
func (h *PantryHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Pantry.DeleteAll(); err != nil {
		errors.InternalServerError(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "All ingredients removed from pantry")
}
