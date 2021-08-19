package services

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"strings"
	"time"
)

type (
	ChatUsers struct {
		ChatId    int
		LastQuery time.Time
		Users     []domain.User
	}

	Chats map[int]ChatUsers
)

var (
	chats = make(Chats, 0)
)

func NewUsersMessage(members *[]tgbotapi.User, chatId int64) {
	if members == nil || len(*members) < 1 {
		return
	}

	var users = make([]string, len(*members))
	for _, member := range *members {
		users = append(users, member.FirstName)
	}

	var msg string
	if len(users) == 1 {
		msg = "Seja bem-vindo(a) <b>%s</b> %s"
	} else {
		msg = "Sejam bem-vindo(a)s: <b>%s</b> %s"
	}
	convite := `

Para ligar o associado (membro juvenil ou adulto) ao seu usuário aqui no grupo, use o comando
/associado REGISTRO_UEB
(número que consta na carteirinha).`
	msg = fmt.Sprintf(msg, strings.Join(users, ", "), convite)
	topic := tgbotapi.NewMessage(chatId, msg)

	bot, _ := GetBot()
	bot.Send(topic)
}
