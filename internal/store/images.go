package store

import "time"

type Image struct {
	ImageID    uint      `gorm:"primaryKey"`
	FilePath   string    `gorm:"not null"`
	UploadedAt time.Time `gorm:"autoCreateTime"`
}


