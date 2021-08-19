package services

import (
	"fmt"
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
		secao, err := MappaGetSecao(chat.CodSecao)
		msgSecao := "Seção não identificada"
		if err == nil {
			grupo, err := MappaGetGrupoFromCode(secao.CodGrupo)
			if err == nil {
				msgSecao = fmt.Sprintf("%s do GE %s", secao.ToString(), grupo.ToString())
			}
		}

		_, err = SendTextMessage(chat.ID, fmt.Sprintf("%s Escoteirando Bot ativo!\n%s\n%s", consts.Robot, msgSecao, msgRun), 0)
		if err == nil {
			log.Printf("Greeting sent to chat #%d", chat.ID)
		} else {
			log.Printf("Greeting not sent to chat #%d %v", chat.ID, err)
		}
	}
}
