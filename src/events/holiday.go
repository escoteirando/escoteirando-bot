package events

import (
	"github.com/guionardo/escoteirando-bot/src/eventbin"
	"github.com/guionardo/escoteirando-bot/src/services"
	"time"
)

func init() {

	eventbin.
		RegisterSchedule(eventbin.CreateSchedule("Holiday", time.Duration(24)*time.Hour, HolidayEvent, eventbin.TruePreCondition).
			RoundStart().
			CanStartNow())

}

func HolidayEvent(eventbin.Schedule) error {
	return services.PublishTodayHolidays()
}
