package mappa

type EspecialidadeItem struct {
	Id                  int    `json:"id"`
	CodigoEspecialidade int    `json:"codigoEspecialidade"`
	Descricao           string `json:"descricao"`
	Numero              int    `json:"numero"`
}
