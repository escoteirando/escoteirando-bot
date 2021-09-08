package dto

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/mappa"
)

func CreateGrupo(grupo mappa.Grupo) domain.MappaGrupo {
	return domain.MappaGrupo{
		ID:               grupo.Codigo,
		CodigoRegiao:     grupo.CodigoRegiao,
		Nome:             grupo.Nome,
		CodigoModalidade: grupo.CodigoModalidade,
	}
}
