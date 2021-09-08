package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"time"
)

func CommandBirthdays(ctx domain.Context, _ *tgbotapi.Message) {
	d1 := time.Now().Add(time.Duration(-14*24) * time.Hour)
	d2 := time.Now().Add(time.Duration(+14*24) * time.Hour)
	birthdays := GetBirthdaysFromSession(ctx.CodSecao, d1, d2)
	SendBirthdaysMessage(ctx.ChatId, birthdays, d1, d2)
}
