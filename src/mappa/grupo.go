package mappa

type Grupo struct {
	Codigo           int    `json:"codigo"`
	CodigoRegiao     string `json:"codigoRegiao"`
	Nome             string `json:"nome"`
	CodigoModalidade int    `json:"codigoModalidade"`
}
