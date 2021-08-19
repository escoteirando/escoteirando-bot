package mappa

type Escotista struct {
	Codigo          int    `json:"codigo"`
	CodigoAssociado int    `json:"codigoAssociado"`
	Username        string `json:"username"`
	NomeCompleto    string `json:"nomeCompleto"`
	Ativo           string `json:"ativo"`
	CodigoGrupo     int    `json:"codigoGrupo"`
	CodigoRegiao    string `json:"codigoRegiao"`
	CodigoFoto      int    `json:"codigoFoto"`
}
