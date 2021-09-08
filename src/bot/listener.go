package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type (
	OnMessageFunc func(bot *tgbotapi.BotAPI, update tgbotapi.Update)
)

func (bot *Bot) StartListening(chatOffset int, onMessage OnMessageFunc, onCallBackQuery OnMessageFunc) {
	if onMessage == nil {
		log.Fatalln("Invalid onMessage argument")
	}
	if onCallBackQuery == nil {
		log.Fatalln("Invalid onCallBackQuery argument")
	}

	u := tgbotapi.NewUpdate(chatOffset)
	u.Timeout = 60
	updates, err := bot.Instance.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Failed to get updates channel: %v", err)
	}
	for update := range updates {
		log.Printf("Message: %v", update)
		if update.CallbackQuery != nil {
			log.Printf("Callback: %v", update)
			onCallBackQuery(bot.Instance, update)
			continue
		}

		if bot.UpdateChatOffsetFunc != nil {
			bot.UpdateChatOffsetFunc(update.Message.MessageID + 1)
		}
		if update.Message == nil || update.Message.From.UserName == bot.Self.UserName { // ignore any non-Message Updates
			continue
		}

		onMessage(bot.Instance, update)
	}
}


