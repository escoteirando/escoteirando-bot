package services

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/repository"
)

var (
	CurrentConfig *domain.Config
)

func readConfig() *domain.Config {
	if CurrentConfig == nil {
		config, err := repository.GetConfig()
		if err == nil {
			CurrentConfig = &config
		} else {
			CurrentConfig = &domain.Config{
				LastMessageOffset: 0,
				AdminChat:         0,
			}
			err = repository.SetConfig(CurrentConfig)
		}
	}
	return CurrentConfig
}

func GetChatOffset() int {
	return readConfig().LastMessageOffset
}

func SetChatOffset(offset int) bool {
	readConfig()
	CurrentConfig.LastMessageOffset = offset
	return repository.SetConfig(CurrentConfig) == nil
}

func GetAdminChatId() int64 {
	return readConfig().AdminChat
}

func SetAdminChatId(chatId int64) bool {
	readConfig()
	CurrentConfig.AdminChat = chatId
	return repository.SetConfig(CurrentConfig) == nil
}
