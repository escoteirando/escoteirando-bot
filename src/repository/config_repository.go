package repository

import "github.com/guionardo/escoteirando-bot/src/domain"

func GetConfig() (domain.Config, error) {
	var config domain.Config
	result := GetDB().First(&config)
	err := ParseResponse(result, "GetConfig", 0)
	return config, err
}

func SetConfig(config *domain.Config) error {
	CanIWrite()
	result := GetDB().Save(config)
	YouCanWrite()
	return ParseResponse(result, "SetConfig", 0)
}
