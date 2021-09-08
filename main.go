package main

import (
	"github.com/guionardo/escoteirando-bot/src/bot"
	"github.com/guionardo/escoteirando-bot/src/eventbin"
	"github.com/guionardo/escoteirando-bot/src/events"
	"github.com/guionardo/escoteirando-bot/src/services"
)

func main() {
	services.Setup()

	//events.MappaSyncEscotistas(eventbin.Schedule{
	//	Name:               "",
	//	LastRun:            time.Time{},
	//	Interval:           0,
	//	WorkFunc:           nil,
	//	LastError:          nil,
	//	RoundStarting:      false,
	//	CanStartImediately: false,
	//})
	//eid, err := repository.GetAllEscotistasIds()
	//if err == nil {
	//	fmt.Printf("Escotistas: %v", eid)
	//}

	botInstance := bot.GetCurrentBot()
	botInstance.Start(services.GetChatOffset(), services.BotOnMessage,services.BotOnCallBackQuery)

	events.InitEvents()

	eventBin := eventbin.CreateEventBin().AddSchedules()

	//var wg sync.WaitGroup
	//wg.Add(6)
	//go workers.Listener(&wg)
	//go workers.BackgroundWorker(&wg)
	//go workers.MessageSenderWorker(&wg, messageChannel)
	//go workers.RunnerWorker(&wg)
	//go workers.AutoDescructMessageWorker(&wg, autoDescructChannel)
	//go workers.StatsWorker(&wg)
	//wg.Wait()
	//go workers.RunAllWorkers()
	//_:=workers.WorkerRunner()

	eventBin.Run()
}
