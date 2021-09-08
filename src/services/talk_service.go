package services

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	bot2 "github.com/guionardo/escoteirando-bot/src/bot"
	"github.com/guionardo/escoteirando-bot/src/consts"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"log"
	"strings"
	"time"
)

func ParseUserConversation(ctx domain.Context, message *tgbotapi.Message, update tgbotapi.Update) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	talkType := TryGetTalkType(message.Text)
	switch talkType {
	case consts.IDontKnow:
		responseIDontKnow(ctx, update, message)
		break
	case consts.TalkTypeGoodDay:
		responseGoodDay(ctx, update, message)
		break
	}
	//bot2.GetCurrentBot().SendTextReply(ctx.ChatId, message.Text, message.MessageID)

}

func responseGoodDay(ctx domain.Context, update tgbotapi.Update, message *tgbotapi.Message) {
	lastTalk, err := repository.GetLastTalk(consts.TalkTypeGoodDay, message.From.ID,ctx.ChatId)
	if err == nil && utils.SameDay(lastTalk.CreatedAt, time.Now()) {
		// Já deu bom-dia, boa-tarde ou boa-noite para este usuário hoje
		log.Printf("GoodDay %s", message.From.UserName)
		return
	}
	resposta := fmt.Sprintf("%s %s", GreetingTime(), message.From.FirstName)
	msg := bot2.GetCurrentBot().SendTextReply(ctx.ChatId, resposta, message.MessageID)
	if msg.MessageID != 0 {
		repository.AddTalk(consts.TalkTypeGoodDay, message.From.ID, message.From.UserName, ctx.ChatId, message.MessageID)
	}
}

func responseIDontKnow(ctx domain.Context, update tgbotapi.Update, message *tgbotapi.Message) {
	botName := bot2.GetCurrentBotInstance().Self.UserName
	if strings.Contains(message.Text, fmt.Sprintf("@%s", botName)) {
		bot2.GetCurrentBot().
			SendTextReply(ctx.ChatId,
				fmt.Sprintf("Desculpe, %s. Não entendi", message.From.FirstName), message.MessageID)
	}
	log.Printf("Não entendi: %s", message.Text)
}

func TryGetTalkType(text string) int {
	text = strings.ToLower(consts.RemoverAcentos(text))
	if text == "bom dia" || text == "boa tarde" || text == "boa noite" {
		return consts.TalkTypeGoodDay
	}

	return consts.IDontKnow
	//doc, err := prose.NewDocument(text)
	//doc.

}

func GreetingTime() string {
	if time.Now().Hour() < 12 {
		return consts.Morning + " Bom dia"
	}
	if time.Now().Hour() < 18 {
		return consts.Afternoon + " Boa tarde"
	}
	return consts.Night + " Boa noite"
}
