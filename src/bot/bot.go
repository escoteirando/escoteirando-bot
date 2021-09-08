package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"log"
	"time"
)

type Bot struct {
	Instance                *tgbotapi.BotAPI
	Self                    tgbotapi.User
	User                    tgbotapi.User
	MessageChannel          chan domain.TelegramMessage
	AutoDestructChannel     chan domain.TelegramMessage
	UpdateChatOffsetFunc    func(int)
	LastSendTime            time.Time
	IntervalBetweenMessages time.Duration
	MaxRetrySend            int
	AdminChatId             int64
}

var (
	bot Bot
)

func GetCurrentBot() *Bot {
	if bot.Instance == nil {
		log.Fatalf("Bot was not instantiated")
	}
	return &bot
}

func GetCurrentBotInstance() *tgbotapi.BotAPI {
	return GetCurrentBot().Instance
}

func CreateBot(botToken string, debug bool, adminChatId int64) Bot {
	bot = Bot{}
	instance, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to get bot %v", err)
	}
	instance.Debug = debug
	bot.Instance = instance
	bot.Self = bot.Instance.Self
	bot.User, err = instance.GetMe()
	if err != nil {
		log.Fatalf("Failed to get bot user %v", err)
	}
	bot.MessageChannel = make(chan domain.TelegramMessage, 10)
	bot.AutoDestructChannel = make(chan domain.TelegramMessage, 10)
	bot.IntervalBetweenMessages = time.Second
	bot.MaxRetrySend = 10
	bot.AdminChatId = adminChatId
	log.Printf("Bot Instance %s", bot.Self.UserName)
	return bot
}

func (bot *Bot) Start(chatOffset int, onMessage OnMessageFunc, onCallBackQuery OnMessageFunc) {
	go bot.SenderStart()
	go bot.AutoDestructWorkerStart()

	go bot.StartListening(chatOffset, onMessage, onCallBackQuery)
}
