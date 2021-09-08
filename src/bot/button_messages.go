package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *Bot) SendButtonMessage(chatId int64, caption string, message string, command string) (tgbotapi.Message, error) {
	msg := CreateButtonMessage(chatId, caption, message, command)
	response, err := bot.PushMessageAndWait(chatId, msg, 0)
	if err == nil && response.MessageID < 1 {
		err = fmt.Errorf("failed to send message %d:%s", chatId, message)
	}
	return response, nil
}

func CreateButtonMessage(chatId int64, caption string, message string, command string) tgbotapi.MessageConfig {
	var button tgbotapi.InlineKeyboardButton

	if isValidUrl(command) {
		button = tgbotapi.NewInlineKeyboardButtonURL(message, command)
	} else {
		button = tgbotapi.NewInlineKeyboardButtonData(message, command)
	}
	var cmdKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			button),
	)

	msg := tgbotapi.NewMessage(chatId, caption)
	msg.ReplyMarkup = cmdKeyboard
	return msg
}

func CreateCommandButtons(chatId int64, caption string, buttons []CommandButton) tgbotapi.MessageConfig {
	var keyboardRow []tgbotapi.InlineKeyboardButton
	keyboardRows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, command := range buttons {
		if keyboardRow == nil {
			keyboardRow = make([]tgbotapi.InlineKeyboardButton, 0)
		}
		keyboardRow = append(keyboardRow, createCommandButton(command))
		if len(keyboardRow) == 2 {
			keyboardRows = append(keyboardRows, keyboardRow)
			keyboardRow = nil
		}
	}
	if keyboardRow != nil {
		keyboardRows = append(keyboardRows, keyboardRow)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
	msg := tgbotapi.NewMessage(chatId, caption)
	msg.ReplyMarkup = keyboard
	return msg
}
