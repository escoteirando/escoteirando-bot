package services

import (
	"encoding/json"
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"log"
	"sort"
	"time"
)

func MappaGetMarcacoes(ctx domain.Context) ([]mappa.Marcacao, error) {
	// Obter a data da última consulta de marcações
	if ctx.CodSecao < 1 {
		return nil, fmt.Errorf("seção não foi configurada neste chat")
	}
	ultimaMarcacao := repository.GetUltimaConsultaMarcacao(ctx.CodSecao)
	marcacoes, err := mappaGetMarcacoes(ctx, ultimaMarcacao)
	if err != nil {
		return nil, err
	}
	if len(marcacoes) == 0 {
		return make([]mappa.Marcacao, 0), nil
	}
	err = MappaValidateProgressoes(marcacoes)
	if err != nil {
		log.Printf("Reloading progressoes: %v", err)
		err = MappaUpdateProgressoes()
		if err != nil {
			log.Printf("Falha na atualização de progressões: %v", err)
			return marcacoes, err
		}
	}
	if err = repository.SetMarcacoes(marcacoes); err != nil {
		return nil, err
	}
	repository.SetUltimaConsultaMarcacao(ctx.CodSecao, time.Now())
	return marcacoes, nil
}

func mappaGetMarcacoes(ctx domain.Context, desde time.Time) ([]mappa.Marcacao, error) {
	if desde.Before(utils.Epoch) {
		desde = utils.Epoch
	}
	marcacoesDesde := desde.Format("2006-01-02T15:04:05.000Z")
	body, _, err := FetchMappa(fmt.Sprintf("/api/marcacoes/v2/updates?dataHoraUltimaAtualizacao=%s&codigoSecao=%d", marcacoesDesde, ctx.CodSecao), ctx)
	var marcacoes mappa.Marcacoes
	if err != nil {
		return marcacoes.Values, err
	}
	err = json.Unmarshal(body, &marcacoes)
	return marcacoes.Values, err
}

/*
def get_marcacoes(self, login: Login, cod_secao: int) -> Marcacoes:
if not self._mappa._set_auth(login):
return
url = '/api/marcacoes/v2/updates?dataHoraUltimaAtualizacao='\
'{0}&codigoSecao={1}'.format(
"1970-01-01T00:00:00.000Z",
cod_secao)

response = self._http.get(
url, description="Marcações", max_age=86400)

if response.is_ok:
try:
marcacoes = Marcacoes(**response.content)
return marcacoes
except Exception as exc:
self.LOG.error('ERROR PARSING MARCACOES %s - %s',
response.content, str(exc))


class Marcacoes(BaseModel):
	dataHora: datetime
	values: List[Marcacao]

class Marcacao(BaseModel):

codigoAtividade: int
codigoAssociado: int
dataAtividade: Union[datetime, None]
dataStatusJovem: Union[datetime, None]
dataStatusEscotista: Union[datetime, None]
statusJovem: Union[str, None]
statusEscotista: Union[str, None]
dataHoraAtualizacao: datetime
codigoUltimoEscotista: Union[int, None]
segmento: str
*/

func MappaGroupMarcacoes(marcacoes []domain.MappaMarcacao) []domain.GroupedMarcacoes {

	groupMap := map[string]*domain.GroupedMarcacoes{}
	cacheSubSecao := repository.NewCachedMappaSubSecao()
	cacheSecao := repository.NewCachedMappaSecao()
	for _, marcacao := range marcacoes {
		subSecao, err := cacheSubSecao.GetSubSecao(marcacao.Associado.CodigoEquipe)
		if err != nil {
			continue
		}
		secao, err := cacheSecao.GetSecao(subSecao.CodSecao)
		if err != nil {
			continue
		}

		key := fmt.Sprintf("%d_%s_%d", secao.ID, marcacao.DataAtividade.Format("20060102"), marcacao.CodigoAtividade)

		if val, ok := groupMap[key]; ok {
			val.Associados = append(val.Associados, marcacao.Associado)
			val.Marcacoes[marcacao.ID] = false
		} else {
			mmap := make(map[uint]bool)
			mmap[marcacao.ID] = true
			groupMap[key] = &domain.GroupedMarcacoes{
				CodSecao:   secao.ID,
				Data:       marcacao.DataAtividade,
				Progressao: marcacao.Progressao,
				Associados: []domain.MappaAssociado{
					marcacao.Associado,
				},
				Marcacoes: mmap,
			}
		}
	}
	keys := make([]string, len(groupMap))
	i := 0
	for key := range groupMap {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	result := make([]domain.GroupedMarcacoes, len(groupMap))
	i = 0
	for _, key := range keys {
		result[i] = *groupMap[key]
		i++
	}

	return result
}
