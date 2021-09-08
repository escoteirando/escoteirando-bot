package dto

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/mappa"
)

func CreateSecao(secao mappa.Secao) domain.MappaSecao {
	return domain.MappaSecao{
		ID:           secao.Codigo,
		Nome:         secao.Nome,
		CodTipoSecao: secao.CodigoTipoSecao,
		CodGrupo:     secao.CodigoGrupo,
		CodRegiao:    secao.CodigoRegiao,
	}
}
