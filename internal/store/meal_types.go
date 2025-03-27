package store

type MealType struct {
    MealTypeID uint   `gorm:"primaryKey"`
    Name       string `gorm:"unique"`
}

