package domain

import "gorm.io/gorm"

type MappaEspecialidade struct {
	gorm.Model
	Descricao        string
	RamoConhecimento string
	Prerequisito     string
	Items            []MappaEspecialidadeItem `gorm:"foreignKey:CodigoEspecialidade"`
}

type MappaEspecialidadeItem struct {
	gorm.Model
	Id                  int
	CodigoEspecialidade int
	Descricao           string
	Numero              int
}
