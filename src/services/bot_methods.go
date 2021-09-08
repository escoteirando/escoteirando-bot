package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func BotOnCallBackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Printf("OnCallBackQuery %v", update)
	ctx := GetContext(update.CallbackQuery.Message)
	ctx.UserIsAdmin = GetUserIsAdmin(update.CallbackQuery.Message.Chat, update.CallbackQuery.From)
	bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
	ParseCommandString(*ctx, update.CallbackQuery.Message, update.CallbackQuery.Data)
}

func BotOnMessage(_ *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Printf("OnMessage %v", update)
	ctx := GetContext(update.Message)
	if ctx == nil {
		return
	}
	ctx.UserIsAdmin = GetUserIsAdmin(update.Message.Chat, update.Message.From)

	NewUsersMessage(update.Message.NewChatMembers, update.Message.Chat.ID)

	if ctx.CodSecao < 1 {
		log.Printf("%s", ctx.ChatName)
	}
	if update.Message.IsCommand() {
		ParseCommand(*ctx, update.Message)
	} else {
		ParseUserConversation(*ctx, update.Message, update)
	}

}
