package services

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/consts"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/rickar/cal/v2"
	"github.com/rickar/cal/v2/aa"
	"log"
	"time"
)

const (
	birthdaySchedule string = "birthday"
	holidaySchedule  string = "holiday"
)

var (
	calendar *cal.BusinessCalendar
	// AnoNovo represents New Year's Day on 1-Jan
	AnoNovo = aa.NewYear.Clone(&cal.Holiday{Name: "Ano Novo", Description: consts.FireWorks, Type: cal.ObservancePublic})

	// Tiradentes represents Tiradentes' Day on 21-Apr
	Tiradentes = &cal.Holiday{
		Name:        "Dia de Tiradentes",
		Description: consts.Brasil,
		Month:       time.April,
		Day:         21,
		Func:        cal.CalcDayOfMonth,
	}

	// Trabalhador represents Labor Day on 21-Apr
	Trabalhador = aa.WorkersDay.Clone(&cal.Holiday{Name: "Dia do Trabalhador", Description: consts.Worker, Type: cal.ObservancePublic})

	// Independencia represents Brazil Independence Day on 07-Sep
	Independencia = &cal.Holiday{
		Name:        "Dia da Independência",
		Description: consts.Brasil,
		Month:       time.September,
		Day:         7,
		Func:        cal.CalcDayOfMonth,
	}

	// NossaSenhoraAparecida represents Our Lady of Aparecida Day - Patroness of Brazil on 12-Oct
	NossaSenhoraAparecida = &cal.Holiday{
		Name:        "Dia de Nossa Senhora Aparecida",
		Description: consts.Cross,
		Month:       time.October,
		Day:         12,
		Func:        cal.CalcDayOfMonth,
	}

	// Finados represents Day of the Dead on 02-Nov
	Finados = &cal.Holiday{
		Name:        "Dia de Finados",
		Description: consts.FlorBranca,
		Month:       time.November,
		Day:         2,
		Func:        cal.CalcDayOfMonth,
	}

	// Republica represents Proclamation of the Republic on 15-Nov
	Republica = &cal.Holiday{
		Name:        "Proclamação da República",
		Description: consts.Brasil,
		Month:       time.November,
		Day:         15,
		Func:        cal.CalcDayOfMonth,
	}

	// CorpusChristi represents Corpus Christi on the 60th day after Easter
	CorpusChristi = aa.CorpusChristi.Clone(&cal.Holiday{Description: consts.Cross})

	// SextaFeiraSanta represents Good Friday - two days before Easter
	SextaFeiraSanta = aa.GoodFriday.Clone(&cal.Holiday{
		Name:        "Sexta-feira Santa",
		Description: consts.Rabbit,
		Type:        cal.ObservancePublic})

	// Carnaval represents Brazilian Carnival - 47 days before Easter
	Carnaval = &cal.Holiday{
		Name:        "Carnaval",
		Description: consts.Clown,
		Type:        cal.ObservancePublic,
		Offset:      -47,
		Func:        cal.CalcEasterOffset,
	}

	// Natal represents Christmas Day on 25-Dec
	Natal = aa.ChristmasDay.Clone(&cal.Holiday{Name: "Natal", Description: consts.ArvoreNatal, Type: cal.ObservancePublic})

	// 23/04 - Dia mundial do escoteiro
	Escoteiro = &cal.Holiday{
		Name:        "Dia Mundial do Escoteiro",
		Description: consts.FlorDeLis,
		Type:        cal.ObservanceOther,
		Month:       4,
		Day:         23,
		Func:        cal.CalcDayOfMonth,
	}

	// 01/08 - Dia mundial do escotismo
	Escotismo = &cal.Holiday{
		Name:        "Dia Mundial do Escotismo",
		Description: consts.FlorDeLis,
		Type:        cal.ObservanceOther,
		Month:       8,
		Day:         1,
		Func:        cal.CalcDayOfMonth,
	}

	// 04/10 - Dia mundial do lobinho
	Lobinho = &cal.Holiday{
		Name:        "Dia Mundial do Lobinho",
		Description: consts.MelhorPossivel,
		Type:        cal.ObservanceOther,
		Month:       10,
		Day:         4,
		Func:        cal.CalcDayOfMonth,
	}

	// 29/06 - Dia do pioneiro
	Pioneiro = &cal.Holiday{
		Name:        "Dia do Pioneiro",
		Description: consts.FlorDeLis,
		Type:        cal.ObservanceOther,
		Month:       6,
		Day:         29,
		Func:        cal.CalcDayOfMonth,
	}

	// 06/08 - Dia do chefe escoteiro
	Chefe = &cal.Holiday{
		Name:        "Dia do Chefe Escoteiro",
		Description: consts.FlorDeLis,
		Type:        cal.ObservanceOther,
		Month:       8,
		Day:         6,
		Func:        cal.CalcDayOfMonth,
	}

	// 18/06 - Dia do Sênior
	Senior = &cal.Holiday{
		Name:        "Dia do Sênior/Guia",
		Description: consts.FlorDeLis,
		Type:        cal.ObservanceOther,
		Month:       6,
		Day:         18,
		Func:        cal.CalcDayOfMonth,
	}

	TesteHoje = &cal.Holiday{
		Name:        "Feriado teste",
		Description: consts.FlorDeLis,
		Type:        cal.ObservanceOther,
		Month:       time.Now().Month(),
		Day:         time.Now().Day(),
		Func:        cal.CalcDayOfMonth,
	}

	// Holidays provides a list of the standard national holidays
	Holidays = []*cal.Holiday{
		AnoNovo,
		Tiradentes,
		Trabalhador,
		Independencia,
		NossaSenhoraAparecida,
		Finados,
		Republica,
		CorpusChristi,
		SextaFeiraSanta,
		Carnaval,
		Natal,
	}
)

func SetupCalendar() {
	calendar = cal.NewBusinessCalendar()
	calendar.Name = "Calendário Escoteiro"
	calendar.AddHoliday(
		TesteHoje,
		AnoNovo,
		Tiradentes,
		Trabalhador,
		Independencia,
		NossaSenhoraAparecida,
		Finados,
		Republica,
		CorpusChristi,
		SextaFeiraSanta,
		Carnaval,
		Natal,
		Escoteiro,
		Escotismo,
		Lobinho,
		Pioneiro,
		Chefe,
		Senior,
	)
}

func GetHolidaysFromDay(day time.Time) []*cal.Holiday {
	var holidays []*cal.Holiday
	if _, _, holiday := calendar.IsHoliday(day); holiday != nil {
		holidays = append(holidays, holiday)
	}
	return holidays
}

func GetBirthDaysFromDay(day time.Time, codSecao int) []*cal.Holiday {
	return GetBirthDaysFromDays(day, day, codSecao)
}

func GetBirthDaysFromDays(from time.Time, to time.Time, codSecao int) []*cal.Holiday {
	var holidays []*cal.Holiday

	birthdays := repository.GetBirthdays(codSecao, from, to)
	for _, birthday := range birthdays {
		holiday := &cal.Holiday{
			Name:  fmt.Sprintf("%s %s %s", consts.Birthday, birthday.DataNascimento.Format("02/01"), birthday.Nome),
			Month: birthday.DataNascimento.Month(),
			Day:   birthday.DataNascimento.Day(),
		}
		holidays = append(holidays, holiday)
	}
	return holidays
}

func PublishHolidays() {
	schedule, canContinue := ScheduleGet(holidaySchedule, time.Duration(24)*time.Hour, true)
	if !canContinue {
		return
	}
	sections, err := repository.GetAllSections()
	if err != nil {
		log.Printf("Error getting sections %v", err)
		return
	}
	for _, section := range sections {

		holidayStr := ""
		for _, holiday := range GetHolidaysFromDay(time.Now()) {
			holidayStr = fmt.Sprintf("%s %s (%02d/%02d) %s\n", holidayStr, holiday.Description, holiday.Day, holiday.Month, holiday.Name)
		}

		if len(holidayStr) == 0 {
			continue
		}

		chats, err := repository.GetChatsFromCodSecao(section.ID)
		if err != nil {
			continue
		}
		for _, chat := range chats {
			SendTextMessage(chat.ID, holidayStr, 0)
		}

	}
	ScheduleUpdate(&schedule)
}

func PublishBirthdays() {
	birthdayInterval := time.Duration(24*7) * time.Hour
	schedule, canRun := ScheduleGet(birthdaySchedule, birthdayInterval, true)
	if !canRun {
		return
	}
	sections, err := repository.GetAllSections()
	if err != nil {
		log.Printf("Error getting sections %v", err)
		return
	}
	for _, section := range sections {

		birthdayStr := ""
		for _, birthday := range GetBirthDaysFromDays(time.Now(), time.Now().Add(birthdayInterval), section.ID) {
			birthdayStr = birthdayStr + birthday.Name + "\n"
		}

		chats, err := repository.GetChatsFromCodSecao(section.ID)
		if err != nil {
			continue
		}
		for _, chat := range chats {
			if len(birthdayStr) > 0 {
				SendTextMessage(chat.ID, fmt.Sprintf("<b>Aniversários</b>\n%s", birthdayStr), 0)
			}
		}
	}
	ScheduleUpdate(&schedule)
}
