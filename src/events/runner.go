package events

import (
	"github.com/guionardo/escoteirando-bot/src/eventbin"
	"github.com/guionardo/escoteirando-bot/src/services"
	"time"
)

func init() {
	eventbin.RegisterStartupEvent("Runner", RunnerStartupEvent)
	eventbin.RegisterSchedule(eventbin.CreateSchedule("Runner", time.Duration(1)*time.Minute, RunnerEvent, eventbin.TruePreCondition).RoundStart())
}

func RunnerStartupEvent() error {
	services.RunnerStartRun()
	go services.SendGreetings()
	return nil
}

func RunnerEvent(eventbin.Schedule) error {
	services.RunnerSaveStatus()
	return nil
}
