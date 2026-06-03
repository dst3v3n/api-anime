package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	domainErrors "github.com/dst3v3n/api-anime/internal/domain/errors"
	"github.com/spf13/viper"
)

var (
	instance Config
	once     sync.Once
)

func InitConfig() error {
	var initErr error

	once.Do(func() {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			initErr = domainErrors.NewNotFoundErr("No Existe Directorio de Configuración", "config")
			return
		}

		setConfigDefault()
		viper.AddConfigPath(userConfigDir)
		viper.SetConfigType("json")

		cliConfigPath := filepath.Join(userConfigDir, ".configAnimeCli.json")

		if _, err := os.Stat(cliConfigPath); err == nil {
			viper.SetConfigName(".configAnimeCli")
		} else {
			viper.SetConfigName(".configApiAnime")
		}

		if err := viper.ReadInConfig(); err != nil {
			var configFileNotFoundErr viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundErr) {
				if err := viper.SafeWriteConfig(); err != nil {
					initErr = domainErrors.NewNotSetConfig("No se pudo guardar la configuración inicial", "config-viper")
					return
				}
			} else {
				initErr = domainErrors.NewNotSetConfig("Error al leer el archivo de configuración existente", "config-viper")
				return
			}
		}

		if err := viper.Unmarshal(&instance); err != nil {
			initErr = domainErrors.NewNotSetConfig("Error al mapear la configuración a la instancia", "config-viper")
			return
		}
	})

	return initErr
}

func getConfig() Config {
	var cfg Config
	_ = viper.Unmarshal(&cfg)
	return cfg
}

func GetConfig() CacheConfig {
	var cfg Config
	_ = viper.Unmarshal(&cfg)
	return cfg.CacheConfig
}

func ResetConfig() {
	viper.Reset()
}

func SetConfigCache(data CacheConfig) error {
	configResult := getConfig()
	configResult.CacheConfig = data

	bytes, err := json.Marshal(configResult)
	if err != nil {
		return err
	}

	var miMapa map[string]any
	if err := json.Unmarshal(bytes, &miMapa); err != nil {
		return err
	}

	if err := viper.MergeConfigMap(miMapa); err != nil {
		return err
	}

	return viper.WriteConfig()
}
