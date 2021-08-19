package services

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/repository"

)


func GetChatOffset() int {
	var db = repository.GetDB()
	var config domain.Config
	result := db.First(&config)
	if result.RowsAffected == 0 || result.Error != nil {
		return 0
	}
	return config.LastMessageOffset
}

func SetChatOffset(offset int) bool {
	var db = repository.GetDB()
	var config domain.Config
	result := db.First(&config)
	if result.RowsAffected == 0 || result.Error != nil {
		config = domain.Config{
			LastMessageOffset: offset,
			AdminGroup:        0,
		}
		result = db.Create(&config)

	} else {
		result = db.Model(&config).Update("LastMessageOffset", offset)
	}
	return result.RowsAffected > 0
}
