package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	bot2 "github.com/guionardo/escoteirando-bot/src/bot"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"log"
)

func Setup() {
	repository.SetupDatabase()
	setupCommands()
	//GetCurrentBot()
	setupBot()
}

func setupBot() {
	setup := utils.GetEnvironmentSetup()
	config, err := repository.GetConfig()
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}
	bot2.CreateBot(setup.BotToken, setup.BotDebug, config.AdminChat)
}

func setupCommands() {
	RegisterCommand(Command{
		Name:         "help",
		Method:       CommandShowHelp,
		JustForAdmin: false,
		Description:  "Ajuda",
	})
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
	RegisterCommand(Command{
		Name:         "imadmin",
		Method:       CommandMainAdmin,
		JustForAdmin: false,
		Description:  "Definir o administrador mestre",
	})
	RegisterCommand(Command{
		Name:         "nomainadmin",
		Method:       CommandNotMainAdmin,
		JustForAdmin: false,
		Description:  "Desabilitar o chat administrador mestre",
	})
	RegisterCommand(Command{
		Name:         "aniversarios",
		Method:       CommandBirthdays,
		JustForAdmin: false,
		Description:  "Aniversariantes",
	})
	RegisterCommand(Command{
		Name:         "stats",
		Method:       CommandStats,
		JustForAdmin: true,
		Description:  "Estatísticas",
	})

}

func CommandShowHelp(ctx domain.Context, _ *tgbotapi.Message) {
	//text := consts.Info + " <b>Comandos</b>\n<pre>"
	commands := make([]CommandButton, 0)

	for _, command := range messageCommands {
		if command.JustForAdmin && !ctx.UserIsAdmin {
			continue
		}
		commands = append(commands, CommandButton{
			Message: command.Description,
			Command: command.Name,
		})

		//text = fmt.Sprintf("%s\n|%12s|%s|", text, command.Name, command.Description)
	}
	//text += "</pre>"
	//SendTextMessage(ctx.ChatId, text, 0)
	SendCommandButtons(ctx.ChatId, "Comandos", commands)
}
