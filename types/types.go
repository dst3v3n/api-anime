package types

import (
	"github.com/dst3v3n/api-anime/internal/config"
	"github.com/dst3v3n/api-anime/internal/domain/dto"
)

type Anime = dto.Anime

type AnimeStruct = dto.AnimeStruct

type AnimeInfo = dto.AnimeInfo

type EpisodeStreamingProvider = dto.EpisodeStreamingProvider

type VideoURL = dto.VideoURL

type RecentEpisode = dto.RecentEpisode

type ConfigCache = config.CacheConfig
