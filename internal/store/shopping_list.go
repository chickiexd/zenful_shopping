package store

import (
	"time"

	"gorm.io/gorm"
)

type ShoppingList struct {
	ShoppingListID uint `gorm:"primaryKey"`
	Name           string    `gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

type ShoppingListItem struct {
	ShoppingListItemID uint `gorm:"primaryKey"`
	ShoppingListID     uint
	IngredientID       uint
	MeasurementUnitID  uint
	Quantity           float64         `gorm:"not null"`
	ShoppingList       ShoppingList    `gorm:"foreignKey:ShoppingListID"`
	Ingredient         Ingredient      `gorm:"foreignKey:IngredientID"`
	MeasurementUnit    MeasurementUnit `gorm:"foreignKey:MeasurementUnitID"`
	CreatedAt          time.Time       `gorm:"autoCreateTime"`
	UpdatedAt          time.Time       `gorm:"autoUpdateTime"`
}

type ShoppingListFoodGroup struct {
	ShoppingListID uint         `gorm:"primaryKey"`
	FoodGroupID    uint         `gorm:"primaryKey"`
	ShoppingList   ShoppingList `gorm:"foreignKey:ShoppingListID"`
	FoodGroup      FoodGroup    `gorm:"foreignKey:FoodGroupID"`
	CreatedAt      time.Time    `gorm:"autoCreateTime"`
	UpdatedAt      time.Time    `gorm:"autoUpdateTime"`
}

type ShoppingListRepository struct {
	db *gorm.DB
}

func (r *ShoppingListRepository) GetAll() ([]ShoppingList, error) {
	var shoppingLists []ShoppingList
	err := r.db.
		Find(&shoppingLists).Error
	return shoppingLists, err
}

func (r *ShoppingListRepository) WithTransaction(tx *gorm.DB) *ShoppingListRepository {
	return &ShoppingListRepository{db: tx}
}

func (r *ShoppingListRepository) Create(shopping_list *ShoppingList) error {
	if err := r.db.Create(shopping_list).Error; err != nil {
		return err
	}
	return nil
}

func (r *ShoppingListRepository) CreateFoodGroupAssociation(shopping_list *ShoppingList, foodGroup *FoodGroup) error {
	if err := r.db.Model(shopping_list).Association("FoodGroups").Append(foodGroup); err != nil {
		return err
	}
	return nil
}

func (r *ShoppingListRepository) CreateIngredientAssociation(shopping_list_item *ShoppingListItem) error {
	if err := r.db.Create(shopping_list_item).Error; err != nil {
		return err
	}
	return nil
}

func (r *ShoppingListRepository) GetByFoodGroupID(id uint) ([]ShoppingList, error) {
	var shoppingLists []ShoppingList
	err := r.db.Joins("JOIN shopping_list_food_groups ON shopping_lists.shopping_list_id = shopping_list_food_groups.shopping_list_id").
		Where("shopping_list_food_groups.food_group_id = ?", id).
		Find(&shoppingLists).Error
	return shoppingLists, err
}

func (r *ShoppingListRepository) GetItemsByIngredientID(id uint) ([]ShoppingListItem, error) {
	var items []ShoppingListItem
	err := r.db.
		Preload("ShoppingList").
		Preload("Ingredient").
		Preload("MeasurementUnit").
		Where("ingredient_id = ?", id).
		Find(&items).Error
	return items, err
}

func (r *ShoppingListRepository) GetItemsByShoppingListID(id uint) ([]ShoppingListItem, error) {
	var items []ShoppingListItem
	err := r.db.
		Preload("ShoppingList").
		Preload("Ingredient").
		Preload("MeasurementUnit").
		Where("shopping_list_id = ?", id).
		Find(&items).Error
	return items, err
}

func (r *ShoppingListRepository) CreateItemAssociation(shopping_list_item *ShoppingListItem) error {
	if err := r.db.Create(shopping_list_item).Error; err != nil {
		return err
	}
	return nil
}

func (r *ShoppingListRepository) DeleteItemAssociation(shopping_list_item *ShoppingListItem) error {
	if err := r.db.Delete(shopping_list_item).Error; err != nil {
		return err
	}
	return nil
}

func (r *ShoppingListRepository) UpdateItemAssociation(shopping_list_item *ShoppingListItem) error {
	if err := r.db.Save(shopping_list_item).Error; err != nil {
		return err
	}
	return nil
}

// func (r *ShoppingListRepository) GetByID(id uint) (*ShoppingList, error) {
// 	var ingredient ShoppingList
// 	err := r.db.First(&ingredient, id).Error
// 	return &ingredient, err
// }
//
// func (r *ShoppingListRepository) GetByName(name string) (*ShoppingList, error) {
// 	var ingredient ShoppingList
// 	err := r.db.Where("name = ?", name).First(&ingredient).Error
// 	return &ingredient, err
// }
//
// func (r *ShoppingListRepository) GetAll() ([]ShoppingList, error) {
// 	var ingredients []ShoppingList
// 	err := r.db.
// 		Preload("FoodGroups").
// 		Preload("MeasurementUnits").
// 		Find(&ingredients).Error
// 	return ingredients, err
// }
//
// func (r *ShoppingListRepository) GetFoodGroupsByID(id uint) ([]FoodGroup, error) {
// 	var foodGroups []FoodGroup
//
// 	if err := r.db.Joins("JOIN ingredient_food_groups ON food_groups.food_group_id = ingredient_food_groups.food_group_id").
// 		Where("ingredient_food_groups.ingredient_id = ?", id).
// 		Find(&foodGroups).Error; err != nil {
// 		return nil, err
// 	}
//
// 	return foodGroups, nil
// }
//
// func (r *ShoppingListRepository) GetMeasurementUnitsByID(id uint) ([]MeasurementUnit, error) {
// 	var measurementUnits []MeasurementUnit
//
// 	if err := r.db.Joins("JOIN ingredient_measurement_units ON measurement_units.measurement_unit_id = ingredient_measurement_units.measurement_unit_id").
// 		Where("ingredient_measurement_units.ingredient_id = ?", id).
// 		Find(&measurementUnits).Error; err != nil {
// 		return nil, err
// 	}
//
// 	return measurementUnits, nil
// }
