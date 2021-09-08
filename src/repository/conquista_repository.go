package repository

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/dto"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

func GetNotSentConquistas(since time.Time) ([]domain.MappaConquista, error) {
	var conquistas []domain.MappaConquista
	response := GetDB().Model(&domain.MappaConquista{}).Order("data_conquista").Where("not notificado_chat and data_conquista > ?", since).Find(&conquistas)

	if response.Error == nil {
		cacheAssociados := NewCachedMappaAssociado()
		for i, conquista := range conquistas {
			associado, err := cacheAssociados.GetAssociado(conquista.CodigoAssociado)
			if err != nil {
				continue
			}
			conquistas[i].Associado = associado
		}
	}
	return conquistas, response.Error
}

func SetConquistasAsSent(idConquistas []uint) error {
	if len(idConquistas) == 0 {
		return nil
	}
	result := GetDB().Model(&domain.MappaConquista{}).Where("id in ?", idConquistas).UpdateColumn("notificado_chat", true)
	if result.Error == nil {
		log.Printf("Conquistas sent: %v", idConquistas)
	} else {
		log.Printf("Error on saving sent conquistas: %v", result.Error)
	}
	return result.Error
}

func GetUltimaConsultaConquista(codSecao int) time.Time {
	var secao domain.MappaSecao
	response := GetDB().Find(&secao, codSecao)
	if response.Error != nil {
		return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return secao.UltimaConsultaConquistas
}

func SetUltimaConsultaConquista(codSecao int, quando time.Time) error {
	var secao domain.MappaSecao
	response := GetDB().Find(&secao, codSecao)
	if response.Error != nil {
		return response.Error
	}
	CanIWrite()
	secao.UltimaConsultaConquistas = quando
	response = GetDB().Save(&secao)
	YouCanWrite()
	return response.Error
}

func SetConquistas(conquistas []mappa.Conquista, codigoSecao int) error {
	dbConquistas := make([]domain.MappaConquista, len(conquistas))
	for i, conquista := range conquistas {
		dbConquistas[i] = dto.CreateConquista(conquista)
		dbConquistas[i].CodigoSecao = codigoSecao
	}
	CanIWrite()
	result := GetDB().
		Clauses(clause.OnConflict{UpdateAll: true}).
		CreateInBatches(dbConquistas, 10)
	YouCanWrite()
	return result.Error
}
