package service

import "github.com/adeteam/git-mirror/definition"

var (
	_config *ConfigService
)

type ConfigService struct {
	Current definition.Config
}

func Config() *ConfigService {
	if _config == nil {
		_config = NewConfigService()
	}

	return _config
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}
