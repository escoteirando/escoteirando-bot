package consts

var TiposSecao = map[int]string{
	1: "Alcateia",
	2: "Tropa",
	3: "Tropa sênior",
	4: "Clã",
}

func TipoSecao(tipo int) string {
	if val, ok := TiposSecao[tipo]; ok {
		return val
	}
	return "Seção Não identificada"
}
