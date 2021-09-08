package services

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	bot2 "github.com/guionardo/escoteirando-bot/src/bot"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"net/url"
	"time"
)

func GetContext(message *tgbotapi.Message) *domain.Context {
	context, err := repository.GetFromChat(message.Chat.ID)
	if err != nil {
		err = repository.SaveChatFromMessage(message)
		if err == nil {
			context, err = repository.GetFromChat(message.Chat.ID)
		}
	}
	if context != nil {
		CheckContextSetup(*context)
	}
	return context
}

func GetContextFromChat(chat domain.Chat) domain.Context {
	return domain.Context{
		CodSecao:       chat.CodSecao,
		ChatName:       fmt.Sprintf("%s:%s", chat.Type, chat.Title),
		ChatId:         chat.ID,
		AuthKey:        chat.AuthKey,
		MappaUserId:    chat.MappaUserId,
		AuthValidUntil: chat.AuthValidUntil,
	}
}

func GetContextForSection(codSecao int) (domain.Context, error) {
	var chat domain.Chat
	var context domain.Context
	response := repository.GetDB().Where("cod_secao = ?", codSecao).First(&chat)
	if response.Error == nil && response.RowsAffected > 0 {
		if chat.AuthValidUntil.Before(time.Now()) {
			return context, fmt.Errorf("mappa authorization is expired")
		}

		context = domain.Context{
			CodSecao:       chat.CodSecao,
			ChatName:       chat.Title,
			ChatId:         chat.ID,
			AuthKey:        chat.AuthKey,
			MappaUserId:    chat.MappaUserId,
			AuthValidUntil: chat.AuthValidUntil,
		}
		return context, nil
	}
	if response.Error == nil {
		response.Error = fmt.Errorf("secao #%d inexistente", codSecao)

	}
	return context, response.Error
}

func GetGenericContext() (domain.Context, error) {
	var chat domain.Chat
	var context domain.Context
	response := repository.GetDB().Where("auth_valid_until > ?", time.Now()).First(&chat)
	if response.Error == nil && response.RowsAffected > 0 {
		context = domain.Context{
			CodSecao:       chat.CodSecao,
			ChatName:       chat.Title,
			ChatId:         chat.ID,
			AuthKey:        chat.AuthKey,
			MappaUserId:    chat.MappaUserId,
			AuthValidUntil: chat.AuthValidUntil,
		}
		return context, nil
	}
	return context, fmt.Errorf("no generic context available")
}

func SetContextSetup(ctx domain.Context, authKey domain.AuthKey) bool {
	db := repository.GetDB()
	var chatModel domain.Chat
	db.First(&chatModel, ctx.ChatId)

	if chatModel.ID != 0 && chatModel.ID == authKey.CId {
		chatModel.AuthKey = authKey.Id
		chatModel.AuthValidUntil = authKey.Created.Add(time.Second * time.Duration(authKey.Ttl))
		chatModel.MappaUserId = authKey.UserId
		chatModel.CodSecao = ctx.CodSecao
		db.Save(&chatModel)
		SetAuthorizationInCache(chatModel)
		return true
	}
	return false
}

func CheckContextSetup(context domain.Context) {
	if 0 == context.ChatId {
		return
	}
	if len(context.AuthKey) == 0 {
		schedule, canRun := ScheduleGet(fmt.Sprintf("context_%d", context.ChatId), time.Duration(6)*time.Hour, false)
		if canRun {
			msg, err := SendButtonMessage(context.ChatId, "Admin! Faça a configuração deste chat", "Setup", "/setup")
			if err == nil {
				AutoDestruct(msg, time.Duration(30)*time.Second)
			}
			ScheduleUpdate(&schedule)
		}
		return
	}
	if context.AuthValidUntil.Before(time.Now()) {
		// Autorização fora da validade
		schedule, canRun := ScheduleGet(fmt.Sprintf("auth_%d", context.ChatId), time.Duration(15)*time.Minute, true)
		if canRun {
			if ctx, err := GetContextForSection(context.CodSecao); err == nil {
				SendAuthCommandMessage(ctx, 0, false)
				ScheduleUpdate(&schedule)
			}
		}

		return
	} else {
		authRestante := time.Now().Sub(context.AuthValidUntil)
		if authRestante.Hours() < 24 {
			schedule, canRun := ScheduleGet(fmt.Sprintf("validAuth_%d", context.ChatId), time.Duration(6)*time.Hour, false)
			if canRun {
				bot2.GetCurrentBot().SendTextMessage(context.ChatId, fmt.Sprintf("Autorização de acesso ao mAPPa válida até: %s", context.AuthValidUntil.Format(time.RFC822)))
				ScheduleUpdate(&schedule)
			}
		}
	}

}

func GetUserIsAdmin(chat *tgbotapi.Chat, from *tgbotapi.User) bool {
	bot := bot2.GetCurrentBotInstance()
	admins, err := bot.GetChatAdministrators(chat.ChatConfig())
	if err != nil {
		return false
	}
	for _, user := range admins {
		if from.ID == user.User.ID {
			return true
		}
	}
	return false
}

func SendContextAuthRenewMessage(chatId int64) {
	sendContextAuthMessage(chatId,
		"ADMINISTRADOR: Acesse o link a seguir para renovar a autorização",
		fmt.Sprintf("authRenew_%d", chatId),
		time.Duration(1)*time.Hour)
}

func SendContextAuthSetupMessage(chatId int64) {
	sendContextAuthMessage(chatId,
		"ADMINISTRADOR: Acesse o link a seguir para autorizar a conexão com os dados da seção",
		fmt.Sprintf("authSetup_%d", chatId),
		time.Duration(1)*time.Hour)
}

func sendContextAuthMessage(chatId int64, message string, scheduleName string, scheduleInterval time.Duration) {
	schedule, canRun := ScheduleGet(scheduleName, scheduleInterval, true)

	if !canRun {
		return
	}
	env := utils.GetEnvironmentSetup()
	ctxFrontend := domain.FrontendContext{
		CId: chatId,
		MId: 0,
	}
	jsonFrontend, _ := json.Marshal(ctxFrontend)
	sEnc := b64.StdEncoding.EncodeToString(jsonFrontend)
	encodedUrl := fmt.Sprintf("%s/%s", env.FrontEndUrl, sEnc)
	decodedUrl, _ := url.QueryUnescape(encodedUrl)

	if msg, err := bot2.GetCurrentBot().SendButtonMessage(chatId,
		message,
		"Autorização",
		decodedUrl); err == nil {

		AutoDestruct(msg, time.Duration(20)*time.Second)
	}
	ScheduleUpdate(&schedule)
}
