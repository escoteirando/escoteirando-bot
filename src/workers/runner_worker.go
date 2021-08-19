package workers

import (
	"github.com/guionardo/escoteirando-bot/src/services"
	"log"
	"sync"
	"time"
)

func RunnerWorker(wg *sync.WaitGroup) {
	log.Printf("Started Runner Worker")
	defer wg.Done()
	services.RunnerStartRun()
	totalRun := services.RunnerRunningTime(false)
	log.Printf("Total running time until now: %v", totalRun)
	for {
		services.RunnerSaveStatus()
		time.Sleep(time.Duration(1) * time.Minute)
	}
}
