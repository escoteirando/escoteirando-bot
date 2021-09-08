package events

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/dto"
	"github.com/guionardo/escoteirando-bot/src/eventbin"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"github.com/guionardo/escoteirando-bot/src/mappa_api"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/services"
	"log"
	"time"
)

var syncedSections map[int]interface{}

func init() {
	eventbin.RegisterSchedule(eventbin.
		CreateSchedule("MappaSyncChats", time.Duration(24)*time.Hour, MappaSyncChats, ContextNeeded).
		RoundStart().
		CanStartNow())

	eventbin.RegisterSchedule(eventbin.
		CreateSchedule("MappaSyncProgressoes", time.Duration(24*7)*time.Hour, MappaSyncProgressoes, ContextNeeded).
		RoundStart().
		CanStartNow())
}

func ContextNeeded() error {
	_, err := services.GetGenericContext()
	return err
}
func MappaSyncProgressoes(eventbin.Schedule) error {
	ctx, err := services.GetGenericContext()
	if err == nil {
		progressoes, err := mappa_api.MappaGetProgressoes(ctx)
		if err == nil {
			err = repository.SaveProgressoes(progressoes)
		}
	}
	if err != nil {
		log.Printf("Error syncing progressoes: %v", err)
	}
	return err
}

func MappaSyncChats(eventbin.Schedule) error {
	syncedSections = make(map[int]interface{})

	log.Printf("MappaSyncChats")
	chats, err := repository.GetAllChats()
	if err != nil {
		return err
	}

	for _, chat := range chats {
		err = mappaSyncChat(chat)
		if err == nil {
			log.Printf("Chat synced: #%d/%d", chat.ID, chat.CodSecao)
			syncedSections[chat.CodSecao] = true
		} else {
			log.Printf("Error syncing chat: $%d/%d %v", chat.ID, chat.CodSecao, err)
		}
	}
	return nil
}
func mappaSyncChat(chat domain.Chat) error {
	if chat.AuthValidUntil.Before(time.Now()) {
		return fmt.Errorf("chat authorization was valid until %v - NEED NEW KEY", chat.AuthValidUntil)
	}
	log.Printf("Syncing chat %v", chat)
	ctx := services.GetContextFromChat(chat)
	_, escotista, err := mappaSyncEscotista(ctx, chat.MappaUserId)
	if err != nil {
		return err
	}
	_, _, err = mappaSyncGrupos(ctx, escotista)
	if err != nil {
		return err
	}
	mappaSecoes, _, err := mappaSyncSecoes(ctx, chat.MappaUserId)
	if err != nil {
		return err
	}

	if err = mappaSyncEquipes(ctx, chat.MappaUserId, mappaSecoes); err != nil {
		return err
	}

	return err
}

func mappaSyncEscotista(ctx domain.Context, userId int) (mappa.Escotista, domain.MappaEscotista, error) {
	mappaEscotista, err := mappa_api.MappaGetEscotista(ctx, userId)
	if err != nil {
		log.Printf("Failed to fetch mappaEscotista@mappa %v", err)
		return mappaEscotista, domain.MappaEscotista{}, err
	}
	log.Printf("Fetched mappaEscotista %v", mappaEscotista)
	escotista := dto.CreateEscotista(mappaEscotista)
	err = repository.SaveEscotista(escotista)
	if err != nil {
		log.Printf("Failed to save escotista@db #%d:%s %v", escotista.ID, escotista.UserName, err)
	} else {
		log.Printf("Saved escotista %v", escotista)
	}
	return mappaEscotista, escotista, err
}

func mappaSyncGrupos(ctx domain.Context, escotista domain.MappaEscotista) ([]mappa.Grupo, []domain.MappaGrupo, error) {
	mappaGrupos, err := mappa_api.MappaGetGrupo(ctx, escotista)
	if err != nil {
		log.Printf("Failed to fetch mappaGrupos@mappa %v", err)
		return mappaGrupos, nil, err
	}
	log.Printf("Fetched mappaGrupos %v", mappaGrupos)
	grupos := make([]domain.MappaGrupo, len(mappaGrupos))
	for index, mappaGrupo := range mappaGrupos {
		grupos[index] = dto.CreateGrupo(mappaGrupo)
		err = repository.SaveGrupo(grupos[index])
		if err != nil {
			log.Printf("Failed to save grupo@db #%d:%s %v", grupos[index].ID, grupos[index].Nome, err)
			break
		} else {
			log.Printf("Saved grupo %v", grupos[index])
		}
	}
	return mappaGrupos, grupos, err
}

func mappaSyncSecoes(ctx domain.Context, userId int) ([]mappa.Secao, []domain.MappaSecao, error) {
	mappaSecoes, err := mappa_api.MappaGetSecoes(ctx, userId)
	if err != nil {
		log.Printf("Failed to fetch secoes@mappa %v", err)
		return mappaSecoes, nil, err
	}
	log.Printf("Fetched mappaSecoes %v", mappaSecoes)
	secoes := make([]domain.MappaSecao, len(mappaSecoes))
	for index, mappaSecao := range mappaSecoes {
		secoes[index] = dto.CreateSecao(mappaSecao)
		err = repository.SaveSecao(secoes[index])
		if err != nil {
			log.Printf("Failed to save secao@db #%d:%s %v", secoes[index].ID, secoes[index].Nome, err)
			break
		} else {
			log.Printf("Saved secao %v", secoes[index])
		}
	}
	return mappaSecoes, secoes, err
}

func mappaSyncEquipes(ctx domain.Context, userId int, mappaSecoes []mappa.Secao) error {
	for _, secao := range mappaSecoes {
		_, _, err := mappaSyncEquipe(ctx, userId, secao.Codigo)
		if err != nil {
			return err
		}
		desde := repository.GetUltimaConsultaMarcacao(secao.Codigo)
		if err = mappaSyncMarcacoes(ctx, secao.Codigo, desde); err != nil {
			return err
		}
		repository.SetUltimaConsultaMarcacao(secao.Codigo, time.Now())
	}
	return nil
}

func mappaSyncMarcacoes(ctx domain.Context, codSecao int, desde time.Time) error {
	newCtx := ctx
	newCtx.CodSecao = codSecao
	mappaMarcacoes, err := mappa_api.MappaGetMarcacoes(newCtx, desde)
	if err != nil {
		return err
	}
	err = repository.SetMarcacoes(mappaMarcacoes)
	return err
}

func mappaSyncEquipe(ctx domain.Context, userId int, codSecao int) ([]mappa.SubSecao, []domain.MappaSubSecao, error) {
	if _, ok := syncedSections[codSecao]; ok {
		log.Printf("Section #%d just synced", codSecao)
		return nil, nil, nil
	}
	mappaEquipe, err := mappa_api.MappaGetEquipe(ctx, userId, codSecao)
	if err != nil {
		log.Printf("Failed to fetch mappaEquipe@mappa %v", err)
		return mappaEquipe, nil, err
	}
	log.Printf("Fetched equipe %v", mappaEquipe)
	subSecoes := make([]domain.MappaSubSecao, len(mappaEquipe))
	for index, mappaSubsecao := range mappaEquipe {
		subSecoes[index] = dto.CreateSubsecao(mappaSubsecao)
		err = repository.SaveSubSecao(subSecoes[index])
		if err != nil {
			log.Printf("Failed to save subsecao@db #%d:%s %v", subSecoes[index].ID, subSecoes[index].Nome, err)
			break
		} else {
			log.Printf("Saved subsecao %v", subSecoes[index])
		}
		_, err := mappaSyncSubSecao(mappaSubsecao)
		if err != nil {
			break
		}
	}
	syncedSections[codSecao] = true
	return mappaEquipe, subSecoes, nil
}

func mappaSyncSubSecao(subsecao mappa.SubSecao) ([]domain.MappaAssociado, error) {
	associados := make([]domain.MappaAssociado, len(subsecao.Associados))
	for index, mappaAssociado := range subsecao.Associados {
		associados[index] = dto.CreateAssociado(mappaAssociado)
		err := repository.SaveAssociado(associados[index])
		if err != nil {
			log.Printf("Failed to save associado@db #%d:%s %v", mappaAssociado.Codigo, mappaAssociado.Nome, err)
			return nil, err
		}
		log.Printf("Saved associado %v", associados[index])
	}
	return associados, nil
}
