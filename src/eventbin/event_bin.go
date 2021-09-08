package eventbin

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"log"
	"os"
	"time"
)

type (
	EventBin struct {
		Schedules          []*Schedule
		RunInterval        time.Duration
		CanStartImediately bool
	}
)

var Schedules = make(map[string]*Schedule, 0)
var StartupEvents = make(map[string]StartupFunc, 0)

func RegisterSchedule(schedule *Schedule) {
	if !utils.IsDisabledFeature(schedule.Name) {
		Schedules[schedule.Name] = schedule
	}
}

func RegisterStartupEvent(name string, f StartupFunc) {
	StartupEvents[name] = f
}

func CreateEventBin() *EventBin {
	return &EventBin{
		Schedules:   make([]*Schedule, 0),
		RunInterval: time.Duration(1) * time.Hour,
	}
}

func (bin *EventBin) AddSchedules() *EventBin {
	scheduleCount := 0
	for _, schedule := range Schedules {
		bin.AddSchedule(schedule)
		scheduleCount++
	}
	if scheduleCount == 0 {
		log.Printf("No schedules was added to event bin")
	}
	return bin
}

func getFirstPendingSchedule(bin *EventBin, withError bool) (*Schedule, error) {
	var currentPending *Schedule
	for _, schedule := range bin.Schedules {
		if !schedule.CanRunNow() {
			continue
		}
		if (withError && schedule.LastError == nil) || (!withError && schedule.LastError != nil) {
			continue
		}
		if currentPending == nil {
			currentPending = schedule
			continue
		}
		if schedule.NextRun().Before(currentPending.NextRun()) {
			currentPending = schedule
		}
	}
	if currentPending == nil {
		return nil, fmt.Errorf("no pending schedule")
	}
	return currentPending, nil
}
func (bin *EventBin) GetFirstPendingSchedule() (*Schedule, error) {
	pendingSchedule, err := getFirstPendingSchedule(bin, false)
	if err == nil {
		return pendingSchedule, nil
	}
	return getFirstPendingSchedule(bin, true)
}

func (bin *EventBin) AddSchedule(schedule *Schedule) *EventBin {
	bin.Schedules = append(bin.Schedules, schedule)
	if schedule.Interval < bin.RunInterval {
		bin.RunInterval = schedule.Interval
	}
	if schedule.CanStartImediately {
		bin.CanStartImediately = true
	}
	return bin
}

func (bin *EventBin) WaitUntilNextSchedule() (*Schedule, error) {
	nextSchedule, err := bin.GetFirstPendingSchedule()
	if err != nil {
		return nil, err
	}
	nextRunTime := nextSchedule.NextRun()
	if nextRunTime.Before(time.Now()) {
		return nextSchedule, nil
	}
	waitingDuration := nextRunTime.Sub(time.Now())
	log.Printf("Waiting %v until %v for schedule %s", waitingDuration, nextRunTime, nextSchedule.Name)
	time.Sleep(waitingDuration)
	return nextSchedule, nil
}

func (bin *EventBin) RunIteration() error {
	nextSchedule, err := bin.WaitUntilNextSchedule()
	if err != nil {
		return err
	}
	if nextSchedule.PreCondition != nil {
		if err := nextSchedule.PreCondition(); err != nil {
			log.Printf("schedule '%s' aborted due pre-condition exception %v", nextSchedule.Name, err)
			nextSchedule.LastError = err
			return nil
		}
	}

	log.Printf("Running schedule %s", nextSchedule.Name)
	err = nextSchedule.Run()
	if err != nil {
		log.Printf("Error running schedule %s: %v", nextSchedule.Name, err)
	}
	nextSchedule.LastError = err
	return err
}

func (bin *EventBin) Run() {
	if len(bin.Schedules) == 0 {
		log.Fatalf("There are no events in event bin")
	}
	for name, f := range StartupEvents {
		log.Printf("Startup %s", name)
		err := f()
		if err != nil {
			log.Fatalf("Failed to run startup event %s = %v", name, err)
		}
	}
	tstart := time.Now().Round(bin.RunInterval)
	if tstart.Before(time.Now()) {
		tstart = tstart.Add(bin.RunInterval)
	}
	waitStart := tstart.Sub(time.Now())
	if waitStart > 0 && !bin.CanStartImediately {
		log.Printf("Waiting to start @ %v", tstart)
		time.Sleep(waitStart)
	}
	log.Printf("Running events every %v", bin.RunInterval)
	for {
		bin.PrintSchedules()
		t0 := time.Now()
		for {
			if err := bin.RunIteration(); err != nil {
				break
			}
		}
		t1 := bin.RunInterval - time.Now().Sub(t0)
		if t1 > 0 {
			log.Printf("Waiting %v for next iteration", t1)
			time.Sleep(bin.RunInterval)
		}

	}
}

func (bin *EventBin) PrintSchedules() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Name", "Last Run", "Next Run", "Last error"})
	for index, schedule := range bin.Schedules {
		var err = ""
		if schedule.LastError != nil {
			err = schedule.LastError.Error()
		}
		t.AppendRows([]table.Row{
			{index, schedule.Name, schedule.LastRun, schedule.NextRun(), err},
		})
	}
	t.Render()
}
