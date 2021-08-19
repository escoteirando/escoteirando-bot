package mappa

type Associado struct {
	Codigo                  int    `json:"codigo"`
	Nome                    string `json:"nome"`
	CodigoFoto              int    `json:"codigoFoto"`
	CodigoEquipe            int    `json:"codigoEquipe"`
	Username                int    `json:"username"`
	NumeroDigito            int    `json:"numeroDigito"`
	DataNascimento          string `json:"dataNascimento"`
	DataValidade            string `json:"dataValidade"`
	NomeAbreviado           string `json:"nomeAbreviado"`
	Sexo                    string `json:"sexo"`
	CodigoRamo              int    `json:"codigoRamo"`
	CodigoCategoria         int    `json:"codigoCategoria"`
	CodigoSegundaCategoria  int    `json:"codigoSegundaCategoria"`
	CodigoTerceiraCategoria int    `json:"codigoTerceiraCategoria"`
	LinhaFormacao           string `json:"linhaFormacao"`
	CodigoRamoAdulto        int    `json:"codigoRamoAdulto"`
	DataAcompanhamento      string `json:"dataAcompanhamento"`
}
