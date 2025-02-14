package handler

import (
	"context"
	"net/http"

	"zenful_shopping_backend/internal/service"
	"zenful_shopping_backend/internal/store"
	"github.com/gin-gonic/gin"
)

type RecipeHandler interface {
	Create(c *gin.Context)
}

type recipeHandler struct {
	service service.Service
}

func (h *recipeHandler) Create(c *gin.Context) {
	var recipe store.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.service.Recipes.Create(context.Background(), &recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Recipe created successfully", "recipe": recipe})
}

