package repository

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/dto"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"gorm.io/gorm/clause"
)

func SaveEspecialidades(especialidades []mappa.Especialidade) error {
	dbEsp := make([]domain.MappaEspecialidade, len(especialidades))
	for i, esp := range especialidades {
		dbEsp[i] = dto.CreateEspecialidade(esp)
	}
	CanIWrite()
	result := GetDB().
		Clauses(clause.OnConflict{UpdateAll: true}).
		CreateInBatches(dbEsp, 10)
	YouCanWrite()
	return result.Error
}

func GetEspecialidade(codEspecialidade int) (domain.MappaEspecialidade, error) {
	var especialidade domain.MappaEspecialidade
	result := GetDB().First(&especialidade, codEspecialidade)
	return especialidade, result.Error
}
