package services

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"log"
	"strings"
	"time"
)

var fetchingInterval = time.Duration(24) * time.Hour

func FetchMarcacoes() {
	sections, err := repository.GetAllSections()
	if err != nil {
		log.Printf("Error getting sections %v", err)
		return
	}
	for _, section := range sections {
		scheduleId := fmt.Sprintf("fetch_marcacoes_%d", section.ID)
		schedule, canRun := ScheduleGet(scheduleId, fetchingInterval, true)
		if !canRun {
			continue
		}
		if err := fetchMarcacoesFromSection(section); err != nil {
			log.Print(err)
		}
		ScheduleUpdate(&schedule)
	}
}

func PublishMarcacoesToChats() {
	since := time.Now().Add(-time.Duration(30*24) * time.Hour)
	marcacoes, err := repository.GetNotSentMarcacoes(since)
	if err != nil {
		log.Printf("Error getting marcacoes %v", err)
		return
	}
	groupMarcacao := MappaGroupMarcacoes(marcacoes)
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

	_, err := SendTextMessage(chat.ID, texto, 0)
	if err == nil {
		idMarcacoesSent := make([]uint, len(grupo.Marcacoes))
		i = 0
		for id := range grupo.Marcacoes {
			idMarcacoesSent[i] = id
			i++
		}
		repository.SetMarcacoesAsSent(idMarcacoesSent)
	}
}

func fetchMarcacoesFromSection(section domain.MappaSecao) error {
	ctx, err := GetContextForSection(section.ID)
	errMessage := ""
	if err == nil {
		marcacoes, err := MappaGetMarcacoes(ctx)
		if err == nil {
			if len(marcacoes) == 0 {
				errMessage = "sem novas marcações"
			}
		}
	}
	if err != nil && len(errMessage) == 0 {
		errMessage = err.Error()
	}
	if len(errMessage) > 0 {
		return fmt.Errorf("Seção ##{section.ID}: %s", errMessage)
	}
	return nil
}
