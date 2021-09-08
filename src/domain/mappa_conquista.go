package domain

import (
	"gorm.io/gorm"
	"time"
)

type MappaConquista struct {
	gorm.Model
	Type                           string
	DataConquista                  time.Time
	CodigoAssociado                int
	CodigoEscotistaUltimaAlteracao int
	NumeroNivel                    int
	CodigoEspecialidade            int
	CodigoSecao                    int
	NotificadoChat                 bool
	Associado                      MappaAssociado `gorm:"foreignKey:CodigoAssociado"`
}

func (conquista *MappaConquista) ToString() string{
	return "Conquista" // TODO: Implementar busca quando for especialidade
}