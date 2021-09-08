package dto

import (
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/mappa"
	"gorm.io/gorm"
	"hash/crc32"
)

func CreateConquista(conquista mappa.Conquista) domain.MappaConquista {
	return domain.MappaConquista{
		Model:                          gorm.Model{ID: GetIdConquista(conquista)},
		Type:                           conquista.Type,
		DataConquista:                  conquista.DataConquista,
		CodigoAssociado:                conquista.CodigoAssociado,
		CodigoEscotistaUltimaAlteracao: conquista.CodigoEscotistaUltimaAlteracao,
		CodigoSecao:                    0,
		CodigoEspecialidade:            conquista.CodigoEspecialidade,
		NumeroNivel:                    conquista.NumeroNivel,
		NotificadoChat:                 false,
		Associado:                      domain.MappaAssociado{},
	}
}

func GetIdConquista(conquista mappa.Conquista) uint {
	hashType := uint(crc32.ChecksumIEEE([]byte(conquista.Type)))
	return uint(conquista.CodigoAssociado)*1000000000 + hashType
}
