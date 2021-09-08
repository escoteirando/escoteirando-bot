package events

import (
	"github.com/guionardo/escoteirando-bot/src/dto"
	"github.com/guionardo/escoteirando-bot/src/eventbin"
	"github.com/guionardo/escoteirando-bot/src/mappa_api"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/services"
	"log"
	"time"
)

const syncEscotistasInterval = time.Duration(24) * time.Hour

func init() {
	eventbin.RegisterSchedule(eventbin.
		CreateSchedule("MappaSyncEscotistas", syncEscotistasInterval, MappaSyncEscotistas, ContextNeeded).
		RoundStart().
		CanStartNow())
}

func MappaSyncEscotistas(eventbin.Schedule) error {
	schedule, canRun := services.ScheduleGet("mappa_sync_escotistas", syncEscotistasInterval, true)
	if !canRun {
		return nil
	}
	escotistas, err := repository.GetAllEscotistasIds()
	if err != nil {
		return err
	}
	ctx, err := services.GetGenericContext()
	if err != nil {
		return err
	}
	for _, eid := range escotistas {
		escotista, err := repository.GetEscotista(eid)
		if err == nil {
			continue
		}

		mappaEscotista, err := mappa_api.MappaGetEscotista(ctx, eid)
		if err == nil && mappaEscotista.Codigo > 0 {
			log.Printf("Synced escotista %v", mappaEscotista)
			escotista = dto.CreateEscotista(mappaEscotista)
			err = repository.SaveEscotista(escotista)
			if err != nil {
				return err
			}
			associado, err := repository.GetAssociado(escotista.CodAssociado)
			if err != nil {
				mappaAssociado, err := mappa_api.MappaGetAssociado(ctx, escotista.CodAssociado)
				if err == nil && mappaAssociado.Codigo > 0 {
					associado = dto.CreateAssociado(mappaAssociado)
					err = repository.SaveAssociado(associado)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	services.ScheduleUpdate(&schedule)
	return nil
}
