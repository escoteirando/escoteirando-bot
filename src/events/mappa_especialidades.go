package events

import (
	"github.com/guionardo/escoteirando-bot/src/eventbin"
	"github.com/guionardo/escoteirando-bot/src/mappa_api"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/services"
	"time"
)

const interval = time.Duration(24*7) * time.Hour

func init() {
	eventbin.RegisterSchedule(eventbin.
		CreateSchedule("MappaFetchEspecialidades", interval, MappaFetchEspecialidades, ContextNeeded).
		RoundStart().
		CanStartNow())
}

func MappaFetchEspecialidades(eventbin.Schedule) error {
	schedule, canRun := services.ScheduleGet("fetch_especialidades", interval, true)
	if !canRun {
		return nil
	}
	ctx, err := services.GetGenericContext()
	if err != nil {
		return err
	}
	especialidades, err := mappa_api.MappaGetEspecialidades(ctx)
	if err != nil {
		return err
	}
	err = repository.SaveEspecialidades(especialidades)
	services.ScheduleUpdate(&schedule)
	return err
}
