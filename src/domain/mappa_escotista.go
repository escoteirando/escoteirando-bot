package domain

import (
	"gorm.io/gorm"
)

type (
	MappaEscotista struct {
		gorm.Model
		ID           int
		CodAssociado int
		UserName     string
		NomeCompleto string
		Ativo        bool
		CodGrupo     int
		CodRegiao    string
		CodFoto      int
	}
)
