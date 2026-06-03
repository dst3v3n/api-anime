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
    service, err := apianime.NewAnimeFlv()
    if err != nil {
        panic(err)
    }
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
    
    service, err := apianime.NewAnimeFlv()
    if err != nil {
        panic(err)
    }
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

### 🌐 Search

Obtén todos los animes disponibles sin filtros de búsqueda con soporte para paginación.

```go
Search(ctx context.Context, page uint) (AnimeResponse, error)
```

**Ejemplo:**

```go
resultados, err := service.Search(ctx, 1)
if err != nil {
    log.Fatal(err)
}

for _, anime := range resultados.Animes {
    fmt.Printf("%s\n", anime.Title)
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

Extrae URLs directas de reproducción desde páginas embebidas de video. Retorna una lista de URLs con sus respectivas resoluciones detectadas.

```go
ExtractURL(service string, ctx context.Context, url string) ([]types.VideoURL, error)
```

**Ejemplo:**

```go
import "github.com/dst3v3n/api-anime/extract"

// URL embebida de un reproductor
embedURL := "https://streamwish.to/e/ss619zjv2ufo"

// Extraer URLs directas
videos, err := extract.ExtractURL("streamwish", ctx, embedURL)
if err != nil {
    log.Fatal(err)
}

for _, video := range videos {
    fmt.Printf("✅ Resolución: %s - URL: %s\n", video.Resolution, video.URL)
}
```

**Respuesta:**

```go
type VideoURL struct {
 URL        string `json:"url"`
 Resolution string `json:"resolution"`
}
```

**Servicios Soportados:**

| Servicio | Estado | Características |
|----------|--------|-----------------|
| **streamstape** | ✅ Disponible | Extracción nativa directa |
| **streamwish** | ✅ Disponible | Captura optimizada de master.m3u8 |
| mega | ⏳ Próxima versión | En desarrollo |

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
    service, err := apianime.NewAnimeFlv()
    if err != nil {
        panic(err)
    }
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
            
            // Extraer URLs directas
            fmt.Println("⏳ Extrayendo URLs directas...")
            videos, err := extract.ExtractURL("streamwish", ctx, link.URL)
            if err != nil {
                fmt.Printf("❌ Error: %v\n", err)
                continue
            }
            
            for _, video := range videos {
                fmt.Printf("✅ Resolución: %s - URL: %s\n", video.Resolution, video.URL)
            }
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

La librería utiliza **Viper** para la gestión de configuración, permitiendo persistencia automática en un archivo JSON en el directorio de configuración del usuario (`.configApiAnime.json`).

### Opción 1: Configuración Programática

Puedes configurar la librería programáticamente antes de crear el servicio.

```go
import (
    "github.com/dst3v3n/api-anime/config"
    "github.com/dst3v3n/api-anime/types"
)

func main() {
    // Configurar caché
    err := config.SetConfig(types.ConfigCache{
        CacheEnabled:  true,
        CacheHost:     "redis.prod.com",
        CachePort:     6380,
        CachePassword: "your-password",
        CacheDB:       0,
        CacheTTL:      120,
    })
    
    if err != nil {
        panic(err)
    }

    service, err := apianime.NewAnimeFlv()
    // ...
}
```

### Opción 2: Configuración por Defecto

Si no se especifica ninguna configuración, la librería se inicializará con valores por defecto y buscará el archivo de configuración en el directorio del sistema.

```go
service, err := apianime.NewAnimeFlv() // Inicializa con valores por defecto o archivo existente
```

### Opciones de Configuración

| Campo | Tipo | Default | Descripción |
|--------|------|---------|-------------|
| `CacheEnabled` | bool | `false` | Activar/desactivar caché |
| `CacheHost` | string | `localhost` | Host de Valkey/Redis |
| `CachePort` | int | `6379` | Puerto del servidor |
| `CacheUsername` | string | `""` | Usuario (si aplica) |
| `CachePassword` | string | `""` | Contraseña (si aplica) |
| `CacheDB` | int | `0` | Base de datos Redis |
| `CacheTTL` | int | `60` | TTL en minutos |

---

## 💾 Sistema de Caché

### ¿Qué Se Cachea?

| Función | Clave de Caché | TTL Default |
|---------|--------|-------------|
| `SearchAnime` | `search-anime-{nombre}-page-{N}` | Configurable (default 60min) |
| `Search` | `search-anime-all` | Configurable (default 60min) |
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
config.SetConfig(types.ConfigCache{CacheEnabled: false})
service, _ := apianime.NewAnimeFlv()

// Reactivar
config.SetConfig(types.ConfigCache{CacheEnabled: true})
service, _ = apianime.NewAnimeFlv()
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
config.SetConfig(types.ConfigCache{CacheEnabled: true})
```

**P: ¿Puedo cambiar el TTL?**  
R: Sí:

```go
config.SetConfig(types.ConfigCache{CacheTTL: 120}) // 2 horas
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
