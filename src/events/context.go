package events

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/eventbin"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/services"
	"time"
)

func init() {
	eventbin.RegisterSchedule(eventbin.CreateSchedule("Context", time.Minute, ContextEvent, eventbin.TruePreCondition))
}

func ContextEvent(eventbin.Schedule) error {
	chats, err := repository.GetAllChats()
	var adminChat int64
	if err == nil {
		config, err := repository.GetConfig()
		if err == nil {
			adminChat = config.AdminChat
		}
	}
	if err != nil {
		return fmt.Errorf("error on context worker: %v", err)
	}
	for _, chat := range chats {
		if chat.AuthValidUntil.After(time.Now()) || chat.ID == adminChat {
			continue
		}
		if len(chat.AuthKey) == 0 {
			// Sem autorização mappa
			services.SendContextAuthSetupMessage(chat.ID)
		} else {
			// Autorização mappa vencida
			services.SendContextAuthRenewMessage(chat.ID)
		}
	}
	return nil
}
