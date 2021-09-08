package mappa_api

import (
	"encoding/json"
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"log"
	"strings"
	"time"
)

func MappaGetEscotista(ctx domain.Context, userId int) (mappa.Escotista, error) {
	body, _, err := fetchMappa(fmt.Sprintf("/api/escotistas/%d", userId), ctx)
	var escotista mappa.Escotista
	if err != nil {
		return escotista, err
	}
	err = json.Unmarshal(body, &escotista)
	return escotista, err
}

func MappaGetAssociado(ctx domain.Context, codAssociado int) (mappa.Associado, error) {
	body, _, err := fetchMappa(fmt.Sprintf("/api/associados/%d", codAssociado), ctx)
	var associado mappa.Associado
	if err != nil {
		return associado, err
	}
	err = json.Unmarshal(body, &associado)
	return associado, err
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

	body, _, err := fetchMappa(grupoUrl, ctx)
	var grupos []mappa.Grupo
	if err != nil {
		return grupos, err
	}
	err = json.Unmarshal(body, &grupos)
	return grupos, err
}

func MappaGetSecoes(ctx domain.Context, userId int) ([]mappa.Secao, error) {
	body, _, err := fetchMappa(fmt.Sprintf("/api/escotistas/%d/secoes", userId), ctx)
	var secoes []mappa.Secao
	if err != nil {
		return secoes, err
	}
	err = json.Unmarshal(body, &secoes)
	return secoes, err
}

func MappaGetEquipe(ctx domain.Context, userId int, codSecao int) ([]mappa.SubSecao, error) {
	filter := strings.ReplaceAll("{\"include\":\"associados\"}", "\"", "%22")
	body, _, err := fetchMappa(fmt.Sprintf("/api/escotistas/%d/secoes/%d/equipes?filter=%s", userId, codSecao, filter), ctx)
	var subsecoes []mappa.SubSecao
	if err != nil {
		return subsecoes, err
	}
	err = json.Unmarshal(body, &subsecoes)
	return subsecoes, err
}

func MappaGetProgressoes(ctx domain.Context) ([]mappa.Progressao, error) {
	log.Printf("Obtendo progressoes")
	body, _, err := fetchMappa("/api/progressao-atividades", ctx)
	if err != nil {
		log.Printf("Erro na obtenção das progressões: %v", err)
		return nil, err
	}
	var progressoes []mappa.Progressao
	err = json.Unmarshal(body, &progressoes)
	return progressoes, err
}

func MappaGetMarcacoes(ctx domain.Context, desde time.Time) ([]mappa.Marcacao, error) {
	if desde.Before(utils.Epoch) {
		desde = utils.Epoch
	}
	marcacoesDesde := desde.Format("2006-01-02T15:04:05.000Z")
	body, _, err := fetchMappa(fmt.Sprintf("/api/marcacoes/v2/updates?dataHoraUltimaAtualizacao=%s&codigoSecao=%d", marcacoesDesde, ctx.CodSecao), ctx)
	var marcacoes mappa.Marcacoes
	if err == nil {
		err = json.Unmarshal(body, &marcacoes)
	}

	return marcacoes.Values, err
}

func MappaGetConquistas(ctx domain.Context, desde time.Time) ([]mappa.Conquista, error) {
	if desde.Before(utils.Epoch) {
		desde = utils.Epoch
	}
	conquistasDesde := desde.Format("2006-01-02T15:04:05.000Z")
	body, _, err := fetchMappa(fmt.Sprintf("/api/associado-conquistas/v2/updates?dataHoraUltimaAtualizacao=%s&codigoSecao=%d",
		conquistasDesde,
		ctx.CodSecao), ctx)
	var conquistas mappa.Conquistas
	if err == nil {
		err = json.Unmarshal(body, &conquistas)
	}
	return conquistas.Values, err
}

func MappaGetEspecialidades(ctx domain.Context) ([]mappa.Especialidade, error) {
	body, _, err := fetchMappa("/api/especialidades?filter[include]=itens", ctx)
	var especialidades []mappa.Especialidade
	if err == nil {
		err = json.Unmarshal(body, &especialidades)
	}
	return especialidades, err
}
