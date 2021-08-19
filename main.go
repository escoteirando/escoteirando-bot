package main

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/services"
	"github.com/guionardo/escoteirando-bot/src/workers"
	"sync"
)

func main() {
	services.Setup()
	messageChannel := make(chan domain.TelegramMessage)
	autoDescructChannel := make(chan domain.TelegramMessage)
	services.SetupMessageChanel(messageChannel)
	services.SetupAutoDestructChannel(autoDescructChannel)
	var wg sync.WaitGroup
	wg.Add(5)
	go workers.Listener(&wg)
	go workers.BackgroundWorker(&wg)
	go workers.MessageSenderWorker(&wg, messageChannel)
	go workers.RunnerWorker(&wg)
	go workers.AutoDescructMessageWorker(&wg, autoDescructChannel)
	wg.Wait()
}
