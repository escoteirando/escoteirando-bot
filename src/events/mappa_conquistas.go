package events

import (
	"fmt"
	bot2 "github.com/guionardo/escoteirando-bot/src/bot"
	"github.com/guionardo/escoteirando-bot/src/consts"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/eventbin"
	"github.com/guionardo/escoteirando-bot/src/mappa_api"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/services"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"log"
	"strings"
	"time"
)

const fetchConquistasInterval = consts.Day

func init() {
	eventbin.RegisterSchedule(eventbin.
		CreateSchedule("MappaPublishConquistas", fetchConquistasInterval, MappaPublishConquistas, ContextNeeded).
		RoundStart().
		CanStartNow())

	eventbin.RegisterSchedule(eventbin.
		CreateSchedule("MappaFetchConquistas", fetchConquistasInterval, MappaFetchConquistas, ContextNeeded).
		RoundStart().
		CanStartNow())
}

func MappaFetchConquistas(eventbin.Schedule) error {
	schedule, canRun := services.ScheduleGet("mappa_fetch_conquistas", fetchConquistasInterval, true)
	if !canRun {
		return nil
	}
	secoes, err := repository.GetAllSections()
	if err != nil {
		return err
	}
	for _, secao := range secoes {
		ctx, err := services.GetContextForSection(secao.ID)
		if err != nil {
			continue
		}
		desde := repository.GetUltimaConsultaConquista(secao.ID)
		conquistas, err := mappa_api.MappaGetConquistas(ctx, desde)
		if err != nil {
			continue
		}
		err = repository.SetConquistas(conquistas, secao.ID)
		if err == nil {
			repository.SetUltimaConsultaConquista(secao.ID, time.Now())
		}

	}
	services.ScheduleUpdate(&schedule)
	return nil
}

func MappaPublishConquistas(eventbin.Schedule) error {
	schedule, canRun := services.ScheduleGet("mappa_publish_conquistas", fetchConquistasInterval, true)
	if !canRun {
		return nil
	}
	since := time.Now().Add(-time.Duration(180*24) * time.Hour)
	conquistas, err := repository.GetNotSentConquistas(since)
	if err != nil {
		log.Printf("Error getting marcacoes after %v %v", since, err)
		return err
	}

	publishedIds := make([]uint, 0)
	for _, conquista := range conquistas {
		chats, err := repository.GetChatsFromCodSecao(conquista.CodigoSecao)
		if err != nil {
			continue
		}
		for _, chat := range chats {
			pId := PublishConquistaToChat(conquista, chat)
			if pId > 0 {
				publishedIds = append(publishedIds, pId)
			}
		}
	}
	if len(publishedIds) > 0 {
		repository.SetConquistasAsSent(publishedIds)
	}

	services.ScheduleUpdate(&schedule)
	return nil
}

func PublishConquistaToChat(conquista domain.MappaConquista, chat domain.Chat) uint {
	if conquista.CodigoEspecialidade > 0 {
		return PublishConquistaEspecialidadeToChat(conquista, chat)
	}
	texto := fmt.Sprintf("<u>%s</u>\n<b>Conquista %s</b>\n<i>%s</i>",
		conquista.DataConquista.Format("02/01/2006"),
		utils.ToCamelCase(conquista.Type),
		conquista.Associado.ToString())
	msg := bot2.GetCurrentBot().SendTextMessage(chat.ID, texto)
	if msg.MessageID != 0 {
		return conquista.ID
	}
	return 0
}

func PublishConquistaEspecialidadeToChat(conquista domain.MappaConquista, chat domain.Chat) uint {
	especialidade, err := repository.GetEspecialidade(conquista.CodigoEspecialidade)
	if err != nil {
		return 0
	}
	texto := fmt.Sprintf("<u>%s</u>\n<b>Especialidade %s %s</b>\n%s\n<i>%s</i>",
		conquista.DataConquista.Format("02/01/2006"),
		especialidade.Descricao,
		strings.Repeat(consts.Star, conquista.NumeroNivel),
		consts.EmojiRamoConhecimento(especialidade.RamoConhecimento)+" "+utils.ToCamelCase(especialidade.RamoConhecimento),
		conquista.Associado.ToString())
	msg := bot2.GetCurrentBot().SendTextMessage(chat.ID, texto)
	if msg.MessageID != 0 {
		return conquista.ID
	}
	return 0
}
