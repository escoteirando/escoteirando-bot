package repository

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"time"
)

func GetBirthdays(codSecao int, dataInicial time.Time, dataFinal time.Time) []domain.MappaAssociado {
	var associados []domain.MappaAssociado
	subSecoes, err := GetSubSections(codSecao)
	if err == nil {
		subSecoesIds := make([]int, len(subSecoes))
		for i, subSecao := range subSecoes {
			subSecoesIds[i] = subSecao.ID
		}

		result := GetDB().Where("strftime('%m%d',data_nascimento) between ? and ?", dataInicial.Format("0102"),
			dataFinal.Format("0102")).Where("codigo_equipe in ?", subSecoesIds).Order("strftime('%m%d',data_nascimento)").Find(&associados)
		ParseResponse(result, "Aniversários", dataInicial.Format("02/01")+" "+dataFinal.Format("02/01"))
	}
	return associados
}

func GetAssociado(codAssociado int) (domain.MappaAssociado, error) {
	var associado domain.MappaAssociado
	result := GetDB().First(&associado, codAssociado)
	if result.Error != nil {
		return associado, result.Error
	}
	if result.RowsAffected < 1 {
		return associado, fmt.Errorf("associado not found #%d", codAssociado)
	}
	return associado, nil
}
