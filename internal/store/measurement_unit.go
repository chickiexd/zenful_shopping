package store

import (
	"time"
)

type MeasurementUnit struct {
	MeasurementUnitID uint         `gorm:"primaryKey"`
	Name              string       `gorm:"not null"`
	CreatedAt         time.Time    `gorm:"autoCreateTime"`
	UpdatedAt         time.Time    `gorm:"autoUpdateTime"`
	Ingredients       []Ingredient `gorm:"many2many:ingredient_measurement_units;joinForeignKey:MeasurementUnitID;joinReferences:IngredientID"`
}
