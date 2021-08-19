package mappa

type Marcacoes struct {
	DataHora string     `json:"dataHora"`
	Values   []Marcacao `json:"values"`
}
