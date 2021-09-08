package services

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	bot2 "github.com/guionardo/escoteirando-bot/src/bot"
	"github.com/guionardo/escoteirando-bot/src/consts"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/dto"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var httpClient = &http.Client{}

func FetchMappa(url string, ctx domain.Context) ([]byte, int, error) {
	env := utils.GetEnvironmentSetup()
	mappaUrl := fmt.Sprintf("%s/%s", env.MappaProxyUrl, url)
	req, err := http.NewRequest("GET", mappaUrl, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %s - %v", url, err)
	}
	req.Header.Set("Authorization", ctx.AuthKey)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %s - %v", url, err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %s - %v", url, err)
	}
	saveBody(url, body)
	return body, res.StatusCode, nil
}

func saveBody(url string, body []byte) {
	if !utils.GetEnvironmentSetup().BotDebug {
		return
	}
	wd, _ := os.Getwd()
	filePath := path.Join(wd, ".requests")

	if fi, err := os.Stat(filePath); !fi.IsDir() {
		err = os.Mkdir(filePath, 0666)
		if err != nil {
			log.Printf("Failed to create path %s %v", filePath, err)
			return
		}
	}
	fileUrl := strings.SplitN(strings.ReplaceAll(url, "/", "_"), "?", 1)[0]
	if len(fileUrl) > 20 {
		fileUrl = fileUrl[0:20]
	}

	fileName := path.Join(filePath, time.Now().Format("20060102_150405_")+fileUrl)
	err := os.WriteFile(fileName, body, 0666)
	if err != nil {
		log.Printf("Failed to write response body %s %v", url, err)
	} else {
		log.Printf("Response body saved %s %s", url, fileName)
	}
}

func MappaSetupAuth(ctx *domain.Context, msg tgbotapi.Message) error {

	msg, err := EditTextMessage(msg, consts.HourglassNotDone+" Identificando escotista")
	if err != nil {
		return fmt.Errorf("erro na identificação do escotista")
	}
	// Salvando o escotista
	escotista, err := mappaSaveEscotista(*ctx, msg)
	if err != nil {
		return err
	}
	if err = mappaSaveAssociado(*ctx, escotista, msg); err != nil {
		return err
	}

	return err
}

func mappaSaveSubSecoes(ctx domain.Context, secoes []domain.MappaSecao, mainMsg tgbotapi.Message) error {
	cbot:=bot2.GetCurrentBot()
	for _, secao := range secoes {
		msg := cbot.SendTextReply(ctx.ChatId, fmt.Sprintf("%s Obtendo detalhes da seção %s...", consts.HourglassNotDone, secao.ToString()), mainMsg.MessageID)
		subSecoes, err := MappaGetEquipe(ctx, ctx.MappaUserId, secao.ID)
		if err != nil {
			return err
		}
		for _, subSecao := range subSecoes {
			msgSec := cbot.SendTextReply(ctx.ChatId, fmt.Sprintf("%s Obtendo detalhes da subseção %s...", consts.HourglassNotDone, subSecao.Nome), msg.MessageID)
			dbSubSecao := dto.CreateSubsecao(subSecao)
			err = repository.SaveSubSecao(dbSubSecao)
			if err != nil {
				bot2.GetCurrentBot().EditMessage(ctx.ChatId, msgSec.MessageID, fmt.Sprintf("%s Erro ao gravar dados da seção %s - %v", consts.Warning, subSecao.Nome, err))
				return err
			}
			bot2.GetCurrentBot().EditMessage(ctx.ChatId, msgSec.MessageID, fmt.Sprintf("%s SubSeção: %s", consts.HourglassNotDone, subSecao.Nome))

			nomesAssociados := make([]string, len(subSecao.Associados))
			for index, associado := range subSecao.Associados {
				dbAssociado, err := MappaSaveAssociadoEquipe(ctx, associado, msg, false)
				if err != nil {
					bot2.GetCurrentBot().EditMessage(ctx.ChatId, msgSec.MessageID, fmt.Sprintf("%s Erro ao gravar dados do associado %s - %v", consts.Warning, associado.Nome, err))
					return err
				}
				nomesAssociados[index] = dbAssociado.ToString()
			}
			bot2.GetCurrentBot().EditMessage(ctx.ChatId, msgSec.MessageID, fmt.Sprintf("%s SubSeção %s:\n%s", consts.ThumsUp, subSecao.Nome, strings.Join(nomesAssociados, "\n")))
		}
	}
	return nil
}

func MappaSaveAssociadoEquipe(ctx domain.Context, associado mappa.Associado, msg tgbotapi.Message, sendMessage bool) (domain.MappaAssociado, error) {
	dbAssociado := domain.MappaAssociado{
		ID:                      associado.Codigo,
		Nome:                    associado.Nome,
		CodigoFoto:              associado.CodigoFoto,
		CodigoEquipe:            associado.CodigoEquipe,
		Username:                associado.Username,
		NumeroDigito:            associado.NumeroDigito,
		DataNascimento:          utils.TimeParse(associado.DataNascimento),
		DataValidade:            utils.TimeParse(associado.DataValidade),
		NomeAbreviado:           associado.NomeAbreviado,
		Sexo:                    associado.Sexo,
		CodigoRamo:              associado.CodigoRamo,
		CodigoCategoria:         associado.CodigoCategoria,
		CodigoSegundaCategoria:  associado.CodigoSegundaCategoria,
		CodigoTerceiraCategoria: associado.CodigoTerceiraCategoria,
		LinhaFormacao:           associado.LinhaFormacao,
		CodigoRamoAdulto:        associado.CodigoRamoAdulto,
		DataAcompanhamento:      utils.TimeParse(associado.DataAcompanhamento),
	}

	if err := repository.SaveAssociado(dbAssociado); err != nil {
		bot2.GetCurrentBot().SendTextReply(ctx.ChatId, fmt.Sprintf("%s Erro ao gravar o associado: %v", consts.Warning, err), msg.MessageID)
		return dbAssociado, err
	}
	if sendMessage {
		bot2.GetCurrentBot().SendTextReply(ctx.ChatId, fmt.Sprintf("%s Associado gravado com sucesso: %s", consts.ThumsUp, dbAssociado.ToString()), msg.MessageID)
	}
	return dbAssociado, nil
}

func mappaSaveSecoes(ctx domain.Context, mainMsg tgbotapi.Message) ([]domain.MappaSecao, error) {
	msg := bot2.GetCurrentBot().SendTextReply(ctx.ChatId, consts.HourglassNotDone+" Obtendo dados das seções...", mainMsg.MessageID)
	secoes, err := MappaGetSecoes(ctx, ctx.MappaUserId)
	if err != nil {
		bot2.GetCurrentBot().EditMessage(ctx.ChatId, msg.MessageID, fmt.Sprintf("%s Erro ao carregar dados das seções: %v", consts.Warning, err))
		return nil, err
	}
	savedSecoes := make([]domain.MappaSecao, 0)
	var nomeSecoes = make([]string, 0)
	for _, secao := range secoes {
		dbSecao := domain.MappaSecao{
			ID:           secao.Codigo,
			Nome:         secao.Nome,
			CodTipoSecao: secao.CodigoTipoSecao,
			CodGrupo:     secao.CodigoGrupo,
			CodRegiao:    secao.CodigoRegiao,
		}
		err = repository.SaveSecao(dbSecao)
		if err != nil {
			bot2.GetCurrentBot().EditMessage(ctx.ChatId, msg.MessageID, fmt.Sprintf("%s Erro ao gravar a seção: %v", consts.Warning, err))
			return nil, err
		}
		savedSecoes = append(savedSecoes, dbSecao)
		nomeSecoes = append(nomeSecoes, dbSecao.ToString())
		bot2.GetCurrentBot().EditMessage(ctx.ChatId, msg.MessageID, fmt.Sprintf("%s Seção gravado com sucesso: %s", consts.ThumsUp, dbSecao.ToString()))
	}
	bot2.GetCurrentBot().EditMessage(ctx.ChatId, msg.MessageID, fmt.Sprintf("%s Seção(ões): %s", consts.ThumsUp, strings.Join(nomeSecoes, ", ")))
	return savedSecoes, nil
}

func mappaSaveGrupos(ctx domain.Context, escotista domain.MappaEscotista, mainMsg tgbotapi.Message) error {
	msg := bot2.GetCurrentBot().SendTextReply(ctx.ChatId, consts.HourglassNotDone+" Obtendo dados dos grupos escoteiros...", mainMsg.MessageID)

	grupos, err := MappaGetGrupo(ctx, escotista)
	if err != nil {
		bot2.GetCurrentBot().EditMessage(ctx.ChatId, msg.MessageID, fmt.Sprintf("%s Erro ao carregar dados do grupo: %v", consts.Warning, err))
		return err
	}
	var nomeGrupos = make([]string, 0)
	for _, grupo := range grupos {
		dbGrupo := domain.MappaGrupo{
			ID:               grupo.Codigo,
			CodigoRegiao:     grupo.CodigoRegiao,
			Nome:             grupo.Nome,
			CodigoModalidade: grupo.CodigoModalidade,
		}
		err = repository.SaveGrupo(dbGrupo)
		if err != nil {
			bot2.GetCurrentBot().EditMessage(ctx.ChatId, msg.MessageID, fmt.Sprintf("%s Erro ao gravar o grupo: %v", consts.Warning, err))
			return err
		}
		bot2.GetCurrentBot().EditMessage(ctx.ChatId, msg.MessageID, fmt.Sprintf("%s Grupo gravado com sucesso: %s", consts.ThumsUp, dbGrupo.ToString()))
		nomeGrupos = append(nomeGrupos, dbGrupo.ToString())
	}
	bot2.GetCurrentBot().EditMessage(ctx.ChatId, msg.MessageID, fmt.Sprintf("%s Grupo(s): %s", consts.ThumsUp, strings.Join(nomeGrupos, ", ")))

	return nil

}

func mappaSaveAssociado(ctx domain.Context, escotista domain.MappaEscotista, mainMsg tgbotapi.Message) error {
	msg := bot2.GetCurrentBot().SendTextReply(ctx.ChatId, consts.HourglassNotDone+" Obtendo dados do escotista associado...", mainMsg.MessageID)
	associado, err := MappaGetAssociado(ctx, escotista.CodAssociado)
	if err != nil {
		bot2.GetCurrentBot().EditMessage(ctx.ChatId, msg.MessageID, fmt.Sprintf("%s Erro ao carregar dados do associado: %v", consts.Warning, err))
		return err
	}
	_, err = MappaSaveAssociadoEquipe(ctx, associado, mainMsg, false)
	if err == nil {
		bot2.GetCurrentBot().EditMessage(ctx.ChatId, msg.MessageID, fmt.Sprintf("%s Associado gravado com sucesso: %s", consts.ThumsUp, associado.Nome))
	} else {
		bot2.GetCurrentBot().EditMessage(ctx.ChatId, msg.MessageID, fmt.Sprintf("%s Erro ao gravar o associado: %v", consts.Warning, err))
	}
	return err
}

func mappaSaveEscotista(ctx domain.Context, mainMsg tgbotapi.Message) (domain.MappaEscotista, error) {
	var escotista mappa.Escotista
	var dbEscotista domain.MappaEscotista

	msg := bot2.GetCurrentBot().SendTextReply(ctx.ChatId, consts.HourglassNotDone+" Obtendo dados do escotista...", mainMsg.MessageID)
	escotista, err := MappaGetEscotista(ctx, ctx.MappaUserId)
	if err != nil {
		bot2.GetCurrentBot().EditMessage(ctx.ChatId, msg.MessageID, fmt.Sprintf("%s Erro ao obter o escotista: %v", consts.Warning, err))
		return dbEscotista, err
	}
	dbEscotista = domain.MappaEscotista{
		ID:           escotista.Codigo,
		CodAssociado: escotista.CodigoAssociado,
		UserName:     escotista.Username,
		NomeCompleto: escotista.NomeCompleto,
		Ativo:        strings.ToUpper(escotista.Ativo) == "S",
		CodGrupo:     escotista.CodigoGrupo,
		CodRegiao:    escotista.CodigoRegiao,
		CodFoto:      escotista.CodigoFoto,
	}
	err = repository.SaveEscotista(dbEscotista)
	if err != nil {
		bot2.GetCurrentBot().EditMessage(ctx.ChatId, msg.MessageID, fmt.Sprintf("%s Erro ao gravar o escotista: %v", consts.Warning, err))
		return dbEscotista, err
	}
	bot2.GetCurrentBot().EditMessage(ctx.ChatId, msg.MessageID, fmt.Sprintf("%s Escotista gravado com sucesso: %s", consts.ThumsUp, dbEscotista.NomeCompleto))
	return dbEscotista, nil
}

func MappaGetEscotista(ctx domain.Context, userId int) (mappa.Escotista, error) {
	body, _, err := FetchMappa(fmt.Sprintf("/api/escotistas/%d", userId), ctx)
	var escotista mappa.Escotista
	if err != nil {
		return escotista, err
	}
	err = json.Unmarshal(body, &escotista)
	return escotista, err
}

func MappaGetAssociado(ctx domain.Context, codAssociado int) (mappa.Associado, error) {
	body, _, err := FetchMappa(fmt.Sprintf("/api/associados/%d", codAssociado), ctx)
	var associado mappa.Associado
	if err != nil {
		return associado, err
	}
	err = json.Unmarshal(body, &associado)

	return associado, err
}

func MappaGetSecoes(ctx domain.Context, userId int) ([]mappa.Secao, error) {
	body, _, err := FetchMappa(fmt.Sprintf("/api/escotistas/%d/secoes", userId), ctx)
	var secoes []mappa.Secao
	if err != nil {
		return secoes, err
	}
	err = json.Unmarshal(body, &secoes)
	return secoes, err
}

func MappaGetEquipe(ctx domain.Context, userId int, codSecao int) ([]mappa.SubSecao, error) {
	filter := strings.ReplaceAll("{\"include\":\"associados\"}", "\"", "%22")
	body, _, err := FetchMappa(fmt.Sprintf("/api/escotistas/%d/secoes/%d/equipes?filter=%s", userId, codSecao, filter), ctx)
	var subsecoes []mappa.SubSecao
	if err != nil {
		return subsecoes, err
	}
	err = json.Unmarshal(body, &subsecoes)
	return subsecoes, err
}
