package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/url"
	"time"
)

type CommandButton struct {
	Message string
	Command string
}

var (
	bot             *tgbotapi.BotAPI
	BotUser         tgbotapi.User
	initialized     = false
	lastSendMessage = time.Now()
)
var minInterval, _ = time.ParseDuration("1s")

func createCommandButton(command CommandButton) tgbotapi.InlineKeyboardButton {
	var button tgbotapi.InlineKeyboardButton
	if isValidUrl(command.Command) {
		button = tgbotapi.NewInlineKeyboardButtonURL(command.Message, command.Command)
	} else {
		button = tgbotapi.NewInlineKeyboardButtonData(command.Message, command.Command)
	}
	return button
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
