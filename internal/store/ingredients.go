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
	FoodGroups       []FoodGroup       `gorm:"many2many:ingredient_food_groups;joinForeignKey:IngredientID;joinReferences:FoodGroupID"`
	Recipes          []Recipe          `gorm:"many2many:recipe_ingredients;"`
}

type IngredientRepository struct {
	db *gorm.DB
}

func (r *IngredientRepository) WithTransaction(tx *gorm.DB) *IngredientRepository {
	return &IngredientRepository{db: tx}
}

func (r *IngredientRepository) Create(ingredient *Ingredient) error {
	if err := r.db.Create(ingredient).Error; err != nil {
		return err
	}
	return nil
}

func (r *IngredientRepository) CreateFoodGroupAssociation(ingredient *Ingredient, foodGroup *FoodGroup) error {
	if err := r.db.Model(ingredient).Association("FoodGroups").Append(foodGroup); err != nil {
		return err
	}
	return nil
}

func (r *IngredientRepository) CreateMeasurementUnitAssociation(ingredient *Ingredient, measurementUnit *MeasurementUnit) error {
	if err := r.db.Model(ingredient).Association("MeasurementUnits").Append(measurementUnit); err != nil {
		return err
	}
	return nil
}

func (r *IngredientRepository) GetByID(id uint) (*Ingredient, error) {
	var ingredient Ingredient
	err := r.db.First(&ingredient, id).Error
	return &ingredient, err
}

func (r *IngredientRepository) GetByName(name string) (*Ingredient, error) {
	var ingredient Ingredient
	err := r.db.Where("name = ?", name).First(&ingredient).Error
	return &ingredient, err
}

func (r *IngredientRepository) GetAll() ([]Ingredient, error) {
	var ingredients []Ingredient
	err := r.db.
		Preload("FoodGroups").
		Preload("MeasurementUnits").
		Find(&ingredients).Error
	return ingredients, err
}

func (r *IngredientRepository) GetFoodGroupsByID(id uint) ([]FoodGroup, error) {
	var foodGroups []FoodGroup

	if err := r.db.Joins("JOIN ingredient_food_groups ON food_groups.food_group_id = ingredient_food_groups.food_group_id").
		Where("ingredient_food_groups.ingredient_id = ?", id).
		Find(&foodGroups).Error; err != nil {
		return nil, err
	}

	return foodGroups, nil
}

func (r *IngredientRepository) GetMeasurementUnitsByID(id uint) ([]MeasurementUnit, error) {
	var measurementUnits []MeasurementUnit

	if err := r.db.Joins("JOIN ingredient_measurement_units ON measurement_units.measurement_unit_id = ingredient_measurement_units.measurement_unit_id").
		Where("ingredient_measurement_units.ingredient_id = ?", id).
		Find(&measurementUnits).Error; err != nil {
		return nil, err
	}

	return measurementUnits, nil
}
