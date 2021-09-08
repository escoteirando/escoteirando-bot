package repository

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/dto"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"gorm.io/gorm/clause"
	"log"
)

func SaveEscotista(escotista domain.MappaEscotista) error {
	CanIWrite()
	GetDB().Save(&escotista)
	YouCanWrite()
	return nil
}

func SaveAssociado(associado domain.MappaAssociado) error {
	CanIWrite()
	GetDB().Save(&associado)
	YouCanWrite()
	return nil
}

func SaveGrupo(grupo domain.MappaGrupo) error {
	CanIWrite()
	GetDB().Save(&grupo)
	YouCanWrite()
	return nil
}

func SaveSecao(secao domain.MappaSecao) error {
	CanIWrite()
	GetDB().Save(&secao)
	YouCanWrite()
	return nil
}
func SaveSubSecao(subSecao domain.MappaSubSecao) error {
	CanIWrite()
	GetDB().Save(&subSecao)
	YouCanWrite()
	return nil
}

func FindAssociado(codigoAssociado string) (domain.MappaAssociado, error) {
	var associado domain.MappaAssociado
	username := codigoAssociado[0 : len(codigoAssociado)-1]
	numeroDigito := codigoAssociado[len(codigoAssociado)-1:]
	result := GetDB().Where("username=? and numero_digito = ?", username, numeroDigito).First(&associado)
	if result.Error == nil && result.RowsAffected > 0 {
		return associado, nil
	}
	return associado, fmt.Errorf("associado não encontrado: %s", codigoAssociado)
}

func VincularAssociado(associado domain.MappaAssociado, user *tgbotapi.User) error {
	usuario := domain.User{}
	result := GetDB().Find(&usuario, user.ID)
	if result.Error != nil || result.RowsAffected < 1 {
		usuario = domain.User{
			ID:           user.ID,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			UserName:     user.UserName,
			LanguageCode: user.LanguageCode,
			Associados: []*domain.MappaAssociado{
				&associado,
			},
		}
		CanIWrite()
		if result = GetDB().Create(&usuario); result.Error == nil {
			result = GetDB().Save(&usuario)
		}
		YouCanWrite()

		return result.Error

	}

	usuario.Associados = append(usuario.Associados, &associado)
	CanIWrite()
	result = GetDB().Save(&usuario)
	YouCanWrite()
	if result.Error != nil {
		return result.Error
	}
	return nil
}



func SaveProgressoes(progressoes []mappa.Progressao) error {
	if len(progressoes) == 0 {
		return fmt.Errorf("empty progressoes list")
	}
	dbProgressoes := make([]domain.MappaProgressao, len(progressoes))
	for i, mappaProgressao := range progressoes {
		dbProgressoes[i] = dto.CreateProgressao(mappaProgressao)
	}
	CanIWrite()
	result := GetDB().
		Clauses(clause.OnConflict{UpdateAll: true}).
		CreateInBatches(dbProgressoes, 10)
	YouCanWrite()
	if result.Error != nil {
		log.Printf("Error saving progressoes: %v", result.Error)
		return result.Error
	}
	log.Printf("Progressoes saved: %d", len(progressoes))
	return nil
}
