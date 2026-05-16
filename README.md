# 🎌 Anime API

**Librería Go de alto rendimiento** para obtener información de animes desde AnimeFlv con caché distribuido opcional
---

## 📋 Descripción

**Anime API** es una librería Go que te permite:

- 🔍 Buscar animes por nombre con paginación
- 📖 Obtener información detallada (sinopsis, géneros, estado, episodios, relacionados)
- 🎬 Conseguir enlaces de episodios desde servicios externos (Mega, StreamTape, StreamWish, etc.)
- 🎥 Extraer URLs directas de reproducción en tiempo real
- 📺 Monitorear animes y episodios recientes

**Incluye caché distribuido opcional** con Valkey/Redis que reduce tiempos de respuesta de **2-3 segundos a <1ms**.

---

## 🌟 Características Principales

| Característica | Descripción |
|---|---|
| 🔍 **Búsqueda inteligente** | Por nombre con paginación automática |
| 📖 **Info completa** | Sinopsis, géneros, estado, total de episodios |
| 🎬 **Enlaces directos** | URLs de múltiples servicios de streaming |
| 🎥 **Extracción URL** | Automatización de navegador para URLs directas ⚡ |
| 📺 **Feed de actualizaciones** | Animes y episodios más recientes |
| 💾 **Caché configurable** | Valkey/Redis, TTL personalizable |
| 🚀 **Ultra rápido** | <1ms en consultas cacheadas (~3000x más rápido) |

---

## 📦 Instalación

### Requisitos Previos

- **Go 1.25.3+**
- **Valkey/Redis** (opcional, solo si usas caché)

### Instalar la librería

```bash
go get github.com/dst3v3n/api-anime
```

---

## 🚀 Inicio Rápido

### 1️⃣ Sin Caché (Más Simple)

```go
package main

import (
    "context"
    "fmt"
    "github.com/dst3v3n/api-anime"
)

func main() {
    service := apianime.NewAnimeFlv()
    ctx := context.Background()
    
    resultados, err := service.SearchAnime(ctx, "One Piece", 1)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("✅ Encontrados %d animes\n", len(resultados.Animes))
    
    for _, anime := range resultados.Animes {
        fmt.Printf("  📺 %s - ⭐ %.1f\n", anime.Title, anime.Punctuation)
    }
}
```

### 2️⃣ Con Caché (Recomendado para Producción)

**Paso 1:** Inicia Valkey/Redis

```bash
docker run -d -p 6379:6379 valkey/valkey:latest
```

**Paso 2:** Configura tu servicio

```go
package main

import (
    "context"
    "fmt"
    "github.com/dst3v3n/api-anime"
    "github.com/dst3v3n/api-anime/config"
)

func main() {
    // Configurar caché
    cfg := config.NewConfigWithDefaults().
        WithEnableCache(true).              // Activar caché
        WithCacheHost("localhost").         // Host
        WithCachePort(6379).                // Puerto
        WithCacheTTL(60)                    // TTL en minutos
    
    config.InitConfig(cfg)
    
    service := apianime.NewAnimeFlv()
    ctx := context.Background()
    
    // Primera búsqueda: ~2s (scraping del sitio)
    fmt.Println("⏳ Primera búsqueda...")
    resultados, _ := service.SearchAnime(ctx, "Naruto", 1)
    fmt.Printf("✅ Encontrados %d animes\n", len(resultados.Animes))
    
    // Segunda búsqueda: <1ms (desde caché!)
    fmt.Println("⚡ Segunda búsqueda (desde caché)...")
    resultados, _ = service.SearchAnime(ctx, "Naruto", 1)
    fmt.Printf("✅ Encontrados %d animes (casi instantáneo)\n", len(resultados.Animes))
}
```

---

## 📚 Referencia Completa de la API

### 🔎 SearchAnime

Busca animes por nombre con soporte para paginación.

```go
SearchAnime(ctx context.Context, anime string, page uint) (AnimeResponse, error)
```

**Ejemplo:**

```go
resultados, err := service.SearchAnime(ctx, "Naruto", 1)
if err != nil {
    log.Fatal(err)
}

for _, anime := range resultados.Animes {
    fmt.Printf("%s (⭐ %.1f/10) - %s\n", 
        anime.Title, 
        anime.Punctuation,
        anime.Sinopsis[:50] + "...")
}
```

**Respuesta:**

```go
type AnimeResponse struct {
    Animes     []types.AnimeStruct
    TotalPages uint
}

type AnimeStruct struct {
    ID          string        // "naruto-shippuden"
    Title       string        // "Naruto Shippuden"
    Sinopsis    string
    Type        CategoryAnime // Anime, OVA, Pelicula, Especial
    Punctuation float64       // Calificación 0-10
    Image       string        // URL de la portada
}
```

---

### 📖 AnimeInfo

Obtén toda la información detallada de un anime.

```go
AnimeInfo(ctx context.Context, idAnime string) (AnimeInfoResponse, error)
```

**Ejemplo:**

```go
info, err := service.AnimeInfo(ctx, "one-piece-tv")
if err != nil {
    log.Fatal(err)
}

fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
fmt.Printf("📺 %s\n", info.Title)
fmt.Printf("🎯 Tipo: %s\n", info.Type)
fmt.Printf("📊 Estado: %s\n", info.Status)          // "En Emisión" / "Finalizado"
fmt.Printf("🎭 Géneros: %v\n", info.Genres)
fmt.Printf("📍 Total episodios: %d\n", len(info.Episodes))
fmt.Printf("⏰ Próximo episodio: %s\n", info.NextEpisode)

fmt.Println("\n🔗 Animes Relacionados:")
for _, rel := range info.AnimeRelated {
    fmt.Printf("  └─ %s (%s)\n", rel.Title, rel.Category)
}
fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
```

**Respuesta:**

```go
type AnimeInfoResponse struct {
    AnimeStruct                   // Información básica
    AnimeRelated []AnimeRelated   // Precuelas, secuelas, spin-offs
    Genres       []string
    Status       StatusAnime      // "En Emisión" / "Finalizado"
    NextEpisode  string
    Episodes     []int            // [1, 2, 3, ..., 1150]
}

type AnimeRelated struct {
    ID       string
    Title    string
    Category string // "Precuela", "Secuela", "Spin-off", etc.
    Image    string
}
```

---

### 🔗 Links

Obtén todos los enlaces de descarga/streaming disponibles para un episodio.

```go
Links(ctx context.Context, idAnime string, episode uint) (LinkResponse, error)
```

**Ejemplo:**

```go
links, err := service.Links(ctx, "one-piece-tv", 1150)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("📺 %s - Episodio %d\n", links.Title, links.Episode)
fmt.Printf("🌐 Servidores disponibles: %d\n\n", len(links.Link))

for i, link := range links.Link {
    fmt.Printf("%d. %s\n", i+1, link.Server)
    fmt.Printf("   🔗 %s\n", link.URL)
    if link.Code != "" {
        fmt.Printf("   📝 Código: %s\n", link.Code)
    }
    fmt.Println()
}
```

**Respuesta:**

```go
type LinkResponse struct {
    ID      string
    Title   string
    Episode uint
    Link    []types.LinkSource
}

type LinkSource struct {
    Server string // "Mega", "StreamTape", "StreamWish", etc.
    URL    string // URL directa/embebida
    Code   string // Código de inserción (si aplica)
}
```

---

### ⚡ ExtractURL

Extrae URLs directas de reproducción desde páginas embebidas de video filtrando por resolución en tiempo real.

```go
ExtractURL(service string, ctx context.Context, url string, resolution string) (string, error)
```

**Ejemplo:**

```go
import "github.com/dst3v3n/api-anime/extract"

// URL embebida de un reproductor
embedURL := "https://streamwish.to/e/ss619zjv2ufo"

// Extraer URL directa especificando la resolución vertical deseada ("480", "720", etc.)
// También puedes enviar "default" o "" para obtener la primera opción disponible.
videoURL, err := extract.ExtractURL("streamwish", ctx, embedURL, "480")
if err != nil {
    log.Fatal(err)
}

fmt.Println("✅ URL directa del video (480p):")
fmt.Println(videoURL)
// Output: https://hgplaycdn.com/stream/.../index-f1-v1-a1.m3u8
```

**Servicios Soportados:**

| Servicio | Estado | Características |
|----------|--------|-----------------|
| **StreamTape** | ✅ Disponible | Extracción nativa directa |
| **StreamWish** | ✅ Disponible | Captura optimizada de master.m3u8 |
| Mega | ⏳ Próxima versión | En desarrollo |
| Google Drive | ⏳ Próxima versión | En desarrollo |

**Requisitos del Sistema:**

Chrome/Chromium debe estar instalado:

```bash
# Linux (Ubuntu/Debian)
sudo apt-get install chromium-browser

# macOS
brew install chromium

# Windows
# Descargar desde: https://www.chromium.org/
```

**Ejemplo Completo - Obtener y Extraer URLs:**

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/dst3v3n/api-anime"
    "github.com/dst3v3n/api-anime/extract"
)

func main() {
    service := apianime.NewAnimeFlv()
    ctx := context.Background()
    
    // 1. Obtener información del anime
    info, _ := service.AnimeInfo(ctx, "one-piece-tv")
    fmt.Printf("📺 %s\n", info.Title)
    fmt.Printf("📊 Total episodios: %d\n\n", len(info.Episodes))
    
    // 2. Obtener enlaces del último episodio
    lastEp := info.Episodes[len(info.Episodes)-1]
    links, _ := service.Links(ctx, "one-piece-tv", uint(lastEp))
    fmt.Printf("🔗 Enlaces disponibles para episodio %d:\n\n", lastEp)
    
    // 3. Buscar StreamWish y extraer URL directa
    for _, link := range links.Link {
        if link.Server == "StreamWish" {
            fmt.Printf("🎯 Servidor: %s\n", link.Server)
            fmt.Printf("📍 URL embebida: %s\n", link.URL)
            
            // Extraer URL directa
            fmt.Println("⏳ Extrayendo URL directa...")
            videoURL, err := extract.ExtractURL("streamwish", ctx, link.URL)
            if err != nil {
                fmt.Printf("❌ Error: %v\n", err)
                continue
            }
            
            fmt.Printf("✅ URL directa: %s\n", videoURL)
            break
        }
    }
}
```

---

### 📺 RecentAnime

Obtén los últimos animes agregados al sitio.

```go
RecentAnime(ctx context.Context) ([]AnimeStruct, error)
```

**Ejemplo:**

```go
recientes, _ := service.RecentAnime(ctx)

fmt.Println("🆕 Últimos animes agregados:\n")
for i, anime := range recientes[:5] {
    fmt.Printf("%d. %s (%s)\n", i+1, anime.Title, anime.Type)
}
```

---

### 🆕 RecentEpisode

Obtén los últimos episodios publicados.

```go
RecentEpisode(ctx context.Context) ([]EpisodeListResponse, error)
```

**Ejemplo:**

```go
episodios, _ := service.RecentEpisode(ctx)

fmt.Println("📢 Últimos episodios publicados:\n")
for i, ep := range episodios[:10] {
    fmt.Printf("%d. %s - Ep. %d\n", i+1, ep.Title, ep.Episode)
}
```

---

## 🔧 Configuración

### Opción 1: Configuración Programática (Recomendado)

```go
import "github.com/dst3v3n/api-anime/config"

cfg := config.NewConfigWithDefaults().
    WithEnableCache(true).                      // Activar
    WithCacheHost("redis.prod.com").            // Host
    WithCachePort(6380).                        // Puerto
    WithCachePassword("your-password").         // Contraseña (opcional)
    WithCacheDB(0).                             // Base de datos (0-15)
    WithCacheTTL(120)                           // TTL en minutos

config.InitConfig(cfg)
service := apianime.NewAnimeFlv()
```

### Opción 2: Variables de Entorno (.env)

```bash
CACHE_ENABLED=true
CACHE_HOST=localhost
CACHE_PORT=6379
CACHE_DB=0
CACHE_PASSWORD=
CACHE_TTL=60
```

```go
service := apianime.NewAnimeFlv()  // Carga automáticamente .env
```

### Opción 3: Archivo .env Personalizado

```go
cfg, err := config.NewConfigFromEnvPath("/ruta/custom/.env")
if err != nil {
    panic(err)
}
config.InitConfig(cfg)
```

### Opciones de Configuración

| Método | Tipo | Default | Rango | Descripción |
|--------|------|---------|-------|-------------|
| `WithEnableCache(bool)` | bool | `false` | `true/false` | Activar/desactivar caché |
| `WithCacheHost(string)` | string | `localhost` | - | Host de Valkey/Redis |
| `WithCachePort(int)` | int | `6379` | 1-65535 | Puerto del servidor |
| `WithCachePassword(string)` | string | `""` | - | Contraseña (si aplica) |
| `WithCacheDB(int)` | int | `0` | 0-15 | Base de datos Redis |
| `WithCacheTTL(int)` | int | `60` | 1-1440 | TTL en minutos |

### Ejemplos de Configuración por Entorno

**🔨 Desarrollo sin caché:**

```go
cfg := config.NewConfigWithDefaults()
config.InitConfig(cfg)
```

**🏠 Desarrollo con caché local:**

```go
cfg := config.NewConfigWithDefaults().
    WithEnableCache(true)
config.InitConfig(cfg)
```

**🚀 Producción con Redis remoto:**

```go
cfg := config.NewConfigWithDefaults().
    WithEnableCache(true).
    WithCacheHost(os.Getenv("REDIS_HOST")).
    WithCachePort(6380).
    WithCachePassword(os.Getenv("REDIS_PASSWORD")).
    WithCacheTTL(30)  // 30 minutos más corto en producción
config.InitConfig(cfg)
```

**🌍 Multi-entorno dinámico:**

```go
func newService(env string) *apianime.AnimeFlv {
    var cfg *config.Config
    
    switch env {
    case "production":
        cfg = config.NewConfigWithDefaults().
            WithEnableCache(true).
            WithCacheHost("redis.prod.com").
            WithCacheTTL(60)
    case "staging":
        cfg = config.NewConfigWithDefaults().
            WithEnableCache(true).
            WithCacheHost("redis.staging.com").
            WithCacheTTL(30)
    default:  // development
        cfg = config.NewConfigWithDefaults().
            WithEnableCache(false)
    }
    
    config.InitConfig(cfg)
    return apianime.NewAnimeFlv()
}
```

---

## 💾 Sistema de Caché

### ¿Qué Se Cachea?

| Función | Clave de Caché | TTL Default |
|---------|--------|-------------|
| `SearchAnime` | `search-anime-{nombre}-page-{N}` | Configurable (default 60min) |
| `AnimeInfo` | `anime-info-{id}` | Configurable (default 60min) |
| `Links` | `links-{id}-{episodio}` | Configurable (default 60min) |
| `RecentAnime` | `recent-anime` | Configurable (default 60min) |
| `RecentEpisode` | `recent-episode` | Configurable (default 60min) |

### Comparativa de Performance

| Operación | Sin Caché | Con Caché | **Mejora** |
|-----------|----------|----------|-----------|
| SearchAnime | ~2.5s | ~0.8ms | **🔥 3100x** |
| AnimeInfo | ~1.8s | ~0.6ms | **🔥 3000x** |
| Links | ~1.5s | ~0.5ms | **🔥 3000x** |

### Invalidar Caché Manualmente

```go
// Desactivar temporalmente
cfg := config.NewConfigWithDefaults().WithEnableCache(false)
config.InitConfig(cfg)
resultados, _ := service.SearchAnime(ctx, "Naruto", 1)

// Reactivar
cfg = config.NewConfigWithDefaults().WithEnableCache(true)
config.InitConfig(cfg)
```

---

## 💡 Casos de Uso Comunes

### Buscar y Explorar Animes

```go
// 1. Buscar
resultados, _ := service.SearchAnime(ctx, "Attack on Titan", 1)

// 2. Obtener detalles del primero
info, _ := service.AnimeInfo(ctx, resultados.Animes[0].ID)

// 3. Ver animes relacionados
fmt.Printf("Relacionados a %s:\n", info.Title)
for _, rel := range info.AnimeRelated {
    fmt.Printf("  └─ %s (%s)\n", rel.Title, rel.Category)
}
```

### Obtener Enlaces de Todos los Episodios

```go
info, _ := service.AnimeInfo(ctx, "shingeki-no-kyojin")

for _, ep := range info.Episodes {
    links, _ := service.Links(ctx, "shingeki-no-kyojin", uint(ep))
    fmt.Printf("Ep.%d: %d servidores disponibles\n", ep, len(links.Link))
    
    for _, link := range links.Link {
        fmt.Printf("  • %s\n", link.Server)
    }
}
```

### Monitorear Nuevos Episodios

```go
for {
    episodios, _ := service.RecentEpisode(ctx)
    
    for _, ep := range episodios {
        fmt.Printf("[NEW] %s - Cap. %s\n", ep.Title, ep.Chapter)
    }
    
    time.Sleep(1 * time.Hour)
}
```

### Descargar Todos los Episodios de un Anime

```go
info, _ := service.AnimeInfo(ctx, "one-piece-tv")
fmt.Printf("Descargando %s (%d episodios)...\n", info.Title, len(info.Episodes))

for _, ep := range info.Episodes {
    links, _ := service.Links(ctx, "one-piece-tv", uint(ep))
    
    for _, link := range links.Link {
        if link.Server == "Mega" {  // Preferir Mega
            fmt.Printf("[%d/%d] Episodio %d: %s\n", 
                ep, len(info.Episodes), ep, link.URL)
            break
        }
    }
}
```

---

## 🧪 Testing

```bash
# Ejecutar todos los tests
go test ./...

# Con cobertura de código
go test ./... -cover

# Tests específicos de extractores
go test ./internal/adapters/scrapers/animeflv -v

# Verbose output
go test ./... -v -run TestSearchAnime
```

---

## ❓ FAQ

**P: ¿Es obligatorio usar Valkey/Redis?**  
R: No. El caché es opcional. Funciona perfectamente sin él.

**P: ¿Cómo activo el caché?**  
R:

```go
cfg := config.NewConfigWithDefaults().WithEnableCache(true)
config.InitConfig(cfg)
```

**P: ¿Puedo cambiar el TTL?**  
R: Sí:

```go
cfg.WithCacheTTL(120)  // 2 horas
```

**P: ¿Funciona con Redis en lugar de Valkey?**  
R: Sí, 100% compatible. Usa los mismos métodos.

**P: ¿Los enlaces caducan?**  
R: Algunos servidores tienen URLs temporales (1-24 horas). Por eso el TTL predeterminado es de 60 minutos.

**P: ¿Puedo usar en producción?**  
R: Sí, pero el scraping depende de la estructura del sitio. Monitorea cambios regularmente.

**P: ¿Necesito Chrome instalado para ExtractURL?**  
R: Sí, es obligatorio. La librería usa automatización de navegador nativo.

---

## ⚖️ Aviso Legal

**Uso único para fines educativos y personales.** Respeta siempre los términos de servicio de AnimeFlv.

### ✅ Obligaciones

- Respeta `robots.txt` del sitio
- Usa solo para proyectos personales/educativos
- Cita la fuente (AnimeFlv)
- Implementa rate limiting en producción

### ❌ Prohibido

- Comercialización sin permiso explícito
- Ataques DDoS o sobrecarga del servidor
- Distribución sin atribución
- Scraping masivo sin respetar tiempos

---

## 📄 Licencia

**MIT** - Libre para usar, modificar y distribuir.  
Ver archivo [LICENSE](LICENSE) para detalles completos.

---

## 👤 Autor

**Steven** ([@dst3v3n](https://github.com/dst3v3n))

---

## 🤝 Contribuir

¡Las contribuciones son bienvenidas!

- 🐛 **Reportar bugs:** [Issues](../../issues)
- 💡 **Sugerir features:** [Discussions](../../discussions)
- 🔧 **Enviar código:** [Pull Requests](../../pulls)

---

## 📞 Soporte

- **GitHub:** [@dst3v3n](https://github.com/dst3v3n)
- **Issues:** [GitHub Issues](../../issues)
- **Discussions:** [GitHub Discussions](../../discussions)

---

**Made with ❤️ by Steven**
