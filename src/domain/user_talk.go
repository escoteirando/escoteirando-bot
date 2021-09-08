package domain

import (
	"gorm.io/gorm"
)

type UserTalk struct {
	gorm.Model
	ChatId    int64
	MessageId int
	UserId    int
	UserName  string
	TalkType  int
}
