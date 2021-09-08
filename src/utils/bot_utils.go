package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func ChatDescription(chat *tgbotapi.Chat) string {
	words := make([]string, 0)
	if len(chat.Title) > 0 {
		words = append(words, chat.Title)
	}
	if len(chat.Description) > 0 {
		words = append(words, chat.Description)
	}
	if len(chat.FirstName) > 0 {
		words = append(words, fmt.Sprintf("%s %s", chat.FirstName, chat.LastName))
	} else {
		if len(chat.UserName) > 0 {
			words = append(words, chat.UserName)
		}
	}
	prefix := "chat"
	switch chat.Type {
	case "private":
		prefix = "Chat privado"
		break
	case "group":
		prefix = "Grupo"
		break
	case "channel":
		prefix = "Canal"
		break
	case "supergroup":
		prefix = "Super grupo"
		break
	}
	return fmt.Sprintf("%s %s", prefix, strings.Join(words, " - "))
}

func GetMessageId(message *tgbotapi.Message) int {
	if message == nil {
		return 0
	}
	return message.MessageID
}
