package workers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/services"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	lastSendMessage = time.Now()
	minInterval, _  = time.ParseDuration("1s")
)

const (
	maxRetrySend = 10
)

func waitToSend() {
	intervalLastSend := time.Now().Sub(lastSendMessage)
	if intervalLastSend > minInterval {
		return
	}
	wait := minInterval - time.Now().Sub(lastSendMessage)
	log.Printf("Waiting %v for next message", wait)
	time.Sleep(wait)
}
func MessageSenderWorker(wg *sync.WaitGroup, messages chan domain.TelegramMessage) {
	log.Print("Started Message Sender Worker")
	bot, _ := services.GetBot()

	defer wg.Done()
	for {
		message, ok := <-messages
		if !ok {
			break
		}
		response := sendMessage(bot, &message)
		message.ResponseChannel <- response
	}
	log.Printf("Stopped Message Sender Worker")
}

const (
	continueSending string = "ok"
	cancelSending   string = "cancel"
)

func processSendError(err error) string {
	if err == nil {
		return continueSending
	}
	if strings.Contains(err.Error(), "Too Many Requests") {
		words := strings.Split(err.Error(), " ")
		if len(words) > 0 {
			delay, err := strconv.Atoi(words[len(words)-1])
			if err == nil {
				log.Printf("Too many requests. Waiting %d seconds", delay)
				time.Sleep(time.Duration(delay) * time.Second)
				return continueSending
			}
		}
	}
	log.Printf("Unexpected sending error: %v", err)
	return cancelSending
}
func sendMessage(bot *tgbotapi.BotAPI, message *domain.TelegramMessage) tgbotapi.Message {
	message.RetryCount++
	var msg tgbotapi.Message
	if message.RetryCount > maxRetrySend {
		log.Print("Exceed sending retry")
		return msg
	}
	waitToSend()
	msg, err := bot.Send(message.Message)
	lastSendMessage = time.Now()
	if err == nil {
		return msg
	}
	whatToDo := processSendError(err)
	switch whatToDo {
	case continueSending:
		return sendMessage(bot, message)
	}

	return msg
}
