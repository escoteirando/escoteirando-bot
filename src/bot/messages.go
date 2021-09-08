package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/url"
	"strings"
)

type CommandButton struct {
	Message string
	Command string
}

func CreateTextMessage(chatId int64, message string, replyToMsgId int) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, message)
	if strings.Contains(message, "</") {
		msg.ParseMode = "html"
	}
	if replyToMsgId != 0 {
		msg.ReplyToMessageID = replyToMsgId
	}
	return msg
}



func CreateDeletedMessage(chatId int64, messageId int) tgbotapi.DeleteMessageConfig {
	msg := tgbotapi.NewDeleteMessage(chatId, messageId)
	return msg
}

func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
func createCommandButton(command CommandButton) tgbotapi.InlineKeyboardButton {
	var button tgbotapi.InlineKeyboardButton
	if isValidUrl(command.Command) {
		button = tgbotapi.NewInlineKeyboardButtonURL(command.Message, command.Command)
	} else {
		button = tgbotapi.NewInlineKeyboardButtonData(command.Message, command.Command)
	}
	return button
}
