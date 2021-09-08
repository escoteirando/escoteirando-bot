package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	continueSending string = "ok"
	cancelSending   string = "cancel"
)

type (
	SenderCallBackFunc func(message tgbotapi.Message, err error)
)

func (bot *Bot) SenderStart() {
	log.Printf("Starting bot sender: %s", bot.Self.UserName)
	for message := range bot.MessageChannel {
		response := bot.sendMessage(&message)
		message.ResponseChannel <- response
	}
}

//func (bot *Bot) SendMessageCallBack(chatId int64, chattable tgbotapi.Chattable, responseTo int, callBack SenderCallBackFunc) {
//	message := &domain.TelegramMessage{
//		ChannelId:       chatId,
//		ResponseTo:      responseTo,
//		Message:         chattable,
//		ResponseChannel: nil,
//	}
//	response := bot.sendMessage(message)
//	callBack(response, nil)
//}

func (bot *Bot) WaitBeforeSending() {
	waitTime := bot.IntervalBetweenMessages - time.Now().Sub(bot.LastSendTime)
	if waitTime <= 0 {
		return
	}
	log.Printf("Waiting %v before sending next message", waitTime)
	time.Sleep(waitTime)
}
// sendMessage to bot
func (bot *Bot) sendMessage(message *domain.TelegramMessage) tgbotapi.Message {
	message.RetryCount++
	var msg tgbotapi.Message
	if message.RetryCount > bot.MaxRetrySend {
		log.Print("Exceed sending retry")
		return msg
	}
	bot.WaitBeforeSending()

	msg, err := bot.Instance.Send(message.Message)
	bot.LastSendTime = time.Now()

	if err == nil {
		return msg
	}
	whatToDo := processSendError(err)
	switch whatToDo {
	case continueSending:
		return bot.sendMessage(message)
	}

	return msg
}

func (bot *Bot) SendTextMessage(chatId int64, message string) tgbotapi.Message {
	return bot.SendTextReply(chatId, message, 0)
}

func (bot *Bot) SendTextReply(chatId int64, message string, replyToMsgId int) tgbotapi.Message {
	msg := CreateTextMessage(chatId, message, replyToMsgId)
	response, _ := bot.PushMessageAndWait(chatId, msg, replyToMsgId)
	return response
}



func processSendError(err error) string {
	if err == nil {
		return continueSending
	}
	if strings.Contains(err.Error(), "Too Many Requests") {
		words := strings.Split(err.Error(), " ")
		if len(words) > 0 {
			delay, err := strconv.Atoi(words[len(words)-1])
			if err == nil {
				log.Printf("Too many requests. Waiting %d seconds", delay)
				time.Sleep(time.Duration(delay) * time.Second)
				return continueSending
			}
		}
	}
	log.Printf("Unexpected sending error: %v", err)
	return cancelSending
}
