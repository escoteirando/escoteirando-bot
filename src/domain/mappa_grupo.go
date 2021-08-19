package domain

import (
	"fmt"
	"gorm.io/gorm"
)

type (
	MappaGrupo struct {
		gorm.Model
		ID               int
		CodigoRegiao     string
		Nome             string
		CodigoModalidade int
	}
	MappaGrupoRequest struct {
		Where MappaGrupoRequestWhere `json:"where"`
	}
	MappaGrupoRequestWhere struct {
		Codigo       int    `json:"codigo"`
		CodigoRegiao string `json:"codigoRegiao"`
	}
)

func (grupo *MappaGrupo) ToString() string {
	return fmt.Sprintf("%d/%s %s", grupo.ID, grupo.CodigoRegiao, grupo.Nome)
}
