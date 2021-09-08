package services

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"log"
	"time"
)

func ScheduleGet(id string, interval time.Duration, roundStart bool) (domain.Schedule, bool) {
	var schedule domain.Schedule
	result := repository.GetDB().First(&schedule, "id = ?", id)
	if result.Error == nil && result.RowsAffected > 0 {
		return schedule, schedule.NextRun.Before(time.Now())
	}
	repository.CanIWrite()
	nextRun := time.Now()
	if roundStart {
		nextRun = utils.RoundStart(nextRun, interval).Add(-interval)
	}
	schedule = domain.Schedule{
		Id:         id,
		LastRun:    utils.Epoch,
		NextRun:    nextRun,
		Interval:   interval,
		RoundStart: roundStart,
	}
	repository.GetDB().Save(&schedule)

	repository.YouCanWrite()
	return schedule, schedule.NextRun.Before(time.Now())
}

func ScheduleUpdate(schedule *domain.Schedule) bool {
	lastRun := time.Now()
	nextRun := schedule.LastRun.Add(schedule.Interval)
	if schedule.RoundStart {
		nextRun = utils.RoundStart(lastRun, schedule.Interval)
	}
	schedule.LastRun = lastRun
	schedule.NextRun = nextRun
	if err := repository.GetDB().Save(&schedule).Error; err != nil {
		log.Printf("Failed to add %s - %v", schedule.ToString(), err)
		SendMessageToAdminChat(fmt.Sprintf("Failed to add schedule %s - %v", schedule.ToString(), err))
		return false
	}
	log.Printf("Added %s", schedule.ToString())
	return true
}
