package domain

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

type TelegramMessage struct {
	ChannelId       int64
	ResponseTo      int
	Message         tgbotapi.Chattable
	ResponseChannel chan tgbotapi.Message
	RetryCount      int
	AutoDestruct    time.Time
	MessageId       int
}