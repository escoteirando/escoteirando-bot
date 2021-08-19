package domain

import (
	"gorm.io/gorm"
	"time"
)

type MappaMarcacao struct {
	gorm.Model
	CodigoAtividade       int `gorm:"primaryKey"`
	CodigoAssociado       int `gorm:"primaryKey"`
	DataAtividade         time.Time
	DataStatusJovem       time.Time
	DataStatusEscotista   time.Time
	StatusJovem           string
	StatusEscotista       string
	DataHoraAtualizacao   time.Time
	CodigoUltimoEscotista int
	Segmento              string
	NotificadoChat        bool
	Associado             MappaAssociado  `gorm:"foreignKey:CodigoAssociado"`
	Progressao            MappaProgressao `gorm:"foreignKey:CodigoAtividade"`
}
