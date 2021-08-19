package repository

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"gorm.io/gorm"
	"log"
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
	err := parseResponse(response, "GetAssociado", codAssociado)
	if err == nil {
		cache.associados[codAssociado] = progressao
	}

	return progressao, err
}

func parseResponse(response *gorm.DB, description string, id interface{}) error {
	var err error
	if response.Error != nil {
		err = response.Error
	} else if response.RowsAffected < 1 {
		err = fmt.Errorf("returns empty")
	}
	if err != nil {
		log.Printf("Query [%s(%v)] exception: %v", description, id, err)
	}
	return err
}
