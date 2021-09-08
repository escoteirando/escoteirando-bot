package repository

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/dto"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"gorm.io/gorm/clause"
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

func GetAllEscotistasIds() ([]int, error) {
	var results []domain.MappaMarcacao
	result := GetDB().
		Distinct("codigo_ultimo_escotista").
		Select("codigo_ultimo_escotista").
		Where("codigo_ultimo_escotista > 0").
		Model(&domain.MappaMarcacao{}).
		Find(&results)

	if result.Error != nil {
		return nil, result.Error
	}
	response := make([]int, result.RowsAffected)
	for i, cod := range results {
		response[i] = cod.CodigoUltimoEscotista
	}
	return response, nil
}

func GetUltimaConsultaMarcacao(codSecao int) time.Time {
	var secao domain.MappaSecao
	response := GetDB().Find(&secao, codSecao)
	if response.Error != nil {
		return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return secao.UltimaConsultaMarcacao
}

func SetUltimaConsultaMarcacao(codSecao int, quando time.Time) error {
	var secao domain.MappaSecao
	response := GetDB().Find(&secao, codSecao)
	if response.Error != nil {
		return response.Error
	}
	CanIWrite()
	secao.UltimaConsultaMarcacao = quando
	response = GetDB().Save(&secao)
	YouCanWrite()
	return response.Error
}

func SetMarcacoes(marcacoes []mappa.Marcacao) error {
	dbMarcacoes := make([]domain.MappaMarcacao, len(marcacoes))
	for i, marcacao := range marcacoes {
		dbMarcacoes[i]=dto.CreateMarcacao(marcacao)
	}
	CanIWrite()
	result := GetDB().
		Clauses(clause.OnConflict{UpdateAll: true}).
		CreateInBatches(dbMarcacoes,10)
	YouCanWrite()
	return result.Error
}