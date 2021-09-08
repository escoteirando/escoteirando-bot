package services

import (
	"fmt"
	bot2 "github.com/guionardo/escoteirando-bot/src/bot"
	"github.com/guionardo/escoteirando-bot/src/consts"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"log"
)

func SendGreetings() {
	chats, err := repository.GetAllChats()
	if err != nil {
		log.Printf("Não foi possível obter os chats registrados")
		return
	}
	msgRun := "Primeira execução"
	lastRun, err := RunnerLastCompleteRun()
	if err == nil {
		msgRun = fmt.Sprintf("Última execução: %s", lastRun.ToString())
	}

	for _, chat := range chats {
		msgSecao := "Seção não identificada"
		if chat.ID == GetAdminChatId() {
			msgSecao = "Chat de administração geral"
		} else {
			secao, err := MappaGetSecao(chat.CodSecao)

			if err == nil {
				grupo, err := MappaGetGrupoFromCode(secao.CodGrupo)
				if err == nil {
					msgSecao = fmt.Sprintf("%s do GE %s", secao.ToString(), grupo.ToString())
				}
			}
		}

		bot2.GetCurrentBot().SendTextMessage(chat.ID, fmt.Sprintf("%s Escoteirando Bot ativo!\n%s\n%s\n%s",
			consts.Robot,
			GreetingTime(),
			msgSecao, msgRun))

		log.Printf("Greeting sent to chat #%d", chat.ID)
	}
}
