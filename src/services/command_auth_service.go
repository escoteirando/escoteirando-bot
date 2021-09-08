package services

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	bot2 "github.com/guionardo/escoteirando-bot/src/bot"
	"github.com/guionardo/escoteirando-bot/src/consts"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"log"
	"net/url"
	"strings"
	"time"
)

func CommandAuth(ctx domain.Context, message *tgbotapi.Message) {
	if len(message.CommandArguments()) > 0 && strings.ToUpper(message.CommandArguments()) != "FORCE" {
		processAuthReturn(ctx, message)
	} else {
		SendAuthCommandMessage(ctx, 0, false)
	}
}

func processAuthReturn(ctx domain.Context, message *tgbotapi.Message) {
	var authKey domain.AuthKey
	sDec, err := b64.StdEncoding.DecodeString(message.CommandArguments())
	if err != nil {
		err = fmt.Errorf("chave de autenticação mal-formada (base64)")
	} else {
		log.Printf("Auth: %v", string(sDec))
		err = json.Unmarshal(sDec, &authKey)
		if err != nil {
			err = fmt.Errorf("chave de autenticação mal-formada (json)")
		} else {
			log.Printf("Auth: %v", authKey)

			if !SetContextSetup(ctx, authKey) {
				err = fmt.Errorf("falha na gravação do contexto do canal de comunicação")
			} else {
				if message.From.ID == BotUser.ID {
					DeleteTextMessage(ctx.ChatId, message.MessageID)
				}

				bot2.GetCurrentBot().SendTextReply(message.Chat.ID, consts.OK+" Contexto atualizado com chave de autorização", message.MessageID)
				SetContextSetup(ctx, authKey)
				return
			}
		}
	}
	if err != nil {
		log.Printf("processAuthReturn: %v", err)
	}

	bot2.GetCurrentBot().SendTextReply(message.Chat.ID, "Chave de autenticação inválida!", message.MessageID)
}

func SendAuthCommandMessageFromWorker(chatId int64) {
	env := utils.GetEnvironmentSetup()
	ctxFrontend := domain.FrontendContext{
		CId: chatId,
		MId: 0,
	}
	jsonFrontend, _ := json.Marshal(ctxFrontend)
	sEnc := b64.StdEncoding.EncodeToString(jsonFrontend)
	encodedUrl := fmt.Sprintf("%s/%s", env.FrontEndUrl, sEnc)
	decodedUrl, _ := url.QueryUnescape(encodedUrl)

	if msg, err := SendButtonMessage(chatId,
		"ADMINISTRADOR: Acesse o link a seguir para autorizar a conexão com os dados da seção",
		"Autorização",
		decodedUrl); err == nil {
		AutoDestruct(msg, time.Duration(20)*time.Minute)
	}
}

func SendAuthCommandMessage(ctx domain.Context, messageId int, force bool) {
	if len(ctx.AuthKey) > 0 && ctx.AuthValidUntil.After(time.Now()) && !force {
		if msg:=bot2.GetCurrentBot().SendTextMessage(ctx.ChatId, fmt.Sprintf("Autorização mAPPA atual ainda é válida até %v\nUse o comando <b>/auth force</b> para forçar a reautenticação", ctx.AuthValidUntil)); msg.MessageID>0 {
			bot2.GetCurrentBot().AutoDestructMessage(msg,time.Duration(30)*time.Second)
		}
		return
	}
	env := utils.GetEnvironmentSetup()
	ctxFrontend := domain.FrontendContext{
		CId: ctx.ChatId,
		MId: messageId,
	}
	jsonFrontend, _ := json.Marshal(ctxFrontend)
	sEnc := b64.StdEncoding.EncodeToString(jsonFrontend)
	encodedUrl := fmt.Sprintf("%s/%s", env.FrontEndUrl, sEnc)
	decodedUrl, _ := url.QueryUnescape(encodedUrl)

	if msg, err := SendButtonMessage(ctx.ChatId,
		"ADMINISTRADOR: Acesse o link a seguir para autorizar a conexão com os dados da seção",
		"Autorização",
		decodedUrl); err == nil {

		AutoDestruct(msg, time.Duration(20)*time.Second)
	}
}
