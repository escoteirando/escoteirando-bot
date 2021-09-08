package repository

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"time"
)

func GetAllSections() ([]domain.MappaSecao, error) {
	var secoes []domain.MappaSecao
	result := GetDB().Find(&secoes)
	return secoes, result.Error
}

func GetSectionOfSubsection(codEquipe int) (domain.MappaSecao, error) {
	var (
		subSecao domain.MappaSubSecao
		secao    domain.MappaSecao
	)
	result := GetDB().First(&subSecao, codEquipe)
	if result.Error == nil && result.RowsAffected > 0 {
		result = GetDB().First(&secao, subSecao.CodSecao)
	}
	return secao, result.Error
}

func GetSubSections(codSecao int) ([]domain.MappaSubSecao, error) {
	var subSecoes []domain.MappaSubSecao
	result := GetDB().Where("cod_secao = ?", codSecao).Find(&subSecoes)
	err := ParseResponse(result, "GetSubSections", codSecao)
	return subSecoes, err
}

func GetUserIdFromSection(codSecao int) (int, error) {
	var chat domain.Chat
	result := GetDB().Model(domain.Chat{}).Where("cod_secao = ? and auth_valid_until>?", codSecao, time.Now()).First(&chat)
	err := ParseResponse(result, "GetUserIdFromSection", codSecao)
	return chat.MappaUserId, err
}

func GetSectionsCountByType() ([]struct {
	Tipo  int
	Count int64
}, error) {
	var counters []struct {
		Tipo  int
		Count int64
	}
	result := GetDB().Table("mappa_secaos").Select("cod_tipo_secao as tipo, COUNT(*) as count").Group("cod_tipo_secao").Find(&counters)
	err := ParseResponse(result, "GetSectionsCountByType", 0)
	return counters, err
}
