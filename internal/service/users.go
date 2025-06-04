package service

import (
	"context"

	"github.com/chickiexd/zenful_shopping/internal/store"
)

type userService struct {
	storage *store.Storage
}

func (s *userService) Create(ctx context.Context, user *store.User) error {
	return s.storage.Users.Create(ctx, user)
}
