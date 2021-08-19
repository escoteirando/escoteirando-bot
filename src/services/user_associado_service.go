package services

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/consts"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/repository"
)

func MappaVinculaAssociado(ctx domain.Context, message *tgbotapi.Message, codigoAssociado string) {
	associado, err := repository.FindAssociado(codigoAssociado)
	if err != nil {
		SendMessage(ctx.ChatId, fmt.Sprintf("%s O código de associado não foi encontrado", consts.Warning), message.MessageID)
		return
	}
	err = repository.VincularAssociado(associado, message.From)
	if err != nil {
		SendMessage(ctx.ChatId, fmt.Sprintf("%s Não foi possível vincular o associado %v", consts.Warning, err), message.MessageID)
	} else {
		SendMessage(ctx.ChatId, fmt.Sprintf("%s O associado %s foi vinculado ao usuário %s", consts.ThumsUp, associado.ToString(), message.From.String()), message.MessageID)
	}
	//if err!=nil{
	//	codAssociado,err:=strconv.Atoi(codigoAssociado[0:len(codigoAssociado)-1])
	//	if err!=nil{
	//		SendMessage(ctx.ChatId,fmt.Sprintf("%s Erro na identificação do código do associado %v",consts.Warning,err),message.MessageID)
	//		return
	//	}
	//	msg,err := SendMessage(ctx.ChatId,fmt.Sprintf("%s Localizando associado...",consts.HourglassNotDone),message.MessageID)
	//	mappaAssociado,err:=MappaGetAssociado(ctx,codAssociado)
	//	if err==nil{
	//		EditMessage(ctx.ChatId,msg.MessageID,fmt.Sprintf("%s Associado localizado: %s",consts.ThumsUp,mappaAssociado.Nome))
	//	}else{
	//		EditMessage(ctx.ChatId,msg.MessageID,fmt.Sprintf("%s O código de associado não foi encontrado",consts.Warning))
	//		return
	//	}
	//	_,err=MappaSaveAssociadoEquipe(ctx,mappaAssociado,*message,false)
	//	if err==nil{
	//		EditMessage(ctx.ChatId,msg.MessageID,fmt.Sprintf("%s Associado %s foi vinculado ao usuário %s",consts.ThumsUp,message.From.String()))
	//	}else{
	//		EditMessage(ctx.ChatId,msg.MessageID,fmt.Sprintf("%s Erro ao gravar o associado  %s  %s",consts.ThumbsDown,associado.Nome,err))
	//	}
	//
	//}
	//TODO: Vincular o associado ao contato

}
