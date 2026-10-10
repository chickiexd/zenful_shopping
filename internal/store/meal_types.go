package store

import "gorm.io/gorm"

type MealType struct {
	MealTypeID uint   `gorm:"primaryKey"`
	Name       string `gorm:"unique"`
}

type MealTypeRepository struct {
	db *gorm.DB
}

func (r *MealTypeRepository) GetAll() ([]MealType, error) {
	var meal_types []MealType
	err := r.db.Order("meal_type_id").Find(&meal_types).Error
	return meal_types, err
}
