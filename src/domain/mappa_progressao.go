package domain

import (
	"fmt"
	"gorm.io/gorm"
)

type MappaProgressao struct {
	gorm.Model
	ID                    int `gorm:"primaryKey"`
	Descricao             string
	CodigoUeb             string
	Ordenacao             int
	CodigoCaminho         int
	CodigoDesenvolvimento int
	NumeroGrupo           int
	CodigoRegiao          string
	CodigoCompetencia     int
	Segmento              string
	AreaDesenvolvimento   string
}

func (progressao MappaProgressao) ToString() string {
	return fmt.Sprintf("[%s] %s", progressao.CodigoUeb, progressao.Descricao)
}
