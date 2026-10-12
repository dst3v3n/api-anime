// Package logger configura e inicializa el sistema de logging de la aplicación usando Zerolog.
// Proporciona una función para crear un logger estructurado con diferentes niveles de log
// según el entorno (development, staging, production), con salida formateada adecuadamente
// para cada ambiente.
package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger inicializa y configura el logger de la aplicación con Zerolog.
// Configura el nivel de logging (Debug en desarrollo, Info en producción),
// el formato de salida (coloreado en consola para desarrollo, JSON en producción),
// y añade campos contextuales como versión, nombre de aplicación y ambiente.
// Retorna una instancia de logger configurada lista para usar en toda la aplicación.
func InitLogger(env string, appName string, version string) zerolog.Logger {
	var output io.Writer

	if env == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		output = os.Stdout
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	}

	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(output).With().
		Timestamp().
		Str("Version", version).
		Str("App", appName).
		Str("Environment", env).
		Logger()

	log.Logger = logger

	return logger
}
