package config

import (
	"github.com/dst3v3n/api-anime/pkg/logger"
	"github.com/rs/zerolog"
)

func Logging() zerolog.Logger {
	c := getConfig()

	return logger.InitLogger(c.LogConfig.LogEnv, c.LogConfig.LogAppName, "2.0.0")
}
