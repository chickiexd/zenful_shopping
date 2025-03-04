package store

import (
	"time"

	"gorm.io/gorm"
)

type Ingredient struct {
	IngredientID     uint   `gorm:"primaryKey"`
	Name             string `gorm:"unique;not null"`
	Description      string
	CreatedAt        time.Time         `gorm:"autoCreateTime"`
	UpdatedAt        time.Time         `gorm:"autoUpdateTime"`
	MeasurementUnits []MeasurementUnit `gorm:"many2many:ingredient_measurement_units;joinForeignKey:IngredientID;joinReferences:MeasurementUnitID"`
	FoodGroups []FoodGroup `gorm:"many2many:ingredient_food_groups;joinForeignKey:IngredientID;joinReferences:FoodGroupID"`
}

type IngredientRepository struct {
	db *gorm.DB
}

// func (r *IngredientRepository) GetOrInsertIngredient(ingredient *Ingredient, tx *gorm.DB) (*Ingredient, error) {
// 	var existingIngredient Ingredient
// 	if err := tx.Where("name = ?", ingredient.Name).First(&existingIngredient).Error; err == nil {
// 		return &existingIngredient, nil
// 	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
// 		return nil, err
// 	}
// 	return r.Create(ingredient, tx)
// }

func (r *IngredientRepository) Create(ingredient *Ingredient) error {
	// if err := db.Create(ingredient).Error; err != nil {
	// 	return err
	// }
	// return nil
	return nil
}

func (r *IngredientRepository) GetByID(id int64) (*Ingredient, error) {
	var ingredient Ingredient
	err := r.db.First(&ingredient, id).Error
	return &ingredient, err
}

func (r *IngredientRepository) GetAll() ([]Ingredient, error) {
	var ingredients []Ingredient
	err := r.db.Find(&ingredients).Error
	return ingredients, err
}


