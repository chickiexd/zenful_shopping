package handler

import (
	"net/http"
	"zenful_shopping_backend/internal/service"
)

type Handler struct {
	Recipes interface {
		Create(http.ResponseWriter, *http.Request)
		GetAll(http.ResponseWriter, *http.Request)
	}
	Ingredients interface {
		GetAll(http.ResponseWriter, *http.Request)
		Create(http.ResponseWriter, *http.Request)
	}
	MeasurementUnits interface {
		GetAll(http.ResponseWriter, *http.Request)
		Create(http.ResponseWriter, *http.Request)
	}
	FoodGroups interface {
		// Create(context.Context, *CreateRecipeRequest) error
		GetAll(http.ResponseWriter, *http.Request)
		Create(http.ResponseWriter, *http.Request)
	}
	ChatGPT interface {
		ParseRecipe(http.ResponseWriter, *http.Request)
	}
}

func NewHandler(service *service.Service) Handler {
	return Handler{
		Recipes: &RecipeHandler{service},
		Ingredients: &IngredientHandler{service},
		MeasurementUnits: &MeasurementUnitHandler{service},
		FoodGroups: &FoodGroupHandler{service},
		// Users:   &userHandler{service},
		ChatGPT: &ChatGPTHandler{service},
	}
}
