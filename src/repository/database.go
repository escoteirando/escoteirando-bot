package repository

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var (
	currentDb *gorm.DB = nil
)

func SetupDatabase(databaseFile string) error {
	env := utils.GetEnvironmentSetup()
	config := gorm.Config{}
	if env.BotDebug {
		config.Logger = logger.Default.LogMode(logger.Info)
	}
	db, err := gorm.Open(sqlite.Open(databaseFile), &config)
	if err != nil {
		return err
	}

	// Migrate the schema
	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&domain.Chat{})
	db.AutoMigrate(&domain.Config{})
	db.AutoMigrate(&domain.MappaEscotista{})
	db.AutoMigrate(&domain.MappaAssociado{})
	db.AutoMigrate(&domain.MappaGrupo{})
	db.AutoMigrate(&domain.MappaSecao{})
	db.AutoMigrate(&domain.MappaSubSecao{})
	db.AutoMigrate(&domain.MappaMarcacao{})
	db.AutoMigrate(&domain.MappaProgressao{})
	db.AutoMigrate(&domain.Run{})
	db.AutoMigrate(&domain.Schedule{})

	log.Printf("Database is open")
	currentDb = db
	return nil
}

func GetDB() *gorm.DB {
	if currentDb == nil {
		log.Panic("Uninitialized database")
	}
	return currentDb
}
