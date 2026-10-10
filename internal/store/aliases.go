package store

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IngredientAlias struct {
	Alias        string    `gorm:"primaryKey"`
	IngredientID uint      `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

type AliasRepository struct {
	db *gorm.DB
}

func (r *AliasRepository) WithTransaction(tx *gorm.DB) *AliasRepository {
	return &AliasRepository{db: tx}
}

// GetIngredientByAlias expects an already normalized alias.
func (r *AliasRepository) GetIngredientByAlias(alias string) (*Ingredient, error) {
	var ingredient Ingredient
	err := r.db.
		Joins("JOIN ingredient_aliases ON ingredient_aliases.ingredient_id = ingredients.ingredient_id").
		Where("ingredient_aliases.alias = ?", alias).
		First(&ingredient).Error
	if err != nil {
		return nil, err
	}
	return &ingredient, nil
}

// Upsert points alias at ingredient_id, replacing an earlier mapping.
func (r *AliasRepository) Upsert(alias *IngredientAlias) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "alias"}},
		DoUpdates: clause.AssignmentColumns([]string{"ingredient_id"}),
	}).Create(alias).Error
}
