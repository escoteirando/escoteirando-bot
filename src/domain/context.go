package domain

import (
	"time"
)

type (
	Context struct {
		CodSecao       int
		ChatName       string
		ChatId         int64
		AuthKey        string
		MappaUserId    int
		AuthValidUntil time.Time
		UserIsAdmin    bool
	}
	FrontendContext struct {
		CId int64 `json:"cId"`
		MId int   `json:"mId"`
	}
)
