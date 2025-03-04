package store

import ( 
	"time" 
)

type FoodGroup struct {
	FoodGroupID uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	Description string
	CreatedAt   time.Time    `gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime"`
	Ingredients []Ingredient `gorm:"many2many:ingredient_food_groups;joinForeignKey:FoodGroupID;joinReferences:IngredientID"`
}
