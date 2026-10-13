package mocks

import "github.com/dst3v3n/api-anime/internal/domain/dto"

const (
	MockSource = "tioanime"

	StatusOnAir  = "En emisión"
	StatusFinish = "Finalizado"
)

func ptr[T any](v T) *T {
	return &v
}

func MockAnimeStruct() dto.AnimeStruct {
	return dto.AnimeStruct{
		AnimeSummary: dto.AnimeSummary{
			ID:    "one-piece-tv",
			Title: "One Piece",
			Image: "/uploads/portadas/1.jpg",
		},
		Synopsis:    ptr("Las aventuras de Monkey D. Luffy y su tripulación en busca del tesoro más grande del mundo, el One Piece, para convertirse en el Rey de los Piratas."),
		Type:        ptr("Anime"),
		Punctuation: ptr(8.9),
	}
}

func MockAnimeStructMinimal() dto.AnimeStruct {
	return dto.AnimeStruct{
		AnimeSummary: dto.AnimeSummary{
			ID:    "anime-sin-datos",
			Title: "Anime sin datos",
			Image: "/uploads/portadas/0.jpg",
		},
	}
}

func MockAnimeStructList() []dto.AnimeStruct {
	return []dto.AnimeStruct{
		MockAnimeStruct(),
		{
			AnimeSummary: dto.AnimeSummary{
				ID:    "naruto-shippuden",
				Title: "Naruto: Shippuden",
				Image: "/uploads/portadas/2.jpg",
			},
			Synopsis:    ptr("Naruto regresa después de dos años de entrenamiento para enfrentar nuevas amenazas."),
			Type:        ptr("Anime"),
			Punctuation: ptr(8.7),
		},
		{
			AnimeSummary: dto.AnimeSummary{
				ID:    "shingeki-no-kyojin-the-final-season",
				Title: "Shingeki no Kyojin: The Final Season",
				Image: "/uploads/portadas/3.jpg",
			},
			Synopsis:    ptr("La temporada final de la guerra entre la humanidad y los titanes."),
			Type:        ptr("Anime"),
			Punctuation: ptr(9.1),
		},
		{
			AnimeSummary: dto.AnimeSummary{
				ID:    "kimetsu-no-yaiba-mugen-ressha-hen",
				Title: "Kimetsu no Yaiba: Mugen Ressha-hen",
				Image: "/uploads/portadas/4.jpg",
			},
			Synopsis:    ptr("Tanjiro y sus amigos abordan el Tren Infinito para investigar desapariciones."),
			Type:        ptr("Película"),
			Punctuation: ptr(8.8),
		},
		{
			AnimeSummary: dto.AnimeSummary{
				ID:    "steins-gate-oukoubakko-no-poriomania",
				Title: "Steins;Gate: Oukoubakko no Poriomania",
				Image: "/uploads/portadas/5.jpg",
			},
			Synopsis:    ptr("Un episodio especial que muestra la vida cotidiana después de los eventos de Steins;Gate."),
			Type:        ptr("OVA"),
			Punctuation: ptr(8.4),
		},
	}
}

func MockAnime() dto.Anime {
	return dto.Anime{
		Animes:     MockAnimeStructList(),
		TotalPages: 15,
		Source:     MockSource,
	}
}

func MockAnimeSinglePage() dto.Anime {
	return dto.Anime{
		Animes:     MockAnimeStructList()[:3],
		TotalPages: 1,
		Source:     MockSource,
	}
}

func MockAnimeRecent() dto.Anime {
	return dto.Anime{
		Animes:     MockAnimeStructList(),
		TotalPages: 0,
		Source:     MockSource,
	}
}

func MockEmptyAnime() dto.Anime {
	return dto.Anime{
		Animes:     []dto.AnimeStruct{},
		TotalPages: 0,
		Source:     MockSource,
	}
}

func MockAnimeInfo() dto.AnimeInfo {
	return dto.AnimeInfo{
		AnimeStruct: MockAnimeStruct(),
		AnimeInfoSummary: dto.AnimeInfoSummary{
			Genres:       []string{"Acción", "Aventura", "Comedia", "Drama", "Fantasía", "Shounen"},
			Status:       StatusOnAir,
			FirstEpisode: 1,
			LastEpisode:  1090,
		},
		AnimeRelated: &[]dto.AnimeRelated{
			{
				ID:       "one-piece-film-red",
				Title:    "One Piece Film: Red",
				Category: "Película",
			},
			{
				ID:       "one-piece-3d-mugiwara-chase",
				Title:    "One Piece 3D: Mugiwara Chase",
				Category: "Especial",
			},
		},
		NextEpisode: ptr("2026-10-18"),
		Source:      MockSource,
	}
}

func MockAnimeInfoFinished() dto.AnimeInfo {
	return dto.AnimeInfo{
		AnimeStruct: dto.AnimeStruct{
			AnimeSummary: dto.AnimeSummary{
				ID:    "fullmetal-alchemist-brotherhood",
				Title: "Fullmetal Alchemist: Brotherhood",
				Image: "/uploads/portadas/7.jpg",
			},
			Synopsis:    ptr("Dos hermanos alquimistas buscan la Piedra Filosofal para recuperar sus cuerpos."),
			Type:        ptr("Anime"),
			Punctuation: ptr(9.2),
		},
		AnimeInfoSummary: dto.AnimeInfoSummary{
			Genres:       []string{"Acción", "Aventura", "Drama", "Fantasía", "Magia", "Militar", "Shounen"},
			Status:       StatusFinish,
			FirstEpisode: 1,
			LastEpisode:  64,
		},
		AnimeRelated: &[]dto.AnimeRelated{
			{
				ID:       "fullmetal-alchemist",
				Title:    "Fullmetal Alchemist",
				Category: "Precuela",
			},
		},
		NextEpisode: nil,
		Source:      MockSource,
	}
}

func MockAnimeInfoMinimal() dto.AnimeInfo {
	return dto.AnimeInfo{
		AnimeStruct: MockAnimeStructMinimal(),
		AnimeInfoSummary: dto.AnimeInfoSummary{
			Genres:       []string{},
			Status:       StatusFinish,
			FirstEpisode: 1,
			LastEpisode:  1,
		},
		Source: MockSource,
	}
}

func MockRecentEpisodes() []dto.RecentEpisode {
	return []dto.RecentEpisode{
		{
			ID:      "one-piece-tv",
			Title:   "One Piece",
			Episode: 1090,
			Image:   "/uploads/portadas/1.jpg",
			Source:  MockSource,
		},
		{
			ID:      "jujutsu-kaisen-2nd-season",
			Title:   "Jujutsu Kaisen 2nd Season",
			Episode: 23,
			Image:   "/uploads/portadas/8.jpg",
			Source:  MockSource,
		},
		{
			ID:      "spy-x-family-season-2",
			Title:   "Spy x Family Season 2",
			Episode: 12,
			Image:   "/uploads/portadas/9.jpg",
			Source:  MockSource,
		},
		{
			ID:      "sousou-no-frieren",
			Title:   "Sousou no Frieren",
			Episode: 28,
			Image:   "/uploads/portadas/10.jpg",
			Source:  MockSource,
		},
	}
}

func MockEmptyRecentEpisodes() []dto.RecentEpisode {
	return []dto.RecentEpisode{}
}

func MockEpisodeStreamingProvider() dto.EpisodeStreamingProvider {
	return dto.EpisodeStreamingProvider{
		ID:      "one-piece-tv",
		Title:   "One Piece",
		Episode: 1090,
		Provider: []dto.Provider{
			{Server: "Streamwish", URL: "https://streamwish.to/e/abc123def456"},
			{Server: "Streamtape", URL: "https://streamtape.com/e/xyz789uvw"},
			{Server: "Mega", URL: "https://mega.nz/embed/qrs456tuv"},
			{Server: "Okru", URL: "https://ok.ru/videoembed/123456789"},
		},
		Source: MockSource,
	}
}

func MockEpisodeStreamingProviderList() []dto.EpisodeStreamingProvider {
	return []dto.EpisodeStreamingProvider{
		MockEpisodeStreamingProvider(),
		{
			ID:      "one-piece-tv",
			Title:   "One Piece",
			Episode: 1089,
			Provider: []dto.Provider{
				{Server: "Streamwish", URL: "https://streamwish.to/e/prev123"},
				{Server: "Mega", URL: "https://mega.nz/embed/prev456"},
			},
			Source: MockSource,
		},
	}
}

func MockEpisodeStreamingProviderEmpty() dto.EpisodeStreamingProvider {
	return dto.EpisodeStreamingProvider{
		ID:       "one-piece-tv",
		Title:    "One Piece",
		Episode:  1091,
		Provider: []dto.Provider{},
		Source:   MockSource,
	}
}
