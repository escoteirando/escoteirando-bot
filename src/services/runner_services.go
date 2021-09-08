package services

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"log"
	"time"
)

var currentRunId uint
var lastRun domain.Run

func RunnerStartRun() {
	// Obter a última execução
	var lr domain.Run
	result := repository.GetDB().Order("id desc").Find(&lr).Limit(1)
	if result.Error == nil {
		lastRun = lr
	}
	run := domain.Run{
		Start: time.Now(),
		End:   time.Now(),
	}
	repository.CanIWrite()
	result = repository.GetDB().Create(&run)
	repository.YouCanWrite()
	if result.Error == nil {
		currentRunId = run.ID
		log.Printf("Started @ %v", run.Start)
	}
}

func RunnerSaveStatus() {
	repository.CanIWrite()
	result := repository.GetDB().Model(&domain.Run{}).Where("id = ?", currentRunId).Update("end", time.Now())
	repository.YouCanWrite()
	if result.Error == nil {
		log.Printf("runner status: %v", time.Now())
	}
}

func RunnerLastCompleteRun() (domain.Run, error) {
	if lastRun.ID == 0 {
		return lastRun, fmt.Errorf("no last run")
	}
	return lastRun, nil
}

func RunnerRunningTime(justThisRunning bool) time.Duration {
	if justThisRunning {
		var running domain.Run
		result := repository.GetDB().Find(&running, currentRunId)
		if result.Error == nil {
			return running.End.Sub(running.Start)
		}
	} else {
		var runnings []domain.Run
		result := repository.GetDB().Where("id < ?", currentRunId).Find(&runnings)
		if result.Error == nil {
			var totalRunning time.Duration
			for _, run := range runnings {
				totalRunning += run.End.Sub(run.Start)
			}
			return totalRunning
		}
	}
	return time.Duration(0)
}

func FirstRunningTime() time.Time {
	var running domain.Run
	result := repository.GetDB().First(&running)
	if result.Error == nil && result.RowsAffected > 0 {
		return running.Start
	}
	return utils.Epoch
}
