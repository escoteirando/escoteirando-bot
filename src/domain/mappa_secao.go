package domain

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/consts"
	"gorm.io/gorm"
	"time"
)

type MappaSecao struct {
	gorm.Model
	ID                       int
	Nome                     string
	CodTipoSecao             int
	CodGrupo                 int
	CodRegiao                string
	UltimaConsultaMarcacao   time.Time
	UltimaConsultaConquistas time.Time
}

func (secao *MappaSecao) ToString() string {
	return fmt.Sprintf("%s %s", consts.TipoSecao(secao.CodTipoSecao), secao.Nome)
}
