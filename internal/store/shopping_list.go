package store

import (
	"time"

	"gorm.io/gorm"
)

type ShoppingList struct {
	ShoppingListID uint      `gorm:"primaryKey"`
	Name           string    `gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	FoodGroups     []FoodGroup `gorm:"many2many:shopping_list_food_groups;joinForeignKey:ShoppingListID;joinReferences:FoodGroupID"`
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

func (r *ShoppingListRepository) GetByID(id uint) (*ShoppingList, error) {
	var shoppingList ShoppingList
	err := r.db.First(&shoppingList, id).Error
	if err != nil {
		return nil, err
	}
	return &shoppingList, nil
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

func (r *ShoppingListRepository) DeleteItemAssociationByID(id uint) error {
	if err := r.db.Delete(&ShoppingListItem{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *ShoppingListRepository) DeleteAllItemsByShoppingListID(id uint) error {
	if err := r.db.Where("shopping_list_id = ?", id).Delete(&ShoppingListItem{}).Error; err != nil {
		return err
	}
	return nil
}
