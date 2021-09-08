package services

import (
	bot2 "github.com/guionardo/escoteirando-bot/src/bot"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"time"
)

func GetBirthdaysFromSession(codSecao int, from time.Time, to time.Time) []domain.Birthday {
	var birthdays []domain.Birthday
	if utils.GetEnvironmentSetup().BotDebug {
		birthdays = append(birthdays, domain.Birthday{
			Name: "Teste Aniversário",
			Date: time.Now(),
		})
	}
	for _, bd := range repository.GetBirthdays(codSecao, from, to) {
		birthday := domain.Birthday{
			Name: bd.Nome,
			Date: bd.DataNascimento,
		}
		birthdays = append(birthdays, birthday)
	}
	return birthdays
}

func SendBirthdaysMessage(chatId int64, birthdays []domain.Birthday, from time.Time, to time.Time) {
	msg := "<b>Aniversários entre " + from.Format("02/01") + " e " + to.Format("02/01") + "</b>\n"
	if len(birthdays) == 0 {
		msg += "<i>Sem aniversários</i>"
	} else {
		for _, birthday := range birthdays {
			msg += birthday.Name + "\n"
		}
	}
	bot2.GetCurrentBot().SendTextMessage(chatId, msg)
}
