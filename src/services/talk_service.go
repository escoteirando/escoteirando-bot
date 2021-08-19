package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"log"
)

func ParseUserConversation(ctx domain.Context, message *tgbotapi.Message, ) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	
	SendTextMessage(ctx.ChatId, message.Text, message.MessageID)

}
