package events

import (
	"fmt"
	bot2 "github.com/guionardo/escoteirando-bot/src/bot"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/eventbin"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/services"
	"log"
	"strings"
	"time"
)

func init() {
	//eventbin.RegisterSchedule(eventbin.
	//	CreateSchedule("MappaPublishMarcacoes", time.Duration(24)*time.Hour, MappaPublishMarcacoes, ContextNeeded).
	//	RoundStart().
	//	CanStartNow())
}

func MappaPublishMarcacoes(eventbin.Schedule) error {
	since := time.Now().Add(-time.Duration(30*24) * time.Hour)
	marcacoes, err := repository.GetNotSentMarcacoes(since)
	if err != nil {
		log.Printf("Error getting marcacoes after %v %v", since, err)
		return err
	}
	groupMarcacao := services.MappaGroupMarcacoes(marcacoes)
	secoes := repository.NewCachedMappaSecao()
	for _, grupo := range groupMarcacao {
		secao, err := secoes.GetSecao(grupo.CodSecao)
		if err != nil {
			continue
		}
		chats, err := repository.GetChatsFromCodSecao(secao.ID)
		if err != nil {
			continue
		}
		for _, chat := range chats {
			PublishMarcacoesToChat(&grupo, chat)
		}
	}
	return nil
}

func PublishMarcacoesToChat(grupo *domain.GroupedMarcacoes, chat domain.Chat) {
	associadosMap := make(map[string]bool)

	for _, associado := range grupo.Associados {
		associadosMap[associado.ToString()] = true
	}
	associados := make([]string, len(associadosMap))
	i := 0
	for associado := range associadosMap {
		associados[i] = associado
		i++
	}
	texto := fmt.Sprintf("<u>%s</u>\n<b>%s</b>\n\n<i>%s</i>", grupo.Data.Format("02/01/2006"), grupo.Progressao.ToString(), strings.Join(associados, "\n"))

	msg := bot2.GetCurrentBot().SendTextMessage(chat.ID, texto)
	if msg.MessageID != 0 {
		idMarcacoesSent := make([]uint, len(grupo.Marcacoes))
		i = 0
		for id := range grupo.Marcacoes {
			idMarcacoesSent[i] = id
			i++
		}
		repository.SetMarcacoesAsSent(idMarcacoesSent)
	}
}
