package mappa

type Marcacao struct {
	CodigoAtividade       int    `json:"codigoAtividade"`
	CodigoAssociado       int    `json:"codigoAssociado"`
	DataAtividade         string `json:"dataAtividade"`
	DataStatusJovem       string `json:"dataStatusJovem"`
	DataStatusEscotista   string `json:"dataStatusEscotista"`
	StatusJovem           string `json:"statusJovem"`
	StatusEscotista       string `json:"statusEscotista"`
	DataHoraAtualizacao   string `json:"dataHoraAtualizacao"`
	CodigoUltimoEscotista int    `json:"codigoUltimoEscotista"`
	Segmento              string `json:"segmento"`
}
