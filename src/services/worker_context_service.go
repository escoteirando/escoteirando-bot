package services

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"time"
)

func GetExpiredAuthSections() ([]domain.Chat, error) {
	var chats []domain.Chat
	result := repository.GetDB().Where("auth_valid_until < ?", time.Now()).Find(&chats)
	return chats, result.Error
}

func GetNoSetupChats() []domain.Chat {
	var chats []domain.Chat
	result := repository.GetDB().Where("auth_key is null").Find(&chats)
	err := repository.ParseResponse(result, "GetNoSetupChats", 0)
	if err != nil {
		return make([]domain.Chat, 0)
	}
	return chats
}
