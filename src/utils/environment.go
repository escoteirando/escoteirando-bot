package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

type (
	EnvironmentSetup struct {
		BotToken         string
		BotDebug         bool
		FrontEndUrl      string
		MappaProxyUrl    string
		DatabaseFile     string
		DisabledFeatures []string
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
			BotToken:         os.Getenv("BOT_TOKEN"),
			BotDebug:         os.Getenv("BOT_DEBUG") == "1",
			FrontEndUrl:      strings.TrimSuffix(os.Getenv("FRONTEND_URL"), "/"),
			MappaProxyUrl:    strings.TrimSuffix(os.Getenv("MAPPA_PROXY_URL"), "/"),
			DatabaseFile:     os.Getenv("BOT_DATABASE"),
			DisabledFeatures: strings.Split(os.Getenv("DISABLED_FEATURES"), ","),
		}
		log.Printf("Loaded setup: %v", currentEnvironment)
	}
	return currentEnvironment
}

func IsDisabledFeature(feature string) bool {
	for _, f := range currentEnvironment.DisabledFeatures {
		if strings.ToUpper(strings.Trim(f, " ")) == strings.ToUpper(feature) {
			log.Printf("Disabled feature %s", feature)
			return true
		}
	}
	return false
}
