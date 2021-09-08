package repository

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/domain"
)

func GetEscotista(escotistaId int) (domain.MappaEscotista, error) {
	var escotista domain.MappaEscotista
	result := GetDB().First(&escotista, escotistaId)
	if result.Error != nil {
		return escotista, result.Error
	}
	if result.RowsAffected < 1 {
		return escotista, fmt.Errorf("escotista not found #%d", escotistaId)
	}
	return escotista, nil
}
