package store

import (
	"time"

	"gorm.io/gorm"
)

type FoodGroup struct {
	FoodGroupID uint         `gorm:"primaryKey" json:"food_group_id"`
	Name        string       `gorm:"unique;not null" json:"name"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	Ingredients []Ingredient `gorm:"many2many:ingredient_food_groups;joinForeignKey:FoodGroupID;joinReferences:IngredientID" json:"ingredients"`
}

type FoodGroupRepository struct {
	db *gorm.DB
}

func (r *FoodGroupRepository) WithTransaction(tx *gorm.DB) *FoodGroupRepository {
	return &FoodGroupRepository{db: tx}
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

func (r *FoodGroupRepository) GetByID(id uint) (*FoodGroup, error) {
	var food_group FoodGroup
	err := r.db.First(&food_group, id).Error
	return &food_group, err
}
