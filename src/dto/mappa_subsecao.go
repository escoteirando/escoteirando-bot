package dto

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/mappa"
)

func CreateSubsecao(subsecao mappa.SubSecao) domain.MappaSubSecao {
	return domain.MappaSubSecao{
		ID:           subsecao.Codigo,
		Nome:         subsecao.Nome,
		CodSecao:     subsecao.CodigoSecao,
		CodLider:     subsecao.CodigoLider,
		CodViceLider: subsecao.CodigoViceLider,
	}
}
