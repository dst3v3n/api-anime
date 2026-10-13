package tioanime

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
	domainErrors "github.com/dst3v3n/api-anime/internal/domain/errors"
)

var (
	episodesRegex    = regexp.MustCompile(`var episodes = (\[.*?\]);`)
	animeInfoRegex   = regexp.MustCompile(`var anime_info = (\[.*?\]);`)
	linkEpisodeRegex = regexp.MustCompile(`var videos = (\[.*?\]);`)
)

func scriptEpisodeList(scriptContent string) (int, int, error) {
	const op = "scriptEpisodeList"

	matches := episodesRegex.FindStringSubmatch(scriptContent)
	if len(matches) < 2 {
		return 0, 0, domainErrors.E(service, op, domainErrors.KindExternal,
			errors.New("var episodes no encontrada en el script"))
	}

	var episodes []int
	if err := json.Unmarshal([]byte(matches[1]), &episodes); err != nil {
		return 0, 0, domainErrors.E(service, op, domainErrors.KindExternal, err)
	}

	if len(episodes) == 0 {
		return 0, 0, domainErrors.E(service, op, domainErrors.KindExternal,
			errors.New("la lista de episodios está vacía"))
	}

	lastEpisode := episodes[0]
	firstEpisode := episodes[len(episodes)-1]

	return firstEpisode, lastEpisode, nil
}

func scriptInfo(scriptContent string) (string, error) {
	const op = "scriptInfo"

	matches := animeInfoRegex.FindStringSubmatch(scriptContent)
	if len(matches) < 2 {
		return "", domainErrors.E(service, op, domainErrors.KindExternal,
			errors.New("var anime_info no encontrada en el script"))
	}

	var info []string
	if err := json.Unmarshal([]byte(matches[1]), &info); err != nil {
		return "", domainErrors.E(service, op, domainErrors.KindExternal, err)
	}

	if len(info) < 4 {
		return "", nil
	}

	return info[3], nil
}

func scriptProvider(scriptContent string) ([]dto.Provider, error) {
	const op = "scriptLinksEpisode"
	var provider []dto.Provider

	matches := linkEpisodeRegex.FindStringSubmatch(scriptContent)
	if len(matches) < 2 {
		return nil, domainErrors.E(service, op, domainErrors.KindExternal,
			errors.New("var videos no encontrada en el script"))
	}

	var videos []video
	if err := json.Unmarshal([]byte(matches[1]), &videos); err != nil {
		return nil, domainErrors.E(service, op, domainErrors.KindExternal, err)
	}

	for _, v := range videos {
		provider = append(provider, dto.Provider{
			Server: v.Server,
			URL:    v.URL,
		})
	}

	return provider, nil
}
