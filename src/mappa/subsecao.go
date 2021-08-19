package mappa

type SubSecao struct {
	Codigo          int         `json:"codigo"`
	Nome            string      `json:"nome"`
	CodigoSecao     int         `json:"codigoSecao"`
	CodigoLider     int         `json:"codigoLider"`
	CodigoViceLider int         `json:"codigoViceLider"`
	Associados      []Associado `json:"associados"`
}
