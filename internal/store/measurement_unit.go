package store

import (
	"time"

	"gorm.io/gorm"
)

type MeasurementUnit struct {
	MeasurementUnitID uint         `gorm:"primaryKey" json:"measurement_unit_id"`
	Name              string       `gorm:"not null" json:"name"`
	Abbreviation      string       `json:"abbreviation"`
	CreatedAt         time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	Ingredients       []Ingredient `gorm:"many2many:ingredient_measurement_units;joinForeignKey:MeasurementUnitID;joinReferences:IngredientID" json:"ingredients"`
}

type MeasurementRepository struct {
	db *gorm.DB
}

func (r *MeasurementRepository) WithTransaction(tx *gorm.DB) *MeasurementRepository {
	return &MeasurementRepository{db: tx}
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

func (r *MeasurementRepository) GetByID(id uint) (*MeasurementUnit, error) {
	var measurement MeasurementUnit
	err := r.db.First(&measurement, id).Error
	return &measurement, err
}

func (r *MeasurementRepository) GetByName(name string) (*MeasurementUnit, error) {
	var measurement MeasurementUnit
	err := r.db.Where("name = ?", name).First(&measurement).Error
	return &measurement, err
}
