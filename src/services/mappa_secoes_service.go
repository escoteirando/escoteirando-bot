package services

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/repository"
)

func MappaGetSecao(codSecao int) (domain.MappaSecao, error) {
	var secao domain.MappaSecao
	result := repository.GetDB().Find(&secao, codSecao)
	return secao, result.Error
}

//func MappaSyncSecao(codSecao int) (domain.MappaSecao,error){
//return nil,nil
//}
// TODO: Implementar sync de seção
