package repository

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/domain"
)

func GetFromChat(chatId int64) (*domain.Context, error) {
	chat, err := GetChat(chatId)
	if err == nil {
		return &domain.Context{
			CodSecao:       chat.CodSecao,
			ChatName:       fmt.Sprintf("%s:%s", chat.Type, chat.Title),
			ChatId:         chat.ID,
			AuthKey:        chat.AuthKey,
			MappaUserId:    chat.MappaUserId,
			AuthValidUntil: chat.AuthValidUntil,
		}, nil
	}
	return nil, err
}


