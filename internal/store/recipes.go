package store

import (
	"time"
	"gorm.io/gorm"
)

type Recipe struct {
	RecipeID     uint `gorm:"primaryKey"`
	Title        string
	Description  string
	Public       bool
	CookTime     int
	Servings     int
	Image        string
	UserID       uint
	MealTypeID   uint
	User         User          `gorm:"foreignKey:UserID"`
	MealType     MealType      `gorm:"foreignKey:MealTypeID"`
	Ingredients  []Ingredient  `gorm:"many2many:recipe_has_ingredients;"`
	Instructions []Instruction `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE;"`
	CreatedAt    time.Time     `gorm:"autoCreateTime"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime"`
}

type RecipeHasIngredient struct {
	RecipeID     uint `gorm:"primaryKey"`
	IngredientID uint `gorm:"primaryKey"`
	Quantity     float64
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID"`
}

type RecipeRepository struct {
	db *gorm.DB
}

func (r *RecipeRepository) Create(recipe *Recipe) error {
	if err := r.db.Create(recipe).Error; err != nil {
		return err
	}
	return nil
}

func (r *RecipeRepository) GetAll() ([]Recipe, error) {
	var recipes []Recipe
	err := r.db.
		Preload("Ingredients").
		Preload("Instructions").
		Preload("User").
		Preload("MealType").
		Find(&recipes).Error
	return recipes, err
}
