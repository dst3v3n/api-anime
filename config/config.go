package config

import (
	"github.com/dst3v3n/api-anime/internal/config"
	"github.com/dst3v3n/api-anime/types"
)

func InitConfig() error {
	return config.InitConfig()
}

func GetConfig() config.CacheConfig {
	return config.GetConfig()
}

func ResetConfig() {
	config.ResetConfig()
}

func SetConfig(data types.ConfigCache) error {
	return config.SetConfigCache(data)
}
