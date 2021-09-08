package dto

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"gorm.io/gorm"
)

func CreateEspecialidade(especialidade mappa.Especialidade) domain.MappaEspecialidade {
	items := make([]domain.MappaEspecialidadeItem, len(especialidade.Itens))
	for i, item := range especialidade.Itens {
		items[i] = domain.MappaEspecialidadeItem{
			Id:                  especialidade.Codigo*100000 + item.Id,
			CodigoEspecialidade: item.CodigoEspecialidade,
			Descricao:           item.Descricao,
			Numero:              item.Numero,
		}
	}
	return domain.MappaEspecialidade{
		Model:            gorm.Model{ID: uint(especialidade.Codigo)},
		Descricao:        especialidade.Descricao,
		RamoConhecimento: especialidade.RamoConhecimento,
		Prerequisito:     especialidade.Prerequisito,
		Items:            items,
	}
}
