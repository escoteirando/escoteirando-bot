package services

import (
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/utils"
)

func Setup() {
	repository.SetupDatabase(utils.GetEnvironmentSetup().DatabaseFile)
	setupCommands()
	SetupCalendar()
	GetBot()
}

func setupCommands() {
	RegisterCommand(Command{
		Name:         "setup",
		Method:       CommandSetup,
		JustForAdmin: false,
		Description:  "Configuração do bot no grupo",
	})
	RegisterCommand(Command{
		Name:         "auth",
		Method:       CommandAuth,
		JustForAdmin: true,
		Description:  "Autorização do mAPPa",
	})
}
