package store

import (
	"gorm.io/gorm"
)

type PantryIngredient struct {
	PantryIngredientID uint `gorm:"primaryKey"`
	IngredientID       uint
	Ingredient         Ingredient `gorm:"foreignKey:IngredientID"`
}

type PantryRepository struct {
	db *gorm.DB
}

func (r *PantryRepository) GetAll() ([]PantryIngredient, error) {
	var pantry_ingredients []PantryIngredient
	err := r.db.
		Preload("Ingredient").
		Find(&pantry_ingredients).Error
	return pantry_ingredients, err
}

func (r *PantryRepository) Create(pantry_ingredient *PantryIngredient) (*PantryIngredient, error) {
	if err := r.db.Create(&pantry_ingredient).Error; err != nil {
		return nil, err
	}
	return pantry_ingredient, nil
}

func (r *PantryRepository) Delete(pantry_ingredient *PantryIngredient) error {
	if err := r.db.Delete(&pantry_ingredient).Error; err != nil {
		return err
	}
	return nil
}

func (r *PantryRepository) DeleteAll() error {
	if err := r.db.Delete(&PantryIngredient{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *PantryRepository) DeleteByIngredientID(id uint) error {
	if err := r.db.Where("ingredient_id = ?", id).Delete(&PantryIngredient{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *PantryRepository) GetByIngredientID(id uint) (*PantryIngredient, error) {
	var pantry_ingredient PantryIngredient
	if err := r.db.Where("ingredient_id = ?", id).First(&pantry_ingredient).Error; err != nil {
		return nil, err
	}
	return &pantry_ingredient, nil
}
