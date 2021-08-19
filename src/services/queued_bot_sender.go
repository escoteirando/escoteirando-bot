package services

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"log"
	"strings"
	"time"
)

var (
	messageChannel      chan domain.TelegramMessage
	autoDescructChannel chan domain.TelegramMessage
)

func AutoDestruct(message tgbotapi.Message, in time.Duration) {
	tgMsg := domain.TelegramMessage{MessageId: message.MessageID,
		ChannelId:    message.Chat.ID,
		AutoDestruct: time.Now().Add(in)}
	autoDescructChannel <- tgMsg
	log.Printf("Autodestruct message %d/%d @ %v", tgMsg.ChannelId, tgMsg.MessageId, tgMsg.AutoDestruct)
}
func SetupMessageChanel(mChannel chan domain.TelegramMessage) {
	messageChannel = mChannel
}
func SetupAutoDestructChannel(mChannel chan domain.TelegramMessage) {
	autoDescructChannel = mChannel
}

func SendTextMessage(chatId int64, message string, replyToMsgId int) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatId, message)
	if strings.Contains(message, "</") {
		msg.ParseMode = "html"
	}
	if replyToMsgId != 0 {
		msg.ReplyToMessageID = replyToMsgId
	}

	responseMessage := enqueueAndWait(chatId, msg)
	if responseMessage.MessageID < 1 {
		return responseMessage, fmt.Errorf("failed to send message %d:%s", chatId, message)
	}
	return responseMessage, nil
}

func SendButtonMessage(chatId int64, caption string, message string, command string) (tgbotapi.Message, error) {
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
	responseMessage := enqueueAndWait(chatId, msg)
	if responseMessage.MessageID < 1 {
		return responseMessage, fmt.Errorf("failed to send message %d:%s", chatId, message)
	}
	return responseMessage, nil
}


func DeleteTextMessage(chatId int64, messageId int) {
	msg := tgbotapi.NewDeleteMessage(chatId, messageId)
	responseMessage := enqueueDeleteAndWait(chatId, msg)
	log.Printf("Deleted message %v", responseMessage)
}

func enqueueDeleteAndWait(chatId int64, msg tgbotapi.DeleteMessageConfig) tgbotapi.Message {
	waitingChannel := make(chan tgbotapi.Message)
	tgMsg := domain.TelegramMessage{
		ChannelId:       chatId,
		ResponseTo:      0,
		Message:         msg,
		Type:            0,
		ResponseChannel: waitingChannel,
	}
	messageChannel <- tgMsg
	responseMessage := <-waitingChannel
	return responseMessage
}
func EditTextMessage(message tgbotapi.Message, newText string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewEditMessageText(message.Chat.ID, message.MessageID, newText)
	if strings.Contains(newText, "</") {
		msg.ParseMode = "html"
	}
	responseMessage := enqueueEditAndWait(message.Chat.ID, msg)
	if responseMessage.MessageID < 1 {
		return responseMessage, fmt.Errorf("failed to edit message %d:%s", message.Chat.ID, newText)
	}
	return responseMessage, nil
}

func enqueueEditAndWait(chatId int64, msg tgbotapi.EditMessageTextConfig) tgbotapi.Message {
	waitingChannel := make(chan tgbotapi.Message)
	tgMsg := domain.TelegramMessage{
		ChannelId:       chatId,
		ResponseTo:      0,
		Message:         msg,
		Type:            0,
		ResponseChannel: waitingChannel,
	}
	messageChannel <- tgMsg
	responseMessage := <-waitingChannel
	return responseMessage
}

func enqueueAndWait(chatId int64, msg tgbotapi.MessageConfig) tgbotapi.Message {
	waitingChannel := make(chan tgbotapi.Message)
	tgMsg := domain.TelegramMessage{
		ChannelId:       chatId,
		ResponseTo:      msg.ReplyToMessageID,
		Message:         msg,
		Type:            0,
		ResponseChannel: waitingChannel,
	}
	messageChannel <- tgMsg
	responseMessage := <-waitingChannel
	return responseMessage
}
