package repository

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
)

type CachedMappaAssociado struct {
	associados map[int]domain.MappaAssociado
}

func NewCachedMappaAssociado() CachedMappaAssociado {
	cached := CachedMappaAssociado{associados: map[int]domain.MappaAssociado{}}
	return cached
}

func (cache *CachedMappaAssociado) GetAssociado(codAssociado int) (domain.MappaAssociado, error) {
	if val, ok := cache.associados[codAssociado]; ok {
		return val, nil
	}
	var progressao domain.MappaAssociado
	response := GetDB().First(&progressao, codAssociado)
	err := ParseResponse(response, "GetAssociado", codAssociado)
	if err == nil {
		cache.associados[codAssociado] = progressao
	}

	return progressao, err
}


