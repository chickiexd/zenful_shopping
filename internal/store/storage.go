package store

import (
	"context"

	"gorm.io/gorm"
)

type Storage struct {
	DB      *gorm.DB
	Recipes interface {
		WithTransaction(*gorm.DB) *RecipeRepository
		GetAll() ([]Recipe, error)
		Create(*Recipe) error
		GetByID(uint) (*Recipe, error)
	}
	Ingredients interface {
		WithTransaction(*gorm.DB) *IngredientRepository
		Create(*Ingredient) error
		GetByID(uint) (*Ingredient, error)
		GetByName(string) (*Ingredient, error)
		GetAll() ([]Ingredient, error)
		GetMeasurementUnitsByID(uint) ([]MeasurementUnit, error)
		GetFoodGroupsByID(uint) ([]FoodGroup, error)
	}
	Pantry interface {
		WithTransaction(*gorm.DB) *PantryRepository
		GetAll() ([]PantryIngredient, error)
		Create(*PantryIngredient) (*PantryIngredient, error)
		Delete(*PantryIngredient) error
		DeleteByIngredientID(uint) error
		DeleteAll() error
		GetByIngredientID(uint) (*PantryIngredient, error)
	}
	MeasurementUnits interface {
		WithTransaction(*gorm.DB) *MeasurementRepository
		Create(*MeasurementUnit) error
		GetAll() ([]MeasurementUnit, error)
		GetByID(uint) (*MeasurementUnit, error)
		GetByName(string) (*MeasurementUnit, error)
	}
	FoodGroups interface {
		WithTransaction(*gorm.DB) *FoodGroupRepository
		Create(*FoodGroup) error
		GetAll() ([]FoodGroup, error)
		GetByID(uint) (*FoodGroup, error)
		GetByName(string) (*FoodGroup, error)
	}
	Aliases interface {
		WithTransaction(*gorm.DB) *AliasRepository
		GetIngredientByAlias(string) (*Ingredient, error)
		Upsert(*IngredientAlias) error
	}
	MealTypes interface {
		GetAll() ([]MealType, error)
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
	ShoppingLists interface {
		GetAll() ([]ShoppingList, error)
		WithTransaction(*gorm.DB) *ShoppingListRepository
		Create(*ShoppingList) error
		CreateFoodGroupAssociation(*ShoppingList, *FoodGroup) error
		CreateIngredientAssociation(*ShoppingListItem) error
		GetByFoodGroupID(uint) ([]ShoppingList, error)
		GetItem(uint, uint, uint) (*ShoppingListItem, error)
		GetItemsByIngredientAndUnit(uint, uint) ([]ShoppingListItem, error)
		GetItemsByShoppingListID(uint) ([]ShoppingListItem, error)
		CreateItemAssociation(*ShoppingListItem) error
		DeleteItemAssociation(*ShoppingListItem) error
		UpdateItemAssociation(*ShoppingListItem) error
		DeleteItemAssociationByID(uint) error
		DeleteAllItemsByShoppingListID(uint) error
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
		ShoppingLists:    &ShoppingListRepository{db},
		Pantry:           &PantryRepository{db},
		Aliases:          &AliasRepository{db},
		MealTypes:        &MealTypeRepository{db},
	}
}
