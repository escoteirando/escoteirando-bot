package bot

import (
	"log"
)

func (bot *Bot) SendAdminTextMessage(message string) {
	if bot.AdminChatId == 0 {
		log.Printf("No admin chat is defined")
		return
	}
	msg := CreateTextMessage(bot.AdminChatId, message, 0)
	bot.PushMessage(bot.AdminChatId, msg, 0)
}
