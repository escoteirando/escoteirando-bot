package mappa

type Secao struct {
	Codigo          int    `json:"codigo"`
	Nome            string `json:"nome"`
	CodigoTipoSecao int    `json:"codigoTipoSecao"`
	CodigoGrupo     int    `json:"codigoGrupo"`
	CodigoRegiao    string `json:"codigoRegiao"`
}
