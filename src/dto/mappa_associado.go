package dto

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"github.com/guionardo/escoteirando-bot/src/utils"
)

func CreateAssociado(associado mappa.Associado) domain.MappaAssociado{
	return domain.MappaAssociado{
		ID:                      associado.Codigo,
		Nome:                    associado.Nome,
		CodigoFoto:              associado.CodigoFoto,
		CodigoEquipe:            associado.CodigoEquipe,
		Username:                associado.Username,
		NumeroDigito:            associado.NumeroDigito,
		DataNascimento:          utils.TimeParse(associado.DataNascimento),
		DataValidade:            utils.TimeParse(associado.DataValidade),
		NomeAbreviado:           associado.NomeAbreviado,
		Sexo:                    associado.Sexo,
		CodigoRamo:              associado.CodigoRamo,
		CodigoCategoria:         associado.CodigoCategoria,
		CodigoSegundaCategoria:  associado.CodigoSegundaCategoria,
		CodigoTerceiraCategoria: associado.CodigoTerceiraCategoria,
		LinhaFormacao:           associado.LinhaFormacao,
		CodigoRamoAdulto:        associado.CodigoRamoAdulto,
		DataAcompanhamento:      utils.TimeParse(associado.DataAcompanhamento),
		Users:                   nil,
	}
}
