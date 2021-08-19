package domain

import "gorm.io/gorm"

type (
	Config struct {
		gorm.Model
		LastMessageOffset int
		AdminGroup int
	}
)
