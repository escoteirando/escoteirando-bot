package mappa

type Especialidade struct {
	Codigo           int                 `json:"codigo"`
	Descricao        string              `json:"descricao"`
	RamoConhecimento string              `json:"ramoConhecimento"`
	Prerequisito     string              `json:"prerequisito"`
	Itens            []EspecialidadeItem `json:"itens"`
}
