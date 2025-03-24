package store

import (
	"time"

	"gorm.io/gorm"
)

type MeasurementUnit struct {
	MeasurementUnitID uint         `gorm:"primaryKey"`
	Name              string       `gorm:"not null"`
	CreatedAt         time.Time    `gorm:"autoCreateTime"`
	UpdatedAt         time.Time    `gorm:"autoUpdateTime"`
	Ingredients       []Ingredient `gorm:"many2many:ingredient_measurement_units;joinForeignKey:MeasurementUnitID;joinReferences:IngredientID"`
}

type MeasurementRepository struct {
	db *gorm.DB
}

func (r *MeasurementRepository) Create(measurement_unit *MeasurementUnit) error {
	if err := r.db.Create(measurement_unit).Error; err != nil {
		return err
	}
	return nil
}

func (r *MeasurementRepository) GetAll() ([]MeasurementUnit, error) {
	var measurements []MeasurementUnit
	err := r.db.
		Preload("Ingredients").
		Find(&measurements).Error
	return measurements, err
}
