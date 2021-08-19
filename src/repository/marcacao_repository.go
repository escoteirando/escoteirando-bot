package repository

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"log"
	"time"
)

func GetMarcacaoCount(codSecao int) (int64, error) {
	var count int64
	response := GetDB().Model(&domain.MappaMarcacao{}).Where("cod_secao = ?", codSecao).Count(&count)
	return count, response.Error
}

func GetMarcacoesAfter(codSecao int, offset int) ([]domain.MappaMarcacao, error) {
	var marcacoes []domain.MappaMarcacao
	response := GetDB().Model(&domain.MappaMarcacao{}).Where("cod_secao = ?", codSecao).Offset(offset).Find(&marcacoes)
	return marcacoes, response.Error
}

func GetNotSentMarcacoes(since time.Time) ([]domain.MappaMarcacao, error) {
	var marcacoes []domain.MappaMarcacao
	response := GetDB().Model(&domain.MappaMarcacao{}).Order("data_atividade").Where("status_escotista = ? and not notificado_chat and data_atividade > ?", "confirmadoEscotista", since).Find(&marcacoes)

	if response.Error == nil {
		cacheProgressoes := NewCachedMappaProgressao()
		cacheAssociados := NewCachedMappaAssociado()
		for i, marcacao := range marcacoes {
			progressao, err := cacheProgressoes.GetProgressao(marcacao.CodigoAtividade)
			if err != nil {
				continue
			}
			marcacoes[i].Progressao = progressao

			associado, err := cacheAssociados.GetAssociado(marcacao.CodigoAssociado)
			if err != nil {
				continue
			}
			marcacoes[i].Associado = associado
		}
	}
	return marcacoes, response.Error
}

func SetMarcacoesAsSent(idMarcacoes []uint) error {
	if len(idMarcacoes) == 0 {
		return nil
	}
	result := GetDB().Model(&domain.MappaMarcacao{}).Where("id in ?", idMarcacoes).UpdateColumn("notificado_chat", true)
	if result.Error == nil {
		log.Printf("Marcacoes sent: %v", idMarcacoes)
	} else {
		log.Printf("Error on saving sent marcacoes: %v", result.Error)
	}
	return result.Error
}
