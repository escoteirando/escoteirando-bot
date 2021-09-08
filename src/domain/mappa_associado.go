package domain

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type MappaAssociado struct {
	gorm.Model
	ID                      int
	Nome                    string
	CodigoFoto              int
	CodigoEquipe            int
	Username                int
	NumeroDigito            int
	DataNascimento          time.Time
	DataValidade            time.Time
	NomeAbreviado           string
	Sexo                    string
	CodigoRamo              int
	CodigoCategoria         int
	CodigoSegundaCategoria  int
	CodigoTerceiraCategoria int
	LinhaFormacao           string
	CodigoRamoAdulto        int
	DataAcompanhamento      time.Time
	Users                   []*User `gorm:"many2many:user_associados;"`
}

var tiposMembro = map[int][]string{
	2: {"Lobo", "Loba"},
	3: {"Escoteiro", "Escoteira"},
	4: {"Sênior", "Guia"},
	5: {"Pioneiro", "Pioneira"},
}

func (associado *MappaAssociado) ToString() string {
	tipoMembro := ""
	if associado.CodigoRamoAdulto > 0 {
		tipoMembro = "Escotista"
	} else {
		tm := tiposMembro[associado.CodigoRamo]
		if len(tm) == 0 {
			tipoMembro = "Associado"
		} else {
			if associado.Sexo == "F" {
				tipoMembro = tm[1]
			} else {
				tipoMembro = tm[0]
			}
		}
	}
	return fmt.Sprintf("%s %s", tipoMembro, associado.Nome)
}
