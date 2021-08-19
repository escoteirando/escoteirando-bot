package services

import (
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"gorm.io/gorm/clause"
	"log"
)

func mappaGetProgressoes(ctx domain.Context) ([]mappa.Progressao, error) {
	log.Printf("Obtendo progressoes")
	body, _, err := FetchMappa("/api/progressao-atividades", ctx)
	if err != nil {
		log.Printf("Erro na obtenção das progressões: %v", err)
		return nil, err
	}
	var progressoes []mappa.Progressao
	err = json.Unmarshal(body, &progressoes)
	return progressoes, err
}

func MappaGetProgressao(codAtividade int) (domain.MappaProgressao, error) {
	var progressao domain.MappaProgressao
	if codAtividade < 1 {
		return progressao, fmt.Errorf("codAtividade must be positive [%v]", codAtividade)
	}
	progressao, err := mappaGetProgressao(codAtividade)

	if err != nil {
		log.Printf("Progressão atividade #%d não encontrada. Tentando recarregar do mappa", codAtividade)
		err = MappaUpdateProgressoes()
		progressao, err = mappaGetProgressao(codAtividade)
	}
	return progressao, err
}

func MappaValidateProgressoes(marcacoes []mappa.Marcacao) error {
	var count int64
	result := repository.GetDB().Model(&domain.MappaProgressao{}).Count(&count)
	if result.Error != nil {
		return result.Error
	}
	if count < 1 {
		return fmt.Errorf("Empty progressoes table")
	}
	set := hashset.New()
	for _, marcacao := range marcacoes {
		set.Add(marcacao.CodigoAtividade)
	}
	if set.Size() == 0 {
		return fmt.Errorf("Empty codAtividade")
	}

	result = repository.GetDB().Model(&domain.MappaProgressao{}).Where("id NOT IN ?", set.Values()).Count(&count)
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("There are %d unloaded progressoes", count)
	}
	return nil
}

func mappaGetProgressao(codAtividade int) (domain.MappaProgressao, error) {
	var progressao domain.MappaProgressao
	result := repository.GetDB().First(&progressao, codAtividade)
	if result.Error != nil {
		return progressao, result.Error
	}
	if result.RowsAffected > 0 {
		return progressao, nil
	}
	return progressao, fmt.Errorf("progressao #%v not found", codAtividade)
}

func MappaUpdateProgressoes() error {
	ctx, err := GetGenericContext()
	if err == nil {
		progressoes, err := mappaGetProgressoes(ctx)
		if err == nil && len(progressoes) > 0 {
			log.Printf("Obtidas %d progressões", len(progressoes))
			dbProgressoes := make([]domain.MappaProgressao, len(progressoes))

			for i, progressao := range progressoes {
				dbProgressoes[i] = domain.MappaProgressao{
					ID:                    progressao.Codigo,
					Descricao:             progressao.Descricao,
					CodigoUeb:             progressao.CodigoUeb,
					Ordenacao:             progressao.Ordenacao,
					CodigoCaminho:         progressao.CodigoCaminho,
					CodigoDesenvolvimento: progressao.CodigoDesenvolvimento,
					NumeroGrupo:           progressao.NumeroGrupo,
					CodigoRegiao:          progressao.CodigoRegiao,
					CodigoCompetencia:     progressao.CodigoCompetencia,
					Segmento:              progressao.Segmento,
					AreaDesenvolvimento:   progressao.AreaDesenvolvimento,
				}
			}

			result := repository.GetDB().Clauses(clause.OnConflict{
				UpdateAll: true,
			}).CreateInBatches(dbProgressoes, 10)
			if result.Error != nil {
				log.Printf("Erro na persistência das progressões: %v", result.Error)
				return result.Error
			}
			log.Printf("Persistência das progressões com sucesso: %d", len(progressoes))
			err = nil
		}
	}

	return err

}

//def get_progressoes(self, login: Login,
//ramo: Ramo = Ramo.Todos) -> List[Progressao]:
//""" Retorna todas as progressões disponíveis """
//if not self._mappa._set_auth(login):
//return
//if ramo is Ramo.Alcateia:
//caminhos = [1, 2, 3]
//elif ramo is Ramo.TropaEscoteira:
//caminhos = [4, 5, 6]
//elif ramo is Ramo.TropaSenior:
//caminhos = [11, 12]
//elif ramo is Ramo.ClaPioneiro:
//caminhos = [15, 16]
//else:
//caminhos = [1, 2, 3, 4, 5, 6, 11, 12,
//13, 14, 15, 16, 17, 18, 19, 20]
//filter = {"filter":
//{"where": {
//"numeroGrupo": None,
//"codigoRegiao": None,
//"codigoCaminho": {
//"inq": caminhos}
//}}}
//response = self._http.get(
//'/api/progressao-atividades', params=filter)
//progressoes = []
//if response.is_ok:
//try:
//for prg in response.content:
//progressoes.append(Progressao(**prg))
//except Exception as exc:
//self.LOG.error('EXCEPTION PARSING PROGRESSOES %s : %s',
//response.content, str(exc))
//
//return progressoes
