package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"log"
	"net/url"
	"strings"
	"time"
)

var (
	bot             *tgbotapi.BotAPI
	initialized     = false
	lastSendMessage = time.Now()
)
var minInterval, _ = time.ParseDuration("1s")


func GetBot() (*tgbotapi.BotAPI, error) {
	var err error
	if !initialized {
		bot, err = tgbotapi.NewBotAPI(utils.GetEnvironmentSetup().BotToken)
		if err == nil {
			bot.Debug = utils.GetEnvironmentSetup().BotDebug
			log.Printf("Authorized on account %s", bot.Self.UserName)
			initialized = true
		} else {
			log.Fatalf("Failed to get bot: %v", err)
		}
	}
	return bot, err
}

func waitToSend() {
	intervalLastSend := time.Now().Sub(lastSendMessage)
	if intervalLastSend > minInterval {
		return
	}
	wait := minInterval - time.Now().Sub(lastSendMessage)
	log.Printf("Waiting %v for next message", wait)
	time.Sleep(wait)
}

func SendMessage(chatId int64, message string, replyToMsgId int) (tgbotapi.Message, error) {
	bot, msgSent, err := getBotAndMessage()
	if err != nil {
		return msgSent, err
	}
	msg := tgbotapi.NewMessage(chatId, message)
	if strings.Contains(message, "</") {
		msg.ParseMode = "html"
	}
	if replyToMsgId != 0 {
		msg.ReplyToMessageID = replyToMsgId
	}

	msgSent, err = bot.Send(msg)
	if err != nil {
		log.Printf("Erro ao enviar mensagem: %v", err)
	}
	lastSendMessage = time.Now()
	return msgSent, err
}

func EditMessage(chatId int64, messageId int, newText string) (tgbotapi.Message, error) {
	bot, msgSent, err := getBotAndMessage()
	if err != nil {
		return msgSent, err
	}
	msg := tgbotapi.NewEditMessageText(chatId, messageId, newText)
	msgSent, err = bot.Send(msg)
	if err != nil {
		log.Printf("Erro ao enviar mensagem: %v", err)
	}
	lastSendMessage = time.Now()
	return msgSent, err
}

func DeleteMessage(chatId int64, messageId int) (tgbotapi.Message, error) {
	bot, msgSent, err := getBotAndMessage()
	if err != nil {
		return msgSent, err
	}
	msg := tgbotapi.NewDeleteMessage(chatId, messageId)
	msgSent, err = bot.Send(msg)
	if err != nil {
		log.Printf("Erro ao enviar mensagem: %v", err)
	}
	lastSendMessage = time.Now()
	return msgSent, err
}

func getBotAndMessage() (*tgbotapi.BotAPI, tgbotapi.Message, error) {
	waitToSend()
	var msgSent tgbotapi.Message
	bot, err := GetBot()
	if err != nil {
		log.Printf("Erro ao criar bot: %s", err.Error())
	}
	return bot, msgSent, err
}
func SendCommandButton(chatId int64, caption string, message string, command string) (tgbotapi.Message, error) {
	bot, msgSent, err := getBotAndMessage()
	if err != nil {
		return msgSent, err
	}

	var button tgbotapi.InlineKeyboardButton

	if isValidUrl(command) {
		button = tgbotapi.NewInlineKeyboardButtonURL(message, command)
	} else {
		button = tgbotapi.NewInlineKeyboardButtonData(message, command)
	}
	var cmdKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			button),
	)

	msg := tgbotapi.NewMessage(chatId, caption)
	msg.ReplyMarkup = cmdKeyboard
	msgSent, err = bot.Send(msg)
	if err != nil {
		log.Printf("Erro ao enviar mensagem: #{err}")
	}
	lastSendMessage = time.Now()
	return msgSent, err
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
