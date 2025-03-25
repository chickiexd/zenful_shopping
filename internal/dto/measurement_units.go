package dto

type CreateMeasurementUnit struct {
	MeasurementUnitID uint   `json:"measurement_unit_id"`
	Name string `json:"name"`
}
