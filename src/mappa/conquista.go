package mappa

import "time"

type Conquista struct {
	Type                           string    `json:"type"`
	DataConquista                  time.Time `json:"dataConquista"`
	CodigoAssociado                int       `json:"codigoAssociado"`
	CodigoEscotistaUltimaAlteracao int       `json:"codigoEscotistaUltimaAlteracao"`
	NumeroNivel                    int       `json:"numeroNivel"`
	CodigoEspecialidade            int       `json:"codigoEspecialidade"`
}
