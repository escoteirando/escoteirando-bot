package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

type (
	EnvironmentSetup struct {
		BotToken      string
		BotDebug      bool
		FrontEndUrl   string
		MappaProxyUrl string
		DatabaseFile  string
	}
)

var currentEnvironment = EnvironmentSetup{}

func GetEnvironmentSetup() EnvironmentSetup {
	if len(currentEnvironment.MappaProxyUrl) == 0 {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		currentEnvironment = EnvironmentSetup{
			BotToken:      os.Getenv("BOT_TOKEN"),
			BotDebug:      os.Getenv("BOT_DEBUG") == "1",
			FrontEndUrl:   strings.TrimSuffix(os.Getenv("FRONTEND_URL"), "/"),
			MappaProxyUrl: strings.TrimSuffix(os.Getenv("MAPPA_PROXY_URL"), "/"),
			DatabaseFile:  os.Getenv("BOT_DATABASE"),
		}
		log.Printf("Loaded setup: %v", currentEnvironment)
	}
	return currentEnvironment
}
