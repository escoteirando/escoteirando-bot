package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
)

type Command struct {
	Name         string
	Method       func(ctx domain.Context, message *tgbotapi.Message)
	JustForAdmin bool
	Description  string
}

var (
	messageCommands = make(map[string]Command, 0)
)

func RegisterCommand(command Command) {
	messageCommands[command.Name] = command
}

func ParseCommandString(ctx domain.Context, message *tgbotapi.Message, command string) bool {
	if command, ok := messageCommands[command]; ok {
		if command.JustForAdmin {

		}
		command.Method(ctx, message)
		return true
	}
	return false
}
func ParseCommand(ctx domain.Context, message *tgbotapi.Message) {
	ParseCommandString(ctx, message, message.Command())
}
