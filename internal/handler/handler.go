package handler

import (
	"net/http"
	"zenful_shopping_backend/internal/service"
)

type Handler struct {
	// Recipes interface {
	// 	// Create(context.Context, *CreateRecipeRequest) error
	// 	GetAll(*context.Context) error
	// }
	Ingredients interface {
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
		// Recipes: &RecipeHandler{service},
		Ingredients: &IngredientHandler{service},
		// Users:   &userHandler{service},
		ChatGPT: &ChatGPTHandler{service},
	}
}
