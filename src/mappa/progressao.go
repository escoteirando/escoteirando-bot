package mappa

type Progressao struct {
	Codigo                int    `json:"codigo"`
	Descricao             string `json:"descricao"`
	CodigoUeb             string `json:"codigoUeb"`
	Ordenacao             int    `json:"ordenacao"`
	CodigoCaminho         int    `json:"codigoCaminho"`
	CodigoDesenvolvimento int    `json:"codigoDesenvolvimento"`
	NumeroGrupo           int    `json:"numeroGrupo"`
	CodigoRegiao          string `json:"codigoRegiao"`
	CodigoCompetencia     int    `json:"codigoCompetencia"`
	Segmento              string `json:"segmento"`
	AreaDesenvolvimento   string `json:"areaDesenvolvimento"`
}

func (progressao *Progressao) Area()string{
	/*if isinstance(self.areaDesenvolvimento, AreaDesenvolvimento):
	return self.areaDesenvolvimento
	if self.codigoUeb and self.codigoUeb[0] in 'FICASE':
	self.areaDesenvolvimento = AreaDesenvolvimento(
		self.codigoUeb.upper()[0])
	else:
	self.areaDesenvolvimento = AreaDesenvolvimento.ERRO

	return self.areaDesenvolvimento

	 */
	return progressao.AreaDesenvolvimento
	//TODO: Implementar identificação da área de acordo com os dados.
}