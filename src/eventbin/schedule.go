package eventbin

import (
	"log"
	"time"
)

type (
	WorkFunc    func(schedule Schedule) error
	StartupFunc func() error
	Schedule    struct {
		Name               string
		LastRun            time.Time
		Interval           time.Duration
		WorkFunc           WorkFunc
		LastError          error
		RoundStarting      bool
		CanStartImediately bool
		PreCondition       StartupFunc
	}
)

func TruePreCondition() error {
	return nil
}

func CreateSchedule(name string, interval time.Duration, workFunc WorkFunc, preConditionFunc StartupFunc) *Schedule {
	log.Printf("Created schedule %s (%v)", name, interval)
	return &Schedule{
		Name:          name,
		LastRun:       time.Time{},
		Interval:      interval,
		WorkFunc:      workFunc,
		LastError:     nil,
		RoundStarting: false,
		PreCondition:  preConditionFunc,
	}
}

func (schedule *Schedule) RoundStart() *Schedule {
	schedule.RoundStarting = true
	nextRun := schedule.NextRun().Round(schedule.Interval)
	if nextRun.Before(time.Now()) {
		nextRun = nextRun.Add(schedule.Interval)
	}
	schedule.LastRun = nextRun.Add(-schedule.Interval)
	log.Printf("Setting next run of schedule '%s' @ %v", schedule.Name, nextRun)
	return schedule
}

func (schedule *Schedule) CanStartNow() *Schedule {
	schedule.LastRun = time.Time{}
	schedule.CanStartImediately = true
	log.Printf("Schedule %s can start imediately", schedule.Name)
	return schedule
}

func (schedule *Schedule) CanRunNow() bool {
	return schedule.LastRun.Add(schedule.Interval).Before(time.Now())
}

func (schedule *Schedule) NextRun() time.Time {
	nextRun := schedule.LastRun.Add(schedule.Interval)
	if schedule.RoundStarting {
		nextRun = nextRun.Round(schedule.Interval)
	}
	if nextRun.Before(time.Now()) {
		nextRun = time.Now()
	}
	return nextRun
}

func (schedule *Schedule) Run() error {
	err := schedule.WorkFunc(*schedule)
	schedule.LastRun = time.Now()
	schedule.LastError = err
	if err != nil {
		log.Printf("Error when running schedule %s - %v", schedule.Name, err)
	}
	return err
}
