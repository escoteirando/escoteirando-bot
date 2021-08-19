package services

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"time"
)

func GetContext(message *tgbotapi.Message) domain.Context {
	context := domain.Context{
		CodSecao:    0,
		ChatName:    "",
		ChatId:      0,
		AuthKey:     "",
		MappaUserId: 0,
	}
	db := repository.GetDB()
	var chatModel domain.Chat
	db.First(&chatModel, message.Chat.ID)
	if chatModel.ID != 0 {
		context.CodSecao = chatModel.CodSecao
		context.ChatName = fmt.Sprintf("%s:%s", chatModel.Type, chatModel.Title)
		context.ChatId = message.Chat.ID
		context.AuthKey = chatModel.AuthKey
		context.MappaUserId = chatModel.MappaUserId
		context.CodSecao = chatModel.CodSecao
		context.AuthValidUntil = chatModel.AuthValidUntil
	} else {
		chatModel = domain.Chat{
			ID:            message.Chat.ID,
			Title:         message.Chat.Title,
			Type:          message.Chat.Type,
			AllAdmin:      message.Chat.AllMembersAreAdmins,
			CodSecao:      0,
			LastSetupCall: time.Unix(0, 0),
		}
		db.Create(&chatModel)
	}
	CheckContextSetup(chatModel)
	return context
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

func CheckContextSetup(chat domain.Chat) {
	if 0 == chat.ID {
		return
	}
	if len(chat.AuthKey) == 0 {
		schedule, canRun := ScheduleGet(fmt.Sprintf("context_%d", chat.ID), time.Duration(6)*time.Hour, false)
		if canRun {
			msg, err := SendButtonMessage(chat.ID, "Admin! Faça a configuração deste chat", "Setup", "/setup")
			if err == nil {
				AutoDestruct(msg, time.Duration(30)*time.Second)
			}
			ScheduleUpdate(&schedule)
		}
		return
	}
	if chat.AuthValidUntil.Before(time.Now()) {
		// Autorização fora da validade
		schedule, canRun := ScheduleGet(fmt.Sprintf("auth_%d", chat.ID), time.Duration(15)*time.Minute, true)
		if canRun {
			if ctx, err := GetContextForSection(chat.CodSecao); err == nil {
				SendAuthCommandMessage(ctx, 0, false)
				ScheduleUpdate(&schedule)
			}
		}

		return
	} else {
		authRestante := time.Now().Sub(chat.AuthValidUntil)
		if authRestante.Hours() < 24 {
			schedule, canRun := ScheduleGet(fmt.Sprintf("validAuth_%d", chat.ID), time.Duration(6)*time.Hour, false)
			if canRun {
				SendTextMessage(chat.ID, fmt.Sprintf("Autorização de acesso ao mAPPa válida até: %s", chat.AuthValidUntil.Format(time.RFC822)), 0)
				ScheduleUpdate(&schedule)
			}
		}
	}

}
