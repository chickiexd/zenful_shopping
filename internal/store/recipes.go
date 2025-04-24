package store

import (
	"gorm.io/gorm"
	"time"
)

type Recipe struct {
	RecipeID          uint `gorm:"primaryKey"`
	Title             string
	Description       string
	Public            bool
	CookTime          int
	Servings          int
	ImagePath         string
	MealTypeID        uint
	MealType          MealType           `gorm:"foreignKey:MealTypeID"`
	RecipeIngredients []RecipeIngredient `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE;"`
	Ingredients       []Ingredient       `gorm:"many2many:recipe_ingredients;joinForeignKey:RecipeID;joinReferences:IngredientID"`
	Instructions      []Instruction      `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE;"`
	CreatedAt         time.Time          `gorm:"autoCreateTime"`
	UpdatedAt         time.Time          `gorm:"autoUpdateTime"`
}

type RecipeIngredient struct {
	RecipeID          uint `gorm:"primaryKey"`
	IngredientID      uint `gorm:"primaryKey"`
	MeasurementUnitID uint
	Quantity          float64
	Ingredient        Ingredient      `gorm:"foreignKey:IngredientID"`
	MeasurementUnit   MeasurementUnit `gorm:"foreignKey:MeasurementUnitID"`
}

type RecipeRepository struct {
	db *gorm.DB
}

func (r *RecipeRepository) WithTransaction(tx *gorm.DB) *RecipeRepository {
	return &RecipeRepository{db: tx}
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
		Preload("RecipeIngredients.Ingredient").
		Preload("RecipeIngredients.MeasurementUnit").
		Preload("Instructions").
		Preload("MealType").
		Find(&recipes).Error
	return recipes, err
}

func (r *RecipeRepository) CreateIngredientAssociation(recipeIngredient *RecipeIngredient) error {
	if err := r.db.Create(recipeIngredient).Error; err != nil {
		return err
	}
	return nil
}

func (r *RecipeRepository) GetByID(id uint) (*Recipe, error) {
	var recipe Recipe
	err := r.db.
		Preload("RecipeIngredients.Ingredient").
		Preload("RecipeIngredients.MeasurementUnit").
		Preload("Instructions").
		Preload("MealType").
		First(&recipe, id).Error
	return &recipe, err
}
