package events

import (
	"github.com/guionardo/escoteirando-bot/src/eventbin"
	"github.com/guionardo/escoteirando-bot/src/services"
	"time"
)

func init() {
	eventbin.RegisterSchedule(
		eventbin.CreateSchedule("Birthday", time.Duration(24)*time.Hour, BirthdayEvent, eventbin.TruePreCondition).
			RoundStart().
			CanStartNow())
}

func BirthdayEvent(eventbin.Schedule) error {
	return services.PublishTodayBirthdays()
}
