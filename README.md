# Anime API

Librería Go para obtener información de animes mediante scraping de TioAnime, con caché distribuido opcional sobre Valkey (tipo JSON nativo).

> [!NOTE]
> AnimeFlv fue eliminado. Las versiones anteriores de esta librería usaban AnimeFlv como fuente, pero ese servicio dejó de estar disponible, por lo que su soporte se retiró por completo. La fuente actual es TioAnime. La arquitectura permite añadir otras fuentes sin cambiar la API pública.

---

## Contenido

1. [Descripción](#descripción)
2. [Instalación](#instalación)
3. [Puesta en marcha](#puesta-en-marcha)
4. [Uso](#uso)
5. [Configuración](#configuración)
6. [Caché con Valkey y JSON](#caché-con-valkey-y-json)
7. [Cómo funciona](#cómo-funciona)
8. [Pruebas](#pruebas)
9. [Preguntas frecuentes](#preguntas-frecuentes)
10. [Aviso legal](#aviso-legal)
11. [Licencia y contribuciones](#licencia-y-contribuciones)

---

## Descripción

La librería se inicializa con `apianime.NewTioanime()` y ofrece:

- Búsqueda de animes con paginación.
- Información detallada de cada anime.
- Servidores de reproducción disponibles por episodio.
- Listados de animes y episodios recientes.
- Caché opcional en Valkey que reduce una consulta de cientos de milisegundos a aproximadamente un milisegundo.

Si el caché está deshabilitado o falla, la librería sigue funcionando consultando directamente la fuente.

---

## Instalación

### Requisitos

| Componente | Versión | Obligatorio |
|---|---|---|
| Go | 1.25.3 o superior | Sí |
| Valkey con el módulo JSON | 8.x o superior | Solo si activas el caché |

### Instalar la librería

```bash
go get github.com/dst3v3n/api-anime
```

### Instalar Valkey con soporte JSON (solo si usas caché)

El caché guarda las respuestas con el tipo JSON nativo de Valkey, por lo que el servidor necesita el módulo `valkey-json`. Consulta la sección [Caché con Valkey y JSON](#caché-con-valkey-y-json) para las opciones de instalación.

---

## Puesta en marcha

1. Instale la librería (ver [Instalación](#instalación)).
2. Si va a usar caché, inicie un Valkey con el módulo JSON:

```bash
docker run -d --name valkey -p 6379:6379 valkey/valkey-bundle
```

3. Active el caché en el archivo de configuración (`~/.config/.configApiAnime.json`):

```json
{
  "cacheconfig": {
    "enablecache": true,
    "cachehost": "localhost",
    "cacheport": 6379,
    "cachedb": 0,
    "cachettl": 60
  }
}
```

4. Cree el servicio y realice consultas como se describe en [Uso](#uso).
5. Verifique que las entradas se guardan con los comandos de la sección [Verificación](#verificación).

---

## Uso

### Inicialización

El punto de entrada es `apianime.NewTioanime()`, que devuelve un `*Tioanime` listo para consultar.

```go
package main

import (
 "context"
 "fmt"
 "log"

 apianime "github.com/dst3v3n/api-anime"
)

func main() {
 svc, err := apianime.NewTioanime()
 if err != nil {
  log.Fatal(err)
 }

 ctx := context.Background()

 resultado, err := svc.Search(ctx, "naruto", 1)
 if err != nil {
  log.Fatal(err)
 }

 for _, a := range resultado.Animes {
  fmt.Printf("%s (%s)\n", a.Title, a.ID)
 }
}
```

`NewTioanime()` realiza tres tareas:

1. Carga la configuración desde el archivo `.configApiAnime.json` y lo crea con valores por defecto si no existe.
2. Construye el caché: un adaptador Valkey si `enablecache` es `true`, o uno que no almacena nada si es `false`.
3. Crea el cliente de TioAnime con su límite de peticiones por segundo.

Devuelve un error si la configuración no puede cargarse o si el caché está habilitado y no es posible conectar con Valkey. La inicialización debe hacerse una sola vez; el mismo `*Tioanime` se reutiliza en toda la aplicación.

Todos los métodos reciben un `context.Context` como primer argumento, que permite cancelar la consulta o aplicar un tiempo máximo.

### Métodos

| Método | Descripción | Retorno |
|---|---|---|
| `Search(ctx, anime, page)` | Busca animes por nombre, con paginación | `types.Anime` |
| `AnimeInfo(ctx, idAnime)` | Información detallada de un anime | `types.AnimeInfo` |
| `EpisodeStreamingProvider(ctx, idAnime, episode)` | Servidores de reproducción de un episodio | `types.EpisodeStreamingProvider` |
| `RecentAnime(ctx)` | Animes agregados recientemente | `types.Anime` |
| `RecentEpisode(ctx)` | Episodios publicados recientemente | `[]types.RecentEpisode` |

Todos devuelven además un `error`. Los tipos de respuesta se importan desde `github.com/dst3v3n/api-anime/types`.

#### Search

```go
Search(ctx context.Context, anime string, page uint) (types.Anime, error)
```

Busca por nombre. `page` comienza en 1.

```go
resultado, err := svc.Search(ctx, "one piece", 1)
if err != nil {
 return err
}

fmt.Printf("Resultados en la página: %d\n", len(resultado.Animes))
for _, a := range resultado.Animes {
 fmt.Println(a.ID, a.Title)
}
```

El campo `ID` de cada resultado es el identificador que se usa en `AnimeInfo` y `EpisodeStreamingProvider`.

#### AnimeInfo

```go
AnimeInfo(ctx context.Context, idAnime string) (types.AnimeInfo, error)
```

El identificador se normaliza (minúsculas y sin espacios laterales) antes de consultar. Un identificador vacío devuelve un error de entrada inválida sin llegar a consultar la fuente.

```go
info, err := svc.AnimeInfo(ctx, "one-piece-tv")
if err != nil {
 return err
}

fmt.Println(info.Title, info.Status)
fmt.Println(info.Genres)
fmt.Printf("Episodios: del %d al %d\n", info.FirstEpisode, info.LastEpisode)

if info.NextEpisode != nil {
 fmt.Println("Próximo episodio:", *info.NextEpisode)
}
if info.AnimeRelated != nil {
 for _, r := range *info.AnimeRelated {
  fmt.Printf("%s (%s)\n", r.Title, r.Category)
 }
}
```

#### EpisodeStreamingProvider

```go
EpisodeStreamingProvider(ctx context.Context, idAnime string, episode uint) (types.EpisodeStreamingProvider, error)
```

El identificador no puede estar vacío y `episode` debe ser mayor que cero; de lo contrario se devuelve un error de entrada inválida sin consultar la fuente.

```go
ep, err := svc.EpisodeStreamingProvider(ctx, "one-piece-tv", 1)
if err != nil {
 return err
}

fmt.Printf("%s, episodio %d\n", ep.Title, ep.Episode)
for _, p := range ep.Provider {
 fmt.Printf("%s: %s\n", p.Server, p.URL)
}
```

Cada elemento de `Provider` indica el nombre del servidor de video y su URL de reproducción.

#### RecentAnime y RecentEpisode

```go
RecentAnime(ctx context.Context) (types.Anime, error)
RecentEpisode(ctx context.Context) ([]types.RecentEpisode, error)
```

```go
animes, err := svc.RecentAnime(ctx)
episodios, err := svc.RecentEpisode(ctx)

for _, e := range episodios {
 fmt.Printf("%s, episodio %d\n", e.Title, e.Episode)
}
```

### Tipos de respuesta

```go
type Anime struct {
 Animes     []AnimeStruct
 TotalPages uint
 Source     string
}

type AnimeStruct struct {
 ID          string
 Title       string
 Image       string
 Synopsis    *string
 Type        *string
 Punctuation *float64
}

type AnimeInfo struct {
 AnimeStruct
 Genres       []string
 Status       string
 FirstEpisode int
 LastEpisode  int
 AnimeRelated *[]AnimeRelated
 NextEpisode  *string
 Source       string
}

type AnimeRelated struct {
 ID       string
 Title    string
 Category string
}

type EpisodeStreamingProvider struct {
 ID       string
 Title    string
 Episode  uint
 Provider []Provider
 Source   string
}

type Provider struct {
 Server string
 URL    string
}

type RecentEpisode struct {
 ID      string
 Title   string
 Episode uint
 Image   string
 Source  string
}
```

Notas sobre los tipos:

- Los campos con puntero (`Synopsis`, `Type`, `Punctuation`, `NextEpisode`, `AnimeRelated`) son opcionales. Compruebe que no sean `nil` antes de usarlos.
- `Image` contiene la ruta de la portada tal como la publica el sitio de origen. Si es una ruta relativa, debe anteponerse el dominio de la fuente.
- `Source` identifica el origen de los datos (por ejemplo, `tioanime`).
- `TotalPages` es `0` en los listados que no paginan, como `RecentAnime`.
- Los tipos se serializan a JSON con nombres en minúscula y guion bajo (`id`, `title`, `total_pages`, `first_episode`), por lo que pueden usarse directamente en respuestas HTTP.

### Manejo de errores

Todos los métodos devuelven errores normalizados: cada error indica el servicio, la operación y una categoría (`kind`) que clasifica la causa, y conserva el error original para inspeccionarlo con `errors.Is` y `errors.As`. El mensaje tiene el formato `operación: categoría: causa`, por ejemplo:

```
detailService.animeInfo: invalid_input: anime id is empty
```

| Categoría | Cuándo ocurre |
|---|---|
| `not_found` | El recurso no existe en la fuente |
| `invalid_input` | Parámetros inválidos: identificador vacío o episodio igual a `0` |
| `timeout` | El contexto superó su tiempo máximo (`context.DeadlineExceeded`) |
| `cancel` | El contexto fue cancelado (`context.Canceled`) |
| `rate_limited` | La fuente rechazó la petición por exceso de solicitudes |
| `config` | La configuración no pudo cargarse o es inválida |
| `external` | Fallo de red, de la fuente o del HTML recibido (cualquier error no clasificado) |

```go
resultado, err := svc.AnimeInfo(ctx, id)
switch {
case errors.Is(err, context.DeadlineExceeded):
 // la consulta tardó más que el tiempo máximo del contexto
case errors.Is(err, context.Canceled):
 // el llamador canceló la consulta
case err != nil:
 // cualquier otro error: validación, red o cambio en la fuente
}
```

Los casos habituales son:

- Identificador vacío o episodio igual a `0`.
- Identificador inexistente en el sitio.
- Fallo de red o tiempo de espera agotado en el contexto.
- Cambio en la estructura HTML del sitio, que impide interpretar la página.

Un fallo del caché nunca llega al llamador: se registra una advertencia y la consulta continúa contra la fuente.

> [!NOTE]
> Los errores se envuelven una sola vez: si un error ya está normalizado, las capas superiores lo devuelven sin modificar su categoría. Los paquetes de errores son internos (`internal/`), por lo que desde fuera de la librería la categoría se lee en el mensaje y la causa se inspecciona con `errors.Is`/`errors.As` sobre los errores estándar (por ejemplo, los del paquete `context`).

---

## Configuración

La configuración se gestiona con Viper y se persiste en un archivo JSON dentro del directorio de configuración del usuario.

| Sistema | Ruta |
|---|---|
| Linux | `~/.config/.configApiAnime.json` |
| macOS | `~/Library/Application Support/.configApiAnime.json` |
| Windows | `%AppData%\.configApiAnime.json` |

Si el archivo no existe, la librería lo crea con valores por defecto la primera vez que se inicializa. Tenga en cuenta este efecto secundario si la usa en entornos con sistema de archivos de solo lectura.

### Archivo de configuración

```json
{
  "appname": "Api-Anime",
  "cacheconfig": {
    "cachedb": 0,
    "cachehost": "localhost",
    "cachepassword": "",
    "cacheport": 6379,
    "cachettl": 60,
    "cacheusername": "",
    "enablecache": false
  },
  "logconfig": {
    "logappname": "Api-Anime",
    "logenv": "development"
  }
}
```

### Opciones de caché

| Clave | Tipo | Valor por defecto | Descripción |
|---|---|---|---|
| `enablecache` | bool | `false` | Activa o desactiva el caché |
| `cachehost` | string | `localhost` | Host del servidor Valkey |
| `cacheport` | int | `6379` | Puerto del servidor |
| `cachedb` | int | `0` | Base de datos lógica |
| `cacheusername` | string | vacío | Usuario, si el servidor usa ACL |
| `cachepassword` | string | vacío | Contraseña, si aplica |
| `cachettl` | int | `60` | Tiempo de vida de cada entrada, en minutos |

---

## Caché con Valkey y JSON

El caché guarda cada respuesta como un documento JSON nativo de Valkey, usando los comandos `JSON.SET` y `JSON.GET`.

> [!IMPORTANT]
> Estos comandos no forman parte del Valkey estándar: los aporta el módulo `valkey-json`. Si el servidor no lo tiene cargado, el caché no funciona. Instale el módulo antes de activar `enablecache`.

Si el caché está habilitado y el servidor no tiene el módulo, cada operación de caché falla con errores como:

```
unknown command 'JSON.SET'
unknown command 'JSON.GET'
```

La librería sigue respondiendo (consulta la fuente en cada llamada), pero el caché no tiene ningún efecto y se registra una advertencia por cada operación fallida.

### Cómo se almacena cada entrada

1. La respuesta se serializa y se guarda con `JSON.SET <clave> $ <documento>`.
2. El tiempo de vida se aplica con `EXPIRE <clave> <segundos>`, calculado a partir de `cachettl`.
3. Ambos comandos se envían dentro de una transacción `MULTI`/`EXEC`, de modo que la clave nunca queda guardada sin expiración.
4. La lectura usa `JSON.GET <clave>` y deserializa el documento directamente sobre el tipo de respuesta.

### Claves utilizadas

| Consulta | Clave |
|---|---|
| Búsqueda | `search-{consulta}-{página}` |
| Información de anime | `anime-info-{id}` |
| Servidores de un episodio | `episode-streaming-{id}-{episodio}` |
| Animes recientes | `recent-anime` |
| Episodios recientes | `recent-episode` |

### Opción A: Docker con la imagen bundle (recomendada)

La imagen `valkey-bundle` incluye el módulo JSON ya cargado.

```bash
docker run -d --name valkey \
  -p 6379:6379 \
  -v valkey-data:/data \
  valkey/valkey-bundle

docker exec -it valkey valkey-cli MODULE LIST
```

La salida de `MODULE LIST` debe incluir un módulo llamado `json`.

### Opción B: Compilar el módulo para una instalación existente

Si ya tiene Valkey instalado desde el gestor de paquetes de su sistema, puede compilar el módulo y cargarlo.

1. Instale las herramientas de compilación (por ejemplo, en Arch Linux):

```bash
sudo pacman -S --needed base-devel cmake git
```

2. Clone y compile:

```bash
git clone https://github.com/valkey-io/valkey-json.git
cd valkey-json
./build.sh
```

El resultado queda en `build/src/libjson.so`.

3. Copie la biblioteca a una ruta estable:

```bash
sudo mkdir -p /usr/local/lib/valkey
sudo cp build/src/libjson.so /usr/local/lib/valkey/
```

4. Añada la carga del módulo en `/etc/valkey/valkey.conf`:

```
loadmodule /usr/local/lib/valkey/libjson.so
```

5. Reinicie el servicio y verifique:

```bash
sudo systemctl restart valkey
valkey-cli MODULE LIST
```

Consideraciones:

- Compile el módulo contra una versión de Valkey compatible con la que ejecuta su servidor. Si el módulo se rechaza al cargar, el registro del servicio (`journalctl -u valkey`) indicará el motivo.
- Una actualización del paquete de Valkey puede dejar el módulo incompatible. Si el servicio no arranca tras actualizar, comente la línea `loadmodule` para recuperarlo y recompile el módulo.
- No use `MODULE LOAD` desde `valkey-cli` para una carga permanente. Valkey lo bloquea por defecto (`enable-module-command no`) y la carga no persistiría tras un reinicio.

### Verificación

Con la caché habilitada, ejecute dos veces cualquier consulta y revise:

```bash
valkey-cli KEYS '*'
valkey-cli TTL recent-anime
valkey-cli JSON.GET recent-anime
```

- `KEYS` debe listar las claves creadas.
- `TTL` debe devolver un valor positivo cercano a `cachettl * 60`. Un `-1` indica una clave sin expiración.
- `JSON.GET` debe devolver el documento almacenado.

### Rendimiento orientativo

| Situación | Tiempo típico |
|---|---|
| Consulta al sitio (sin caché o primera vez) | Cientos de milisegundos |
| Lectura desde Valkey | Alrededor de 1 ms |

Las cifras dependen de la red, del sitio de origen y del rate limiting. Los valores anteriores provienen de mediciones en un entorno local.

### Comportamiento ante fallos

El caché es una optimización y no una dependencia:

- Una clave inexistente se reporta internamente como un fallo de caché (`ErrCacheMiss`), distinto de un error real de Valkey, y la consulta continúa contra la fuente sin registrar advertencias.
- Si la lectura falla (servidor caído, documento corrupto), se registra una advertencia y la consulta continúa contra la fuente.
- Si la escritura falla (incluido un valor que no puede serializarse a JSON), se registra una advertencia y el resultado se devuelve igualmente al llamador.
- Solo un error de la fuente de datos llega al usuario.
- Un fallo al conectar con Valkey durante la creación del servicio sí se devuelve como error, porque indica un problema de configuración.

### Alternativa sin módulo JSON

Si no desea instalar `valkey-json`, deje `enablecache` en `false`. La librería funciona igual, sin caché.

---

## Cómo funciona

El proyecto sigue una arquitectura hexagonal (puertos y adaptadores):

```
tioanime.go                  API pública (paquete apianime)
types/                       Alias públicos de los DTO

internal/
  domain/
    dto/                     Estructuras de datos de respuesta
    services/
      shared/                Patrón cache-aside y utilidades
      tioanime/              Servicios: búsqueda, detalle y recientes
    errors/                  Errores normalizados (Kind), sentinelas y Wrap
  ports/                     Contratos: CachePort, ScraperPort, ...
  adapters/
    scrapers/tioanime/       Cliente y parsers de TioAnime
    cache/                   Adaptadores de caché (Valkey y no-op)
    videoextractor/          Extractores de URLs de video
  config/                    Carga de configuración y logging
  mocks/                     Mocks de los puertos para las pruebas

test/
  unit/                      Pruebas unitarias (caché, extractores, parsers)
  integration/tioanime/      Pruebas de integración de servicio y API pública
```

### Flujo de una consulta

1. El método público delega en el servicio correspondiente.
2. El servicio normaliza el identificador (minúsculas y sin espacios) y construye la clave de caché.
3. El patrón cache-aside consulta el caché. Si hay un acierto, devuelve el resultado.
4. Si hay un fallo de caché, ejecuta el scraper, que respeta un límite de peticiones por segundo.
5. El resultado se guarda en el caché con su tiempo de vida y se devuelve.
6. Cualquier error se normaliza con su servicio, operación y categoría antes de llegar al llamador.

### Decisiones de diseño

- **Puertos pequeños por capacidad.** Los servicios dependen de interfaces, no de implementaciones concretas. Esto permite probar con mocks y añadir nuevas fuentes sin modificar los servicios.
- **Caché deshabilitado como implementación no-op.** Cuando `enablecache` es `false`, se inyecta un adaptador que siempre devuelve un fallo de caché y descarta las escrituras. Los servicios no necesitan comprobar nada.
- **Errores normalizados.** Un tipo de error propio (`Error`) con servicio, operación y categoría (`Kind`) reemplaza a los `fmt.Errorf` dispersos. `Wrap` clasifica el error original (contexto cancelado, tiempo agotado, fallo externo) y evita envolverlo dos veces, y los errores sentinela (`ErrEmptyAnimeID`, `ErrInvalidEpisode`) describen las validaciones de entrada.
- **Fallo de caché diferenciado.** El puerto de caché expone `ErrCacheMiss` para distinguir una clave inexistente de un error real de Valkey.
- **Rate limiting.** El cliente limita las peticiones a la fuente para no sobrecargarla.

---

## Pruebas

```bash
# Todas las pruebas
go test ./...

# Solo pruebas unitarias
go test ./test/unit/...

# Solo pruebas de integración
go test ./test/integration/...

# Con cobertura
go test ./... -cover

# Una prueba concreta
go test ./... -v -run TestServiceSearch
```

| Tipo | Ubicación | Qué cubre |
|---|---|---|
| Unitarias | `test/unit/` | Caché (`Set`/`Get`), extractores de video StreamStape y StreamWish, y parsers de búsqueda, detalle y servidores de reproducción |
| Integración | `test/integration/tioanime/` | Servicios y API pública (`Search`, `AnimeInfo`, `EpisodeStreamingProvider`, recientes y cancelación de contexto) |

Los mocks de los puertos están en `internal/mocks/mocks_tioanime.go`.

> [!NOTE]
> Las pruebas de integración consultan el sitio real de TioAnime, por lo que requieren conexión a internet y pueden fallar si la estructura HTML cambia. Las que usan caché además necesitan un Valkey con el módulo JSON en ejecución.

---

## Preguntas frecuentes

**¿Es obligatorio usar Valkey?**
No. El caché es opcional y la librería funciona sin él.

**¿Funciona con Redis?**
El caché usa comandos `JSON.*`. Redis los ofrece mediante el módulo RedisJSON, pero esta librería está probada únicamente con Valkey y `valkey-json`.

**¿Por qué recibo advertencias `unknown command 'JSON.SET'`?**
El servidor no tiene el módulo JSON cargado. Consulte la sección [Caché con Valkey y JSON](#caché-con-valkey-y-json), o desactive el caché.

**¿Qué ocurre con AnimeFlv?**
Fue eliminado de la librería porque el servicio dejó de estar disponible.

**¿Los enlaces de reproducción caducan?**
Algunos servidores emiten URLs temporales. Ajuste `cachettl` a un valor menor si necesita enlaces más frescos.

**¿Puedo usarla en producción?**
El scraping depende de la estructura HTML del sitio de origen, que puede cambiar sin aviso. Monitorice los errores de las consultas y mantenga la librería actualizada.

**¿Por qué se crea un archivo en mi directorio de configuración?**
La inicialización genera `.configApiAnime.json` con valores por defecto si no existe. Es el mecanismo de persistencia de la configuración.

---

## Aviso legal

Este proyecto tiene fines educativos y personales. Respete los términos de servicio y el archivo `robots.txt` del sitio de origen.

Se espera que quien use la librería:

- Aplique límites de peticiones razonables.
- No sobrecargue el servicio de origen ni realice scraping masivo.
- Cite la fuente de los datos.
- No la use con fines comerciales sin autorización de los titulares correspondientes.

Los autores no se hacen responsables del uso que se haga de la librería.

---

## Licencia y contribuciones

Distribuido bajo licencia MIT. Consulte el archivo [LICENSE](LICENSE).

Las contribuciones son bienvenidas. Consulte [CONTRIBUTING.md](CONTRIBUTING.md) antes de abrir un pull request.

- Reportar errores: [Issues](../../issues)
- Proponer mejoras: [Discussions](../../discussions)
- Enviar código: [Pull Requests](../../pulls)

Autor: Steven ([@dst3v3n](https://github.com/dst3v3n))
