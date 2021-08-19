package workers

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/services"
	"log"
	"sync"
	"time"
)

func AutoDescructMessageWorker(wg *sync.WaitGroup, messages chan domain.TelegramMessage) {
	log.Print("Started AutoDescruct Message Worker")
	defer wg.Done()
	for {
		message, ok := <-messages
		if !ok {
			break
		}
		go func() {
			if message.AutoDestruct.After(time.Now()) {
				time.Sleep(message.AutoDestruct.Sub(time.Now()))
				services.DeleteMessage(message.ChannelId, message.MessageId)
				log.Printf("Destructed message %d/%d", message.ChannelId, message.MessageId)
			}
		}()
		time.Sleep(time.Duration(5) * time.Second)
	}
	log.Printf("End of autodestruct")
}
