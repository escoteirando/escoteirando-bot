package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"log"
	"time"
)

func ReadChannelTimeout(channel chan tgbotapi.Message, timeout time.Duration) (tgbotapi.Message, error) {
	select {
	case res := <-channel:
		return res, nil
	case <-time.After(timeout):
		return tgbotapi.Message{}, fmt.Errorf("timed out waiting for channel (%v)", timeout)
	}
}

// PushMessage Push message to message queue (no control of return)
func (bot *Bot) PushMessage(chatId int64, chattable tgbotapi.Chattable, responseTo int) {
	message := domain.TelegramMessage{
		ChannelId:  chatId,
		ResponseTo: responseTo,
		Message:    chattable,
	}
	bot.MessageChannel <- message
}

// PushMessageAndWait Push message to message queue and waits until response reaches response channel
func (bot *Bot) PushMessageAndWait(chatId int64, chattable tgbotapi.Chattable, responseTo int) (tgbotapi.Message, error) {
	responseChannel := make(chan tgbotapi.Message)
	message := domain.TelegramMessage{
		ChannelId:       chatId,
		ResponseTo:      responseTo,
		Message:         chattable,
		ResponseChannel: responseChannel,
	}
	bot.MessageChannel <- message
	response, err := ReadChannelTimeout(responseChannel, time.Duration(30)*time.Second)
	if err != nil {
		log.Printf("Error on sending message  %v %v", chattable, err)
	}
	return response, err
}
