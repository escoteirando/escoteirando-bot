package repository

import "github.com/guionardo/escoteirando-bot/src/domain"

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
