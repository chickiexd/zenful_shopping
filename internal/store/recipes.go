package store

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Recipe struct {
	RecipeID    uint `gorm:"primaryKey"`
	Title       string
	UserID      uint
	User        User
	Public      bool
	Description string
	CookTime    int
	Servings    int
	ImageID     uint
	Image       Image
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type RecipeRepository struct {
	db *gorm.DB
}

// func NewRecipeRepository(db *gorm.DB) *RecipeRepository {
// 	return &RecipeRepository{db: db}
// }

func (r *RecipeRepository) Create(ctx context.Context, recipe *Recipe) error {
	if err := r.db.WithContext(ctx).Create(recipe).Error; err != nil {
		return err
	}
	return nil
}
