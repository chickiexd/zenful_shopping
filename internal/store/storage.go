package store

import (
	"context"

	"gorm.io/gorm"
)

type Storage struct {
	DB      *gorm.DB
	Recipes interface {
		GetAll() ([]Recipe, error)
		Create(*Recipe) error
	}
	Ingredients interface {
		Create(*Ingredient) error
		GetByID(uint) (*Ingredient, error)
		GetByName(string) (*Ingredient, error)
		GetAll() ([]Ingredient, error)
		GetMeasurementUnitsByID(uint) ([]MeasurementUnit, error)
		GetFoodGroupsByID(uint) ([]FoodGroup, error)
	}
	MeasurementUnits interface {
		WithTransaction(*gorm.DB) *MeasurementRepository
		Create(*MeasurementUnit) error
		GetAll() ([]MeasurementUnit, error)
		GetByID(uint) (*MeasurementUnit, error)
	}
	FoodGroups interface {
		WithTransaction(*gorm.DB) *FoodGroupRepository
		Create(*FoodGroup) error
		GetAll() ([]FoodGroup, error)
		GetByID(uint) (*FoodGroup, error)
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Instructions interface {
		WithTransaction(*gorm.DB) *InstructionRepository
		Create(*Instruction) error
		GetByID(uint) (*Instruction, error)
		GetByRecipeID(uint) ([]Instruction, error)
	}
}

func NewStorage(db *gorm.DB) Storage {
	return Storage{
		DB:               db,
		Recipes:          &RecipeRepository{db},
		Users:            &UserRepository{db},
		Ingredients:      &IngredientRepository{db},
		MeasurementUnits: &MeasurementRepository{db},
		FoodGroups:       &FoodGroupRepository{db},
		Instructions:     &InstructionRepository{db},
	}
}
