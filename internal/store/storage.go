package store

import (
	"context"

	"gorm.io/gorm"
)

type Storage struct {
	Recipes interface {
		Create(context.Context, *Recipe) error
	}
	Ingredients interface {
		Create(*Ingredient) error
		GetByID(int64) (*Ingredient, error)
		GetAll() ([]Ingredient, error)
	}
	Users interface {
		Create(context.Context, *User) error
	}
}

func NewStorage(db *gorm.DB) Storage {
	return Storage{
		Recipes: &RecipeRepository{db},
		Users:   &UserRepository{db},
		Ingredients:   &IngredientRepository{db},
	}
}
