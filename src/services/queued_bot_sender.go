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
	MessageChannel      = make(chan domain.TelegramMessage)
	AutoDestructChannel = make(chan domain.TelegramMessage)
)

func AutoDestruct(message tgbotapi.Message, in time.Duration) {
	tgMsg := domain.TelegramMessage{MessageId: message.MessageID,
		ChannelId:    message.Chat.ID,
		AutoDestruct: time.Now().Add(in)}
	AutoDestructChannel <- tgMsg
	log.Printf("Autodestruct message %d/%d @ %v", tgMsg.ChannelId, tgMsg.MessageId, tgMsg.AutoDestruct)
}

//TODO: Remover
//func SendTextMessage(chatId int64, message string, replyToMsgId int) (tgbotapi.Message, error) {
//	msg := tgbotapi.NewMessage(chatId, message)
//	if strings.Contains(message, "</") {
//		msg.ParseMode = "html"
//	}
//	if replyToMsgId != 0 {
//		msg.ReplyToMessageID = replyToMsgId
//	}
//
//	responseMessage := enqueueAndWait(chatId, msg)
//	if responseMessage.MessageID < 1 {
//		return responseMessage, fmt.Errorf("failed to send message %d:%s", chatId, message)
//	}
//	return responseMessage, nil
//}

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

func SendCommandButtons(chatId int64, caption string, buttons []CommandButton) (tgbotapi.Message, error) {
	var keyboardRow []tgbotapi.InlineKeyboardButton
	keyboardRows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, command := range buttons {
		if keyboardRow == nil {
			keyboardRow = make([]tgbotapi.InlineKeyboardButton, 0)
		}
		keyboardRow = append(keyboardRow, createCommandButton(command))
		if len(keyboardRow) == 2 {
			keyboardRows = append(keyboardRows, keyboardRow)
			keyboardRow = nil
		}
	}
	if keyboardRow != nil {
		keyboardRows = append(keyboardRows, keyboardRow)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
	msg := tgbotapi.NewMessage(chatId, caption)
	msg.ReplyMarkup = keyboard
	responseMessage := enqueueAndWait(chatId, msg)
	if responseMessage.MessageID < 1 {
		return responseMessage, fmt.Errorf("failed to send message %d", chatId)
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
		ResponseChannel: waitingChannel,
	}
	MessageChannel <- tgMsg
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
		ResponseChannel: waitingChannel,
	}
	MessageChannel <- tgMsg
	responseMessage := <-waitingChannel
	return responseMessage
}

func enqueueAndWait(chatId int64, msg tgbotapi.MessageConfig) tgbotapi.Message {
	waitingChannel := make(chan tgbotapi.Message)
	tgMsg := domain.TelegramMessage{
		ChannelId:       chatId,
		ResponseTo:      msg.ReplyToMessageID,
		Message:         msg,
		ResponseChannel: waitingChannel,
	}
	MessageChannel <- tgMsg
	responseMessage := <-waitingChannel
	return responseMessage
}
