package store

import "gorm.io/gorm"

type Instruction struct {
	InstructionID uint `gorm:"primaryKey"`
	RecipeID      uint
	StepNumber    int
	Description   string

	Recipe Recipe `gorm:"foreignKey:RecipeID"`
}

type InstructionRepository struct {
	db *gorm.DB
}

func (r *InstructionRepository) WithTransaction(tx *gorm.DB) *InstructionRepository {
	return &InstructionRepository{db: tx}
}

func (r *InstructionRepository) Create(instruction *Instruction) error {
	if err := r.db.Create(instruction).Error; err != nil {
		return err
	}
	return nil
}

func (r *InstructionRepository) GetByID(id uint) (*Instruction, error) {
	var instruction Instruction
	err := r.db.First(&instruction, id).Error
	return &instruction, err
}

func (r *InstructionRepository) GetByRecipeID(id uint) ([]Instruction, error) {
	var instructions []Instruction
	err := r.db.Where("recipe_id = ?", id).Find(&instructions).Error
	return instructions, err
}
