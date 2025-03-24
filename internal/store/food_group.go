package store

import (
	"time"

	"gorm.io/gorm"
)

type FoodGroup struct {
	FoodGroupID uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	Description string
	CreatedAt   time.Time    `gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime"`
	Ingredients []Ingredient `gorm:"many2many:ingredient_food_groups;joinForeignKey:FoodGroupID;joinReferences:IngredientID"`
}

type FoodGroupRepository struct {
	db *gorm.DB
}

func (r *FoodGroupRepository) Create(food_group *FoodGroup) error {
	if err := r.db.Create(food_group).Error; err != nil {
		return err
	}
	return nil
}

func (r *FoodGroupRepository) GetAll() ([]FoodGroup, error) {
	var food_groups []FoodGroup
	err := r.db.
		Preload("Ingredients").
		Find(&food_groups).Error
	return food_groups, err
}
