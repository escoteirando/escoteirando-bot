package repository

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"gorm.io/gorm"
	"time"
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

func GetUltimaConsultaMarcacao(codSecao int) time.Time {
	var secao domain.MappaSecao
	response := GetDB().Find(&secao, codSecao)
	if response.Error != nil {
		return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return secao.UltimaConsultaMarcacao
}

func SetUltimaConsultaMarcacao(codSecao int, quando time.Time) error {
	var secao domain.MappaSecao
	response := GetDB().Find(&secao, codSecao)
	if response.Error != nil {
		return response.Error
	}
	CanIWrite()
	secao.UltimaConsultaMarcacao = quando
	response = GetDB().Save(&secao)
	YouCanWrite()
	return response.Error
}

func SetMarcacoes(marcacoes []mappa.Marcacao) error {
	dbMarcacoes := make([]domain.MappaMarcacao, len(marcacoes))
	uniqueSequence := utils.NewUniqueSequence()
	for i, marcacao := range marcacoes {
		dbMarcacoes[i] = domain.MappaMarcacao{
			Model:                 gorm.Model{ID: uniqueSequence.GetNext()},
			CodigoAtividade:       marcacao.CodigoAtividade,
			CodigoAssociado:       marcacao.CodigoAssociado,
			DataAtividade:         utils.TimeParse(marcacao.DataAtividade),
			DataStatusJovem:       utils.TimeParse(marcacao.DataStatusJovem),
			DataStatusEscotista:   utils.TimeParse(marcacao.DataStatusEscotista),
			StatusJovem:           marcacao.StatusJovem,
			StatusEscotista:       marcacao.StatusEscotista,
			DataHoraAtualizacao:   utils.TimeParse(marcacao.DataHoraAtualizacao),
			CodigoUltimoEscotista: marcacao.CodigoUltimoEscotista,
			Segmento:              marcacao.Segmento,
			NotificadoChat:        false,
			Associado:             domain.MappaAssociado{},
			Progressao:            domain.MappaProgressao{},
		}
	}
	CanIWrite()
	result := GetDB().Save(&dbMarcacoes)
	YouCanWrite()
	return result.Error
}
