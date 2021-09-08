package repository

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
)

type CachedMappaSubSecao struct {
	secoes map[int]domain.MappaSubSecao
}

func NewCachedMappaSubSecao() CachedMappaSubSecao {
	cached := CachedMappaSubSecao{secoes: map[int]domain.MappaSubSecao{}}
	return cached
}
func (cache *CachedMappaSubSecao) GetSubSecao(codSecao int) (domain.MappaSubSecao, error) {
	if val, ok := cache.secoes[codSecao]; ok {
		return val, nil
	}
	var subSecao domain.MappaSubSecao
	response := GetDB().First(&subSecao, codSecao)
	err := ParseResponse(response, "GetSubSecao", codSecao)
	if err == nil {
		cache.secoes[codSecao] = subSecao
	}
	return subSecao, err
}
