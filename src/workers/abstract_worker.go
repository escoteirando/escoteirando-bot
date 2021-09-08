package workers

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"log"
	"sync"
	"time"
)

type (
	AbstractWorker interface {
	}
	WorkerRunner func(worker *CommonWorker) error
	CommonWorker struct {
		name             string
		interval         time.Duration
		roundStart       bool
		runner           WorkerRunner
		beforeStart      func()
		channel          chan domain.TelegramMessage
		canStartRightNow bool
		waitForStart     func() bool
		dependences      []string
	}
)

var (
	workers        []*CommonWorker
	wg             sync.WaitGroup
	runningWorkers map[string]bool
)

func init() {
	workers = make([]*CommonWorker, 0)
	runningWorkers = make(map[string]bool)
}

func WorkerIsRunning(workerName string) bool {
	_, ok := runningWorkers[workerName]
	return ok
}

func (worker *CommonWorker) ToString() string {
	intervalRun := "Continous"
	if worker.interval > 0 {
		intervalRun = worker.interval.String()
	}
	return fmt.Sprintf("Worker %s (%s)", worker.name, intervalRun)
}

func NewWorker(name string, runner WorkerRunner) *CommonWorker {
	worker := CommonWorker{
		name:             name,
		roundStart:       true,
		runner:           runner,
		interval:         time.Duration(0),
		canStartRightNow: false,
		dependences:      make([]string, 0),
	}

	return &worker
}

func (worker *CommonWorker) Create() {
	wg.Add(1)
	workers = append(workers, worker)
	log.Printf("Created %s", worker.ToString())
}

func (worker *CommonWorker) Interval(interval time.Duration) *CommonWorker {
	worker.interval = interval
	return worker
}

func (worker *CommonWorker) DependsOn(dependences ...string) *CommonWorker {
	worker.dependences = dependences
	return worker
}

func (worker *CommonWorker) NoRoundStart() *CommonWorker {
	worker.roundStart = false
	return worker
}

func (worker *CommonWorker) CanStartRightNow() *CommonWorker {
	worker.canStartRightNow = true
	return worker
}
func (worker *CommonWorker) Channel(channel chan domain.TelegramMessage) *CommonWorker {
	worker.channel = channel
	return worker
}

func (worker *CommonWorker) BeforeStart(beforeStart func()) *CommonWorker {
	worker.beforeStart = beforeStart
	return worker
}

func (worker *CommonWorker) WaitForStart(startWhen func() bool) *CommonWorker {
	worker.waitForStart = startWhen
	return worker
}

func RunAllWorkers() {
	for _, worker := range workers {
		go worker.Run()
	}
	wg.Wait()
}

func (worker *CommonWorker) waitForDependences() {
	for _, dependence := range worker.dependences {
		if !WorkerIsRunning(dependence) {
			log.Printf("Waiting for start %s after dependency %s", worker.ToString(), dependence)
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}

func (worker *CommonWorker) waitForStartFunc() {
	if worker.waitForStart != nil {
		for {
			if worker.waitForStart() {
				break
			}
			log.Printf("Waiting for start %s", worker.ToString())
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}

func (worker *CommonWorker) Run() {
	worker.waitForDependences()
	worker.waitForStartFunc()

	log.Printf("Worker started: %s", worker.ToString())
	defer wg.Done()
	if worker.roundStart && worker.interval > 0 && !worker.canStartRightNow {
		utils.WaitRound(worker.interval, worker.ToString())
	}
	if worker.beforeStart != nil {
		worker.beforeStart()
	}
	var err error
	if worker.interval.Seconds() == 0 {
		err = worker.runner(worker)
	} else {
		for {
			startedAt := time.Now()
			err = worker.runner(worker)
			if err != nil {
				log.Printf("Stopped worker %s due to error %v", worker.name, err)
				break
			}
			if worker.roundStart {
				utils.WaitUntil(utils.RoundStart(startedAt, worker.interval))
			} else {
				utils.WaitUntil(startedAt.Add(worker.interval))
			}
		}
	}
}
