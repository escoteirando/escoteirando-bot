package repository

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
)

type CachedMappaProgressao struct {
	progressoes map[int]domain.MappaProgressao
}

func NewCachedMappaProgressao() CachedMappaProgressao {
	cached := CachedMappaProgressao{progressoes: map[int]domain.MappaProgressao{}}
	return cached
}

func (cache *CachedMappaProgressao) GetProgressao(codAtividade int) (domain.MappaProgressao, error) {
	if val, ok := cache.progressoes[codAtividade]; ok {
		return val, nil
	}
	var progressao domain.MappaProgressao
	response := GetDB().First(&progressao, codAtividade)
	err := parseResponse(response, "GetProgressao", codAtividade)
	if err == nil {
		cache.progressoes[codAtividade] = progressao
	}
	return progressao, err
}
