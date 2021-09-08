package consts

var emojisRamoConhecimento = map[string]string{
	"HABILIDADES_ESCOTEIRAS": RamoHabilidadesEscoteiras,
	"SERVICOS":               RamoServicos,
	"CIENCIA_TECNOLOGIA":     RamoCienciaTecnologia,
	"CULTURA":                RamoCultura,
	"DESPORTO":               RamoDesporto,
}

func EmojiRamoConhecimento(ramo string) string {
	if val, ok := emojisRamoConhecimento[ramo]; ok {
		return val
	}
	return "📚"
}
