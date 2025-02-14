package handler

import (
	"net/http"

	"zenful_shopping_backend/internal/service"
)

type Handler struct {
	Recipes RecipeHandler
	Users   UserHandler
}

func NewHandler(service service.Service) Handler {
	return Handler{
		Recipes: &recipeHandler{service},
		Users:   &userHandler{service},
	}
}

