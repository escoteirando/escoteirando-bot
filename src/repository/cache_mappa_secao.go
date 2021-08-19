package repository

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
)

type CachedMappaSecao struct {
	secoes map[int]domain.MappaSecao
}

func NewCachedMappaSecao() CachedMappaSecao {
	cached := CachedMappaSecao{secoes: map[int]domain.MappaSecao{}}
	return cached
}
func (cache *CachedMappaSecao) GetSecao(codSecao int) (domain.MappaSecao, error) {
	if val, ok := cache.secoes[codSecao]; ok {
		return val, nil
	}
	var secao domain.MappaSecao
	response := GetDB().First(&secao, codSecao)
	err := parseResponse(response, "GetSecao", codSecao)
	if err == nil {
		cache.secoes[codSecao] = secao
	}
	return secao, err
}
