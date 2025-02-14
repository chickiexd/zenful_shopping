package service

import (
	"context"

	"zenful_shopping_backend/internal/store"
)

type userService struct {
	storage store.Storage
}

func (s *userService) Create(ctx context.Context, user *store.User) error {
	return s.storage.Users.Create(ctx, user)
}
