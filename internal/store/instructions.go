package store

type Instruction struct {
	InstructionID uint   `gorm:"primaryKey"`
	RecipeID      uint
	StepNumber    int
	Description   string

	Recipe Recipe `gorm:"foreignKey:RecipeID"`
}

