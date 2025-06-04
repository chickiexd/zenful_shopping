package service

import (
	"log"
	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/chickiexd/zenful_shopping/internal/store"
)

type MeasurementUnitService struct {
	storage *store.Storage
}

func (s *MeasurementUnitService) Create(ingredient *dto.CreateMeasurementUnit) error {
	log.Println("service create ingredient_id")
	return s.storage.MeasurementUnits.Create(nil)
}

func (s *MeasurementUnitService) GetAll() ([]store.MeasurementUnit, error) {
	measurement_units, err := s.storage.MeasurementUnits.GetAll()
	return measurement_units, err
}

// func (s *MeasurementUnitService) GetByID(id uint) (*store.Measurement, error) {
// 	ingredient, err := s.storage.MeasurementUnits.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ingredient, nil
// }
