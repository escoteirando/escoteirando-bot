package events

import (
	"github.com/guionardo/escoteirando-bot/src/eventbin"
	"github.com/guionardo/escoteirando-bot/src/services"
	"time"
)

func init() {
	eventbin.RegisterSchedule(eventbin.CreateSchedule("Stats", time.Duration(24)*time.Hour, StatsEvent, eventbin.TruePreCondition).
		RoundStart().
		CanStartNow())
}

func StatsEvent(eventbin.Schedule) error {
	services.SendMessageToAdminChat(services.GetStatsAsString())
	return nil
}
