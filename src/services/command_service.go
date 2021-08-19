package services

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
)

type Command struct {
	Name         string
	Method       func(ctx domain.Context, message *tgbotapi.Message)
	JustForAdmin bool
	Description  string
}

var (
	commands = make(map[string]Command, 0)
)

func RegisterCommand(command Command) {
	commands[command.Name] = command
}

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonSwitch("2sw", "open 2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

func Ask(caption string, command string, options ...string) tgbotapi.InlineKeyboardMarkup {
	var optionsButtons []tgbotapi.InlineKeyboardButton
	for _, option := range options {
		optionsButtons = append(optionsButtons, tgbotapi.NewInlineKeyboardButtonData(option, command+" "+option))
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonSwitch(caption, "SimNão")),
		tgbotapi.NewInlineKeyboardRow(optionsButtons...))
}

func ParseCommand(ctx domain.Context, message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if command, ok := commands[message.Command()]; ok {
		if command.JustForAdmin {

		}
		command.Method(ctx, message)
		return
	}
	switch message.Command() {
	//case "setup":
	//	env := utils.GetEnvironmentSetup()
	//	ctxFrontend := domain.FrontendContext{
	//		CId: message.Chat.ID,
	//		MId: message.MessageID,
	//	}
	//	jsonFrontend, _ := json.Marshal(ctxFrontend)
	//	sEnc := b64.StdEncoding.EncodeToString(jsonFrontend)
	//	encodedUrl := fmt.Sprintf("%s/%s", env.FrontEndUrl, sEnc)
	//	decodedUrl, _ := url.QueryUnescape(encodedUrl)
	//
	//	SendCommandButton(message.Chat.ID, "ADMINISTRADOR: Acesse o link a seguir para autorizar a conexão com os dados da seção", "Autorização", decodedUrl)
	//
	//	break
	//case "auth":
	//	var authKey domain.AuthKey
	//	sDec, err := b64.StdEncoding.DecodeString(message.CommandArguments())
	//	if err == nil {
	//		log.Printf("Auth: %v", string(sDec))
	//		err = json.Unmarshal(sDec, &authKey)
	//		if err == nil {
	//			log.Printf("Auth: %v", authKey)
	//			msg, err := SendMessage(message.Chat.ID, consts.OK+"Chave 🆗! Iniciando registro", message.MessageID)
	//			if err == nil {
	//				if SetContextSetup(ctx, authKey) {
	//					EditMessage(message.Chat.ID, msg.MessageID, consts.OK+" Contexto atualizado com chave de autorização")
	//					err = MappaSetupAuth(&ctx)
	//					if err == nil {
	//						SetContextSetup(ctx, authKey)
	//					}
	//					return
	//				}
	//			}
	//
	//		}
	//	}
	//	log.Printf("Erro: %v", err)
	//
	//	SendMessage(message.Chat.ID, "Chave de autenticação inválida!", message.MessageID)
	//
	//	break
	case "associado":
		codAssociado := message.CommandArguments()
		MappaVinculaAssociado(ctx, message, codAssociado)

		break
	case "test":
		SendCommandButton(message.Chat.ID, "Teste de botão", "Clique aqui", "https://google.com")

		//msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
		//
		//
		//	msg.ReplyMarkup = Ask("Pergunta","ask","Sim","Não")
		//
		//
		//bot.Send(msg)
		break

	case "marcacoes":
		msg, err := SendTextMessage(message.Chat.ID, "Obtendo marcações da seção", message.MessageID)
		if err == nil {
			marcacoes, err := MappaGetMarcacoes(ctx)
			if err == nil {
				EditMessage(message.Chat.ID, msg.MessageID, fmt.Sprintf("Marcações: %d", len(marcacoes)))
			} else {
				EditMessage(message.Chat.ID, msg.MessageID, fmt.Sprintf("Erro: %v", err))
			}
		}
		break

	}

}
