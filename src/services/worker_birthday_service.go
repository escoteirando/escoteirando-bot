package services

import (
	"fmt"
	bot2 "github.com/guionardo/escoteirando-bot/src/bot"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"log"
	"time"
)

func PublishTodayBirthdays() error {
	sections, err := repository.GetAllSections()
	if err != nil {
		return fmt.Errorf("error getting sections %v", err)
	}
	count := 0
	for _, section := range sections {

		birthdayStr := ""
		for _, birthday := range GetBirthdaysFromSession(section.ID, time.Now(), time.Now()) {
			birthdayStr = birthdayStr + birthday.Name + "\n"
			count++
		}
		if len(birthdayStr) == 0 {
			continue
		}

		chats, err := repository.GetChatsFromCodSecao(section.ID)
		if err != nil {
			continue
		}
		for _, chat := range chats {
			if len(birthdayStr) > 0 {
				bot2.GetCurrentBot().SendTextMessage(chat.ID, fmt.Sprintf("<b>Aniversários</b>\n%s", birthdayStr))
			}
		}
	}
	if count == 0 {
		log.Printf("No birthdays for today")
	}
	return nil
}
