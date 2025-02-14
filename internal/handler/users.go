package handler

import (
	"context"
	"net/http"

	"zenful_shopping_backend/internal/service"
	"zenful_shopping_backend/internal/store"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Create(c *gin.Context)
}

type userHandler struct {
	service service.Service
}

func (h *userHandler) Create(c *gin.Context) {
	var user store.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.service.Users.Create(context.Background(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
}

