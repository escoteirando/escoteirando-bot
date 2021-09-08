package services

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	bot2 "github.com/guionardo/escoteirando-bot/src/bot"
	"github.com/guionardo/escoteirando-bot/src/consts"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/utils"
)

func CommandMainAdmin(ctx domain.Context, message *tgbotapi.Message) {
	cBot := bot2.GetCurrentBot()

	if GetAdminChatId() > 0 {
		cBot.SendTextReply(ctx.ChatId, consts.Warning+" Chat de administração geral já está definido", utils.GetMessageId(message))
		return
	}
	if SetAdminChatId(ctx.ChatId) {
		cBot.SendTextReply(ctx.ChatId, fmt.Sprintf(consts.ThumsUp+" Chat de administração = %s",
			utils.ChatDescription(message.Chat)), utils.GetMessageId(message))
		return
	}
}

func CommandNotMainAdmin(ctx domain.Context, message *tgbotapi.Message) {
	cBot := bot2.GetCurrentBot()
	currentAdmin := GetAdminChatId()
	if currentAdmin == 0 {
		cBot.SendTextReply(ctx.ChatId, "Não há um administrador geral", utils.GetMessageId(message))
		return
	}
	if currentAdmin != ctx.ChatId {
		cBot.SendTextReply(ctx.ChatId, consts.ThumbsDown+" Somente o chat de administração geral pode revogar", utils.GetMessageId(message))
		return
	}
	if SetAdminChatId(0) {
		cBot.SendTextReply(ctx.ChatId, fmt.Sprint("Administração geral foi desabilitada"), utils.GetMessageId(message))
	}
}

func SendMessageToAdminChat(message string) {
	bot2.GetCurrentBot().SendAdminTextMessage(message)
}
