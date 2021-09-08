package services

import (
	"fmt"
	bot2 "github.com/guionardo/escoteirando-bot/src/bot"
	"github.com/guionardo/escoteirando-bot/src/data"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"log"
)

func PublishTodayHolidays() error {
	holidays := data.GetHolidaysFromToday()
	if len(holidays) == 0 {
		log.Printf("No holidays for today")
		return nil
	}
	chats, err := repository.GetAllChats()
	if err != nil {
		return fmt.Errorf("error publishing holidays: %v", err)
	}
	if len(chats) == 0 {
		log.Printf("No chats by now")
		return nil
	}

	holidayStr := ""
	for _, holiday := range holidays {
		holidayStr = fmt.Sprintf("%s %s (%02d/%02d) %s\n", holidayStr, holiday.Description, holiday.Day, holiday.Month, holiday.Name)
	}
	for _, chat := range chats {
		bot2.GetCurrentBot().SendTextMessage(chat.ID, holidayStr)
	}
	log.Printf("Sent %d holiday messages", len(chats))
	return err
}
