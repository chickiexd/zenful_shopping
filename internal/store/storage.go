package store

import (
	"context"

	"gorm.io/gorm"
)

type Storage struct {
	Recipes interface {
		Create(context.Context, *Recipe) error
	}
	Ingredients interface {
		Create(*Ingredient) error
		GetByID(uint) (*Ingredient, error)
		GetByName(string) (*Ingredient, error)
		GetAll() ([]Ingredient, error)
		GetAllTest() ([]Ingredient, error)
		GetMeasurementUnitsByID(uint) ([]MeasurementUnit, error)
		GetFoodGroupsByID(uint) ([]FoodGroup, error)
	}
	MeasurementUnits interface {
		Create(*MeasurementUnit) error
		GetAll() ([]MeasurementUnit, error)
	}
	FoodGroups interface {
		Create(*FoodGroup) error
		GetAll() ([]FoodGroup, error)
	}
	Users interface {
		Create(context.Context, *User) error
	}
}

func NewStorage(db *gorm.DB) Storage {
	return Storage{
		Recipes:     &RecipeRepository{db},
		Users:       &UserRepository{db},
		Ingredients: &IngredientRepository{db},
		MeasurementUnits: &MeasurementRepository{db},
		FoodGroups: &FoodGroupRepository{db},
	}
}
