package config

import (
	"github.com/dst3v3n/api-anime/pkg/logger"
	"github.com/rs/zerolog"
)

// Logging inicializa y retorna un logger de Zerolog configurado según los parámetros de Config.
// El logger se configura con el entorno y nombre de aplicación especificados.
func Logging() zerolog.Logger {
	c := getConfig()

	return logger.InitLogger(c.LogConfig.LogAppName, c.LogConfig.LogEnv)
}
