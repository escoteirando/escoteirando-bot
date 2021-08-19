package workers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/services"
	"log"
	"sync"
)

func Listener(wg *sync.WaitGroup) {
	defer wg.Done()
	bot, err := services.GetBot()

	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(services.GetChatOffset())
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			log.Printf("Callback: %v", update)

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		}
		if update.Message == nil || update.Message.From.UserName == bot.Self.UserName { // ignore any non-Message Updates
			continue
		}

		services.NewUsersMessage(update.Message.NewChatMembers, update.Message.Chat.ID)

		msgId := update.Message.MessageID
		services.SetChatOffset(msgId + 1)

		var ctx = services.GetContext(update.Message)
		if ctx.CodSecao < 1 {
			log.Printf("%s", ctx.ChatName)
		}
		if update.Message.IsCommand() {
			services.ParseCommand(ctx, update.Message, bot)
		} else {
			services.ParseUserConversation(ctx, update.Message)
		}

	}
}
