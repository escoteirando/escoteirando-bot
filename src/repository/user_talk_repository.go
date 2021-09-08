package repository

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
)

func GetLastTalk(talkType int, userId int, chatId int64) (domain.UserTalk, error) {
	var talk domain.UserTalk
	result := GetDB().
		Order("created_at").
		Where("talk_type=? and user_id = ? and chat_id = ?", talkType, userId, chatId).
		Last(&talk)
	return talk, result.Error
}

func AddTalk(talkType int, userId int, userName string, chatId int64, messageId int) error {
	talk := domain.UserTalk{
		ChatId:    chatId,
		MessageId: messageId,
		UserId:    userId,
		UserName:  userName,
		TalkType:  talkType,
	}
	result := GetDB().Save(&talk)
	return result.Error
}
