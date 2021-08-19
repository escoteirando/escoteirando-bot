package domain

import "gorm.io/gorm"

type MappaSubSecao struct {
	gorm.Model
	ID           int
	Nome         string
	CodSecao     int
	CodLider     int
	CodViceLider int
}
