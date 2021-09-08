package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
)

func GetAllChats() ([]domain.Chat, error) {
	var chats []domain.Chat
	response := GetDB().Find(&chats)
	return chats, response.Error
}

func GetChatsFromCodSecao(codSecao int) ([]domain.Chat, error) {
	var chats []domain.Chat
	response := GetDB().Where("cod_secao = ?", codSecao).Find(&chats)
	return chats, response.Error
}

func GetChat(chatId int64) (domain.Chat, error) {
	var chat domain.Chat
	err := ParseResponse(GetDB().First(&chat, chatId), "GetChat", chatId)
	return chat, err
}

func SaveChatFromMessage(message *tgbotapi.Message) error {
	chat := domain.Chat{
		ID:            message.Chat.ID,
		Title:         message.Chat.Title,
		Type:          message.Chat.Type,
		AllAdmin:      message.Chat.AllMembersAreAdmins,
		CodSecao:      0,
	}
	err := ParseResponse(GetDB().Create(&chat), "SaveChatFromMessage", message.Chat.ID)
	return err
}
