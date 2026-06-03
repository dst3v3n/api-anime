// Package animeflv contiene tests unitarios para la capa de caché del sistema.
// Este archivo (cache_test.go) implementa pruebas para verificar el correcto funcionamiento
// de operaciones de caché: almacenamiento (Set), recuperación (Get), verificación de
// existencia (Exists) y eliminación (Delete) de valores en Valkey.
package animeflv

import (
	"context"
	"fmt"
	"testing"

	"github.com/dst3v3n/api-anime/internal/adapters/cache"
	"github.com/dst3v3n/api-anime/internal/config"
	"github.com/dst3v3n/api-anime/internal/domain/dto"
	"github.com/dst3v3n/api-anime/internal/mocks"
	"github.com/valkey-io/valkey-go"
)

func TestCacheSet(t *testing.T) {
	testCases := []struct {
		name        string
		mock        interface{}
		ID          string
		key         string
		wantError   bool
		description string
	}{
		{
			name:        "establecer en caché exitoso recent-anime",
			mock:        mocks.MockAnimeStruct(),
			key:         "recent-anime",
			wantError:   false,
			description: "debe establecer correctamente el valor en caché sin errores",
		},
		{
			name:        "establecer en caché con error recent-anime",
			mock:        nil,
			key:         "recent-anime",
			wantError:   true,
			description: "debe manejar correctamente el error al establecer un valor no serializable en caché",
		},
		{
			name:        "establecer en caché exitoso recent-episode",
			mock:        mocks.MockEpisodeListResponse(),
			key:         "recent-episode",
			wantError:   false,
			description: "debe establecer correctamente el valor en caché sin errores",
		},
		{
			name:        "establecer en caché error recent-episode",
			mock:        "hola mundo",
			ID:          "",
			key:         "recent-episode",
			wantError:   true,
			description: "debe manejar correctamente el error al establecer un valor no serializable en caché",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := config.InitConfig(); err != nil {
				t.Errorf("Error en la inicialización de la configuración %v", err)
			}
			cfg := config.GetConfig()
			t.Log(cfg)
			initAddress := fmt.Sprintf("redis://%s:%d/%d", cfg.CacheHost, cfg.CachePort, cfg.CacheDB)

			client, err := valkey.NewClient(valkey.MustParseURL(initAddress))
			if err != nil {
				t.Fatalf("error initializing Valkey client: %v", err)
			}
			cache := cache.NewValkeyCache(client)

			ctx := context.Background()

			cacheErr := cache.Set(ctx, tc.key, tc.mock)
			if (cacheErr != nil) != tc.wantError {
				t.Errorf("Eror configuración de caché = %v, wantError %v", cacheErr, tc.wantError)
			}

			exist, err := cache.Exists(ctx, tc.key)
			if err != nil {
				t.Fatalf("error checking existence in cache: %v", err)
			}

			if !tc.wantError {
				if !exist {
					t.Errorf("La clave %s debería existir en caché pero no se encontró", tc.key)
				}
			}
		})
	}
}

func TestCacheGet(t *testing.T) {
	testCases := []struct {
		name         string
		key          string
		description  string
		wantError    bool
		variableType interface{}
	}{
		{
			name:         "obtener recent-anime exitoso",
			key:          "recent-anime",
			description:  "debe obtener correctamente el valor en caché sin errores",
			wantError:    false,
			variableType: dto.AnimeStruct{},
		},
		{
			name:         "obtener recent-anime fallido",
			key:          "recent-animes",
			description:  "debe manejar correctamente el error al obtener un valor no existente en caché",
			wantError:    true,
			variableType: dto.AnimeStruct{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := config.InitConfig(); err != nil {
				t.Errorf("Error en la inicialización de la configuración %v", err)
			}
			cfg := config.GetConfig()
			t.Log(cfg)

			initAddress := fmt.Sprintf("redis://%s:%d/%d", cfg.CacheHost, cfg.CachePort, cfg.CacheDB)
			t.Logf("initAdress %s", initAddress)
			client, err := valkey.NewClient(valkey.MustParseURL(initAddress))
			if err != nil {
				t.Fatalf("error initializing Valkey client: %v", err)
			}
			data := tc.variableType

			cache := cache.NewValkeyCache(client)

			ctx := context.Background()

			error := cache.Get(ctx, tc.key, &data)
			if (error != nil) != tc.wantError {
				t.Errorf("Error obteniendo de caché = %v", err)
			}

			if !tc.wantError {
				if data == nil {
					t.Errorf("La variable debería haberse llenado con datos desde la caché")
				}
				t.Logf("Datos obtenidos de caché: %+v", data)
			}
		})
	}
}
