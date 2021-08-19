package services

import (
	"encoding/json"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"strings"
)

func MappaGetGrupoFromCode(codGrupo int) (domain.MappaGrupo,error){
	var grupo domain.MappaGrupo
	result := repository.GetDB().Find(&grupo, codGrupo)
	return grupo, result.Error
}

func MappaGetGrupo(ctx domain.Context, escotista domain.MappaEscotista) ([]mappa.Grupo, error) {
	// 	[]GET /api/grupos?filter={%22where%22:%20{%22codigo%22:%2032,%20%22codigoRegiao%22:%20%22SC%22}} HTTP/1.1
	grupoRequest := domain.MappaGrupoRequest{
		Where: domain.MappaGrupoRequestWhere{
			Codigo:       escotista.CodGrupo,
			CodigoRegiao: escotista.CodRegiao,
		},
	}

	jsonFiltro, _ := json.Marshal(&grupoRequest)
	grupoUrl := "/api/grupos?filter=" + strings.ReplaceAll(string(jsonFiltro), "\"", "%22")

	//var filter = fmt.Sprintf("/api/grupos?filter={%22where%22:%20{%22%d%22:%2032,%20%22codigoRegiao%22:%20%22%s%22}}",escotista.CodGrupo,escotista.CodRegiao)
	//var filter = strings.ReplaceAll(fmt.Sprintf("{\"where\":{\"codigo\":%d,\"codigoRegiao\":\"%s\"}}}", escotista.CodGrupo, escotista.CodRegiao),"\"","%22")
	//url := fmt.Sprintf("/api/grupos?filter=%s", filter)
	body, _, err := FetchMappa(grupoUrl, ctx)
	var grupos []mappa.Grupo
	if err != nil {
		return grupos, err
	}
	err = json.Unmarshal(body, &grupos)
	return grupos, err
}