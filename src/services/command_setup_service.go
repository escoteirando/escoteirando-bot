package services

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"net/url"
	"time"
)

func CommandSetup(_ domain.Context, message *tgbotapi.Message) {
	env := utils.GetEnvironmentSetup()
	ctxFrontend := domain.FrontendContext{
		CId: message.Chat.ID,
		MId: message.MessageID,
	}
	jsonFrontend, _ := json.Marshal(ctxFrontend)
	sEnc := b64.StdEncoding.EncodeToString(jsonFrontend)
	encodedUrl := fmt.Sprintf("%s/%s", env.FrontEndUrl, sEnc)
	decodedUrl, _ := url.QueryUnescape(encodedUrl)

	if msg, err := SendButtonMessage(message.Chat.ID,
		"ADMINISTRADOR: Acesse o link a seguir para autorizar a conexão com os dados da seção",
		"Autorização",
		decodedUrl); err == nil {
	
		AutoDestruct(msg, time.Duration(20)*time.Second)
	}

}
