package config

import (
	"github.com/spf13/viper"
)

// ConfigService exposes configuration to the Wails frontend.
type ConfigService struct{}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (s *ConfigService) Get(key string) any {
	return viper.Get(key)
}

func (s *ConfigService) GetString(key string) string {
	return viper.GetString(key)
}

func (s *ConfigService) GetInt(key string) int {
	return viper.GetInt(key)
}

func (s *ConfigService) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (s *ConfigService) Set(key string, value any) error {
	return Set(key, value)
}

func (s *ConfigService) GetAll() map[string]any {
	return viper.AllSettings()
}

func (s *ConfigService) GetTheme() string        { return Theme() }
func (s *ConfigService) SetTheme(t string) error  { return Set("theme", t) }
func (s *ConfigService) GetADOOrg() string        { return ADOOrganization() }
func (s *ConfigService) GetADOProject() string    { return ADOProject() }
