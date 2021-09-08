package dto

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/mappa"
)

func CreateProgressao(progressao mappa.Progressao) domain.MappaProgressao {
	return domain.MappaProgressao{
		ID:                    progressao.Codigo,
		Descricao:             progressao.Descricao,
		CodigoUeb:             progressao.CodigoUeb,
		Ordenacao:             progressao.Ordenacao,
		CodigoCaminho:         progressao.CodigoCaminho,
		CodigoDesenvolvimento: progressao.CodigoDesenvolvimento,
		NumeroGrupo:           progressao.NumeroGrupo,
		CodigoRegiao:          progressao.CodigoRegiao,
		CodigoCompetencia:     progressao.CodigoCompetencia,
		Segmento:              progressao.Segmento,
		AreaDesenvolvimento:   progressao.AreaDesenvolvimento,
	}
}
