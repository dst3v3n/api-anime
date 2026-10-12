package tioanime

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/shared"
	"github.com/dst3v3n/api-anime/internal/config"
	"github.com/dst3v3n/api-anime/internal/domain/dto"
	domainErrors "github.com/dst3v3n/api-anime/internal/domain/errors"
)

const (
	selectorSearchArticle         = "ul.list-unstyled > li > article.anime"
	selectorLink                  = "a"
	selectorTitle                 = "h3.title"
	selectorImg                   = "img"
	selectorPagination            = "nav > ul.pagination"
	selectorPaginationPenultimate = "li:nth-last-child(2) a"
	selectorAnimeSingle           = "article.anime-single"
	selectorAnimeTitle            = "h1.title"
	selectorAnimeType             = "span.anime-type-peli"
	selectorSynopsis              = "p.sinopsis"
	selectorStatus                = "a.status"
	selectorGenres                = "p.genres > span > a"
	selectorRelatedArticle        = "section.w-history article"
	selectorRelatedLink           = "div.media-body a"
	selectorScript                = "script"
	selectorEpisodeTitle          = "h1.anime-title"
	selectorRecentArticle         = "ul.episodes > li > article"
)

const source = "tioanime"

func parserSearch(htmlElement io.Reader) (dto.Anime, error) {
	const op = "parserSearch"
	log := config.Logging()

	doc, err := goquery.NewDocumentFromReader(htmlElement)
	if err != nil {
		return dto.Anime{}, domainErrors.E(service, op, domainErrors.KindExternal, err)
	}

	var animes []dto.AnimeStruct
	var pagination uint

	doc.Find(selectorSearchArticle).Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Find(selectorLink).Attr("href")
		id, err := shared.ExtractID(href)
		if err != nil {
			log.Warn().Err(domainErrors.E(service, op, domainErrors.KindExternal, err)).
				Str("href", href).Msg("ID no extraído, elemento omitido")
			return
		}

		title, _ := s.Find(selectorTitle).Html()
		image, _ := s.Find(selectorImg).Attr("src")

		animes = append(animes, dto.AnimeStruct{
			AnimeSummary: dto.AnimeSummary{
				ID:    id,
				Title: title,
				Image: image,
			},
		})
	})

	doc.Find(selectorPagination).Each(func(_ int, s *goquery.Selection) {
		text := s.Find(selectorPaginationPenultimate).Text()
		p, err := shared.ParseUint(text)
		if err != nil {
			p = 1
		}
		pagination = p
	})

	return dto.Anime{TotalPages: pagination, Animes: animes, Source: source}, nil
}

func parserAnimeInfo(htmlElement io.Reader, idAnime string) (dto.AnimeInfo, error) {
	const op = "parserAnimeInfo"
	log := config.Logging()

	anime := dto.AnimeInfo{}
	anime.ID = idAnime
	animeRelated := []dto.AnimeRelated{}

	doc, err := goquery.NewDocumentFromReader(htmlElement)
	if err != nil {
		return anime, domainErrors.E(service, op, domainErrors.KindExternal, err)
	}

	doc.Find(selectorAnimeSingle).Each(func(_ int, s *goquery.Selection) {
		anime.Title = s.Find(selectorAnimeTitle).Text()
		anime.Image, _ = s.Find(selectorImg).Attr("src")

		typeAnime := s.Find(selectorAnimeType).Text()

		synopsis := strings.TrimSpace(s.Find(selectorSynopsis).Text())
		synopsis = strings.ReplaceAll(synopsis, `"`, "")

		anime.Type = &typeAnime
		anime.Synopsis = &synopsis
		anime.Status = s.Find(selectorStatus).Text()

		var genres []string
		s.Find(selectorGenres).Each(func(_ int, genreSelect *goquery.Selection) {
			genres = append(genres, genreSelect.Text())
		})
		anime.Genres = genres
	})

	if anime.Title == "" {
		return anime, domainErrors.E(service, op, domainErrors.KindExternal,
			errors.New("no se encontró article.anime-single en el HTML"))
	}

	doc.Find(selectorRelatedArticle).Each(func(_ int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(selectorRelatedLink).Text())
		category := s.Find(selectorAnimeType).Text()
		href, _ := s.Find(selectorRelatedLink).Attr("href")

		id, err := shared.ExtractID(href)
		if err != nil {
			log.Warn().Err(domainErrors.E(service, op, domainErrors.KindExternal, err)).
				Str("href", href).Msg("ID no extraído, anime relacionado omitido")
			return
		}

		animeRelated = append(animeRelated, dto.AnimeRelated{
			ID:       id,
			Title:    title,
			Category: category,
		})
	})

	doc.Find(selectorScript).Each(func(_ int, s *goquery.Selection) {
		textHTML := s.Text()

		if strings.Contains(textHTML, "var episodes") {
			firstEpisode, lastEpisode, err := scriptEpisodeList(textHTML)
			if err != nil {
				log.Warn().Err(err).Str("anime", idAnime).Msg("lista de episodios no extraída")
				return
			}
			anime.FirstEpisode = firstEpisode
			anime.LastEpisode = lastEpisode
		}

		if strings.Contains(textHTML, "var anime_info") {
			nextEpisode, err := scriptInfo(textHTML)
			if err != nil {
				log.Warn().Err(err).Str("anime", idAnime).Msg("info del próximo episodio no extraída")
				return
			}
			if nextEpisode != "" {
				anime.NextEpisode = &nextEpisode
			}
		}
	})

	anime.AnimeRelated = &animeRelated
	anime.Source = source

	return anime, nil
}

func parserEpisodeStreaming(htmlElement io.Reader, idAnime string, episode uint) (dto.EpisodeStreamingProvider, error) {
	const op = "parserEpisodeStreaming"

	episodeStreaming := dto.EpisodeStreamingProvider{}
	episodeStreaming.ID = idAnime
	episodeStreaming.Episode = episode

	doc, err := goquery.NewDocumentFromReader(htmlElement)
	if err != nil {
		return episodeStreaming, domainErrors.E(service, op, domainErrors.KindExternal, err)
	}

	title := doc.Find(selectorEpisodeTitle).Text()
	episodeString := strconv.Itoa(int(episode))
	episodeStreaming.Title = strings.TrimSpace(strings.ReplaceAll(title, episodeString, ""))

	doc.Find(selectorScript).Each(func(_ int, s *goquery.Selection) {
		textHTML := s.Text()
		if strings.Contains(textHTML, "var videos") {
			provider, err := scriptProvider(textHTML)
			if err != nil {
				return
			}
			episodeStreaming.Provider = provider
		}
	})

	episodeStreaming.Source = source

	return episodeStreaming, nil
}

func parserRecentEpisode(htmlElement io.Reader) ([]dto.RecentEpisode, error) {
	const op = "parserRecentEpisode"

	log := config.Logging()
	animeRecent := []dto.RecentEpisode{}

	doc, err := goquery.NewDocumentFromReader(htmlElement)
	if err != nil {
		return animeRecent, domainErrors.E(service, op, domainErrors.KindExternal, err)
	}

	doc.Find(selectorRecentArticle).Each(func(_ int, s *goquery.Selection) {
		titleEpisode := s.Find(selectorTitle).Text()
		split := strings.Split(titleEpisode, " ")
		episode := split[len(split)-1]
		title := strings.TrimSpace(removeLast(titleEpisode, episode))

		href, _ := s.Find(selectorLink).Attr("href")
		href = strings.TrimSpace(removeLast(href, fmt.Sprintf("-%s", episode)))

		id, err := shared.ExtractID(href)
		if err != nil {
			log.Warn().Err(domainErrors.E(service, op, domainErrors.KindExternal, err)).
				Str("href", href).Msg("ID no extraído, elemento omitido")
			return
		}

		img, _ := s.Find(selectorImg).Attr("src")
		episodeUint, _ := shared.ParseUint(episode)

		animeRecent = append(animeRecent, dto.RecentEpisode{
			ID:      id,
			Title:   title,
			Episode: episodeUint,
			Image:   img,
			Source:  source,
		})
	})

	return animeRecent, nil
}
