package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	bot2 "github.com/guionardo/escoteirando-bot/src/bot"
	"github.com/guionardo/escoteirando-bot/src/domain"
)

func CommandStats(ctx domain.Context, _ *tgbotapi.Message) {
	bot2.GetCurrentBot().SendTextMessage(ctx.ChatId,GetStatsAsString())
}