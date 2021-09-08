package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *Bot) EditMessage(chatId int64, messageId int, newText string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewEditMessageText(chatId, messageId, newText)
	response, err := bot.PushMessageAndWait(chatId, msg, 0)
	return response, err
}
