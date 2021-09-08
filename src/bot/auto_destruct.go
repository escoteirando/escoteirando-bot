package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"log"
	"time"
)

func (bot *Bot) AutoDestructWorkerStart() {
	log.Printf("Starting auto destruct worker: %s", bot.Self.UserName)
	for message := range bot.AutoDestructChannel {
		go func(msg domain.TelegramMessage) {
			if msg.AutoDestruct.After(time.Now()) {
				time.Sleep(msg.AutoDestruct.Sub(time.Now()))
			}
			bot.MessageChannel <- domain.TelegramMessage{
				ChannelId: msg.ChannelId,
				Message:   tgbotapi.NewDeleteMessage(msg.ChannelId, msg.MessageId),
				MessageId: msg.MessageId,
			}
			log.Printf("Destructed message %d/%d", msg.ChannelId, msg.MessageId)
		}(message)
	}
}

func (bot *Bot) AutoDestructMessage(msg tgbotapi.Message, duration time.Duration) {
	if msg.Chat == nil || msg.Chat.ID == 0 || msg.MessageID < 1 {
		return
	}
	message := domain.TelegramMessage{
		ChannelId:    msg.Chat.ID,
		AutoDestruct: time.Now().Add(duration),
		MessageId:    msg.MessageID,
	}
	bot.AutoDestructChannel <- message
}