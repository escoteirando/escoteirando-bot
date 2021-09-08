package dto

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"gorm.io/gorm"
)

func CreateMarcacao(marcacao mappa.Marcacao) domain.MappaMarcacao {
	return domain.MappaMarcacao{
		Model:                 gorm.Model{ID: GetIdMarcacao(marcacao)},
		CodigoAtividade:       marcacao.CodigoAtividade,
		CodigoAssociado:       marcacao.CodigoAssociado,
		DataAtividade:         utils.TimeParse(marcacao.DataAtividade),
		DataStatusJovem:       utils.TimeParse(marcacao.DataStatusJovem),
		DataStatusEscotista:   utils.TimeParse(marcacao.DataStatusEscotista),
		StatusJovem:           marcacao.StatusJovem,
		StatusEscotista:       marcacao.StatusEscotista,
		DataHoraAtualizacao:   utils.TimeParse(marcacao.DataHoraAtualizacao),
		CodigoUltimoEscotista: marcacao.CodigoUltimoEscotista,
		Segmento:              marcacao.Segmento,
		NotificadoChat:        false,
	}
}

func GetIdMarcacao(marcacao mappa.Marcacao) uint {
	return uint(marcacao.CodigoAtividade*1000000 + marcacao.CodigoAssociado)
}
