package domain

import (
	"gorm.io/gorm"
	"time"
)

type (
	// Chat "chat":{"id":-477575259,"title":"ChefeBot Lab","type":"group","all_members_are_administrators":true}
	Chat struct {
		gorm.Model
		ID             int64
		Title          string
		Type           string
		AllAdmin       bool
		CodSecao       int
		LastSetupCall  time.Time
		AuthKey        string
		AuthValidUntil time.Time
		MappaUserId    int
		LastAuth       time.Time
	}
)
