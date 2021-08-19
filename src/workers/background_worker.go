package workers

import (
	"github.com/guionardo/escoteirando-bot/src/services"
	"log"
	"sync"
	"time"
)

func BackgroundWorker(wg *sync.WaitGroup) {
	log.Print("Started Background Worker")
	defer wg.Done()

	services.SendGreetings()
	log.Print("Starting background main loop")
	for {
		services.FetchEquipes()
		// Obter marcações das seções
		services.FetchMarcacoes()
		// Enviar marcações pendentes
		services.PublishMarcacoesToChats()

		// Enviar feriados e aniversários
		services.PublishHolidays()
		services.PublishBirthdays()
		time.Sleep(time.Duration(1) * time.Hour)
	}
}
