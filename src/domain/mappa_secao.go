package domain

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type MappaSecao struct {
	gorm.Model
	ID           int
	Nome         string
	CodTipoSecao int
	CodGrupo     int
	CodRegiao    string
	UltimaConsultaMarcacao time.Time
}

var tiposSecao = map[int]string{
	1: "Alcateia",
	2: "Tropa",
	3: "Tropa sênior",
	4: "Clã",
}

func (secao *MappaSecao) ToString() string {

	return fmt.Sprintf("%s %s", tiposSecao[secao.CodTipoSecao], secao.Nome)
}
