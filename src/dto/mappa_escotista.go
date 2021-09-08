package dto

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"strings"
)

func CreateEscotista(escotista mappa.Escotista) domain.MappaEscotista {
	return domain.MappaEscotista{
		ID:           escotista.Codigo,
		CodAssociado: escotista.CodigoAssociado,
		UserName:     escotista.Username,
		NomeCompleto: escotista.NomeCompleto,
		Ativo:        strings.ContainsAny(escotista.Ativo, "1SsYyTt"),
		CodGrupo:     escotista.CodigoGrupo,
		CodRegiao:    escotista.CodigoRegiao,
		CodFoto:      escotista.CodigoFoto,
	}
}
