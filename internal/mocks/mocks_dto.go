// Package mocks contiene funciones para generar datos de prueba (mocks)
// de los DTOs de anime. Útiles para testing unitario e integración.
// Proporciona constructores que generan instancias de estructuras DTO con datos
// realistas para validar comportamientos sin hacer scraping real.
package mocks

import "github.com/dst3v3n/api-anime/v2/internal/domain/dto"

// MockAnimeStruct retorna un anime de ejemplo con datos completos.
func MockAnimeStruct() dto.AnimeStruct {
	return dto.AnimeStruct{
		ID:          "one-piece-tv",
		Title:       "One Piece",
		Sinopsis:    "Las aventuras de Monkey D. Luffy y su tripulación en busca del tesoro más grande del mundo, el One Piece, para convertirse en el Rey de los Piratas.",
		Type:        dto.Anime,
		Punctuation: 8.9,
		Image:       "https://cdn.myanimelist.net/images/anime/6/73245.jpg",
	}
}

// MockAnimeStructList retorna una lista de animes de ejemplo.
func MockAnimeStructList() []dto.AnimeStruct {
	return []dto.AnimeStruct{
		{
			ID:          "one-piece-tv",
			Title:       "One Piece",
			Sinopsis:    "Las aventuras de Monkey D. Luffy y su tripulación.",
			Type:        dto.Anime,
			Punctuation: 8.9,
			Image:       "https://cdn.myanimelist.net/images/anime/6/73245.jpg",
		},
		{
			ID:          "naruto-shippuden",
			Title:       "Naruto: Shippuden",
			Sinopsis:    "Naruto regresa después de dos años de entrenamiento para enfrentar nuevas amenazas.",
			Type:        dto.Anime,
			Punctuation: 8.7,
			Image:       "https://cdn.myanimelist.net/images/anime/5/17407.jpg",
		},
		{
			ID:          "attack-on-titan-final",
			Title:       "Shingeki no Kyojin: The Final Season",
			Sinopsis:    "La temporada final de la guerra entre la humanidad y los titanes.",
			Type:        dto.Anime,
			Punctuation: 9.1,
			Image:       "https://cdn.myanimelist.net/images/anime/1948/120625.jpg",
		},
		{
			ID:          "demon-slayer-movie",
			Title:       "Kimetsu no Yaiba: Mugen Ressha-hen",
			Sinopsis:    "Tanjiro y sus amigos abordan el Tren Infinito para investigar desapariciones.",
			Type:        dto.Pelicula,
			Punctuation: 8.8,
			Image:       "https://cdn.myanimelist.net/images/anime/1704/106947.jpg",
		},
		{
			ID:          "steins-gate-ova",
			Title:       "Steins;Gate: Oukoubakko no Poriomania",
			Sinopsis:    "Un episodio especial que muestra la vida cotidiana después de los eventos de Steins;Gate.",
			Type:        dto.Ova,
			Punctuation: 8.4,
			Image:       "https://cdn.myanimelist.net/images/anime/12/35643.jpg",
		},
	}
}

// MockAnimeResponse retorna una respuesta de búsqueda de animes con paginación.
func MockAnimeResponse() dto.AnimeResponse {
	return dto.AnimeResponse{
		Animes:     MockAnimeStructList(),
		TotalPages: 15,
	}
}

// MockAnimeInfoResponse retorna información detallada de un anime.
func MockAnimeInfoResponse() dto.AnimeInfoResponse {
	return dto.AnimeInfoResponse{
		AnimeStruct: MockAnimeStruct(),
		AnimeRelated: []dto.AnimeRelated{
			{
				ID:       "one-piece-film-red",
				Title:    "One Piece Film: Red",
				Category: "Película",
			},
			{
				ID:       "one-piece-special-3d",
				Title:    "One Piece 3D: Mugiwara Chase",
				Category: "Especial",
			},
		},
		Genres: []string{
			"Acción",
			"Aventura",
			"Comedia",
			"Drama",
			"Fantasía",
			"Shounen",
		},
		Status:      dto.Emision,
		NextEpisode: "2024-01-14",
		Episodes:    generateEpisodeNumbers(1, 1090),
	}
}

// MockAnimeInfoResponseFinished retorna un anime finalizado.
func MockAnimeInfoResponseFinished() dto.AnimeInfoResponse {
	return dto.AnimeInfoResponse{
		AnimeStruct: dto.AnimeStruct{
			ID:          "fullmetal-alchemist-brotherhood",
			Title:       "Fullmetal Alchemist: Brotherhood",
			Sinopsis:    "Dos hermanos alquimistas buscan la Piedra Filosofal para recuperar sus cuerpos.",
			Type:        dto.Anime,
			Punctuation: 9.2,
			Image:       "https://cdn.myanimelist.net/images/anime/1223/96541.jpg",
		},
		AnimeRelated: []dto.AnimeRelated{
			{
				ID:       "fullmetal-alchemist",
				Title:    "Fullmetal Alchemist",
				Category: "Precuela",
			},
		},
		Genres: []string{
			"Acción",
			"Aventura",
			"Drama",
			"Fantasía",
			"Magia",
			"Militar",
			"Shounen",
		},
		Status:      dto.Finalizado,
		NextEpisode: "",
		Episodes:    generateEpisodeNumbers(1, 64),
	}
}

// MockEpisodeListResponse retorna una lista de episodios recientes.
func MockEpisodeListResponse() []dto.EpisodeListResponse {
	return []dto.EpisodeListResponse{
		{
			ID:      "one-piece-tv",
			Title:   "One Piece",
			Chapter: "Cap. 1090",
			Episode: 1090,
			Image:   "https://cdn.myanimelist.net/images/anime/6/73245.jpg",
		},
		{
			ID:      "jujutsu-kaisen-2",
			Title:   "Jujutsu Kaisen 2nd Season",
			Chapter: "Cap. 23",
			Episode: 23,
			Image:   "https://cdn.myanimelist.net/images/anime/1792/138022.jpg",
		},
		{
			ID:      "spy-x-family-2",
			Title:   "Spy x Family Season 2",
			Chapter: "Cap. 12",
			Episode: 12,
			Image:   "https://cdn.myanimelist.net/images/anime/1111/127508.jpg",
		},
		{
			ID:      "frieren",
			Title:   "Sousou no Frieren",
			Chapter: "Cap. 28",
			Episode: 28,
			Image:   "https://cdn.myanimelist.net/images/anime/1015/138006.jpg",
		},
	}
}

// MockLinkResponse retorna enlaces de reproducción para un episodio.
func MockLinkResponse() dto.LinkResponse {
	return dto.LinkResponse{
		ID:      "one-piece-tv",
		Title:   "One Piece",
		Episode: 1090,
		Link: []dto.LinkSource{
			{
				Server: "StreamSB",
				URL:    "https://streamsb.net/e/abc123def456",
				Code:   "abc123def456",
			},
			{
				Server: "Mega",
				URL:    "https://mega.nz/file/xyz789uvw",
				Code:   "xyz789uvw",
			},
			{
				Server: "Fembed",
				URL:    "https://fembed.com/v/qrs456tuv",
				Code:   "qrs456tuv",
			},
			{
				Server: "Okru",
				URL:    "https://ok.ru/video/123456789",
				Code:   "123456789",
			},
			{
				Server: "Zippyshare",
				URL:    "https://www72.zippyshare.com/v/abcd1234/file.html",
				Code:   "abcd1234",
			},
		},
	}
}

// MockLinkResponseMultipleEpisodes retorna enlaces para múltiples episodios.
func MockLinkResponseMultipleEpisodes() []dto.LinkResponse {
	return []dto.LinkResponse{
		MockLinkResponse(),
		{
			ID:      "one-piece-tv",
			Title:   "One Piece",
			Episode: 1089,
			Link: []dto.LinkSource{
				{
					Server: "StreamSB",
					URL:    "https://streamsb.net/e/prev123",
					Code:   "prev123",
				},
				{
					Server: "Mega",
					URL:    "https://mega.nz/file/prev456",
					Code:   "prev456",
				},
			},
		},
	}
}

// MockEmptyAnimeResponse retorna una respuesta vacía (sin resultados).
func MockEmptyAnimeResponse() dto.AnimeResponse {
	return dto.AnimeResponse{
		Animes:     []dto.AnimeStruct{},
		TotalPages: 0,
	}
}

// MockAnimeResponseSinglePage retorna una respuesta con una sola página.
func MockAnimeResponseSinglePage() dto.AnimeResponse {
	return dto.AnimeResponse{
		Animes:     MockAnimeStructList()[:3],
		TotalPages: 1,
	}
}

// generateEpisodeNumbers genera una lista de números de episodios.
func generateEpisodeNumbers(start, end int) []int {
	episodes := make([]int, end-start+1)
	for i := range episodes {
		episodes[i] = start + i
	}
	return episodes
}

// MockAnimeByCategory retorna animes filtrados por categoría.
func MockAnimeByCategory(category dto.CategoryAnime) []dto.AnimeStruct {
	allAnimes := MockAnimeStructList()
	filtered := []dto.AnimeStruct{}

	for _, anime := range allAnimes {
		if anime.Type == category {
			filtered = append(filtered, anime)
		}
	}

	return filtered
}

// MockOVAAnime retorna un ejemplo de OVA.
func MockOVAAnime() dto.AnimeStruct {
	return dto.AnimeStruct{
		ID:          "steins-gate-ova",
		Title:       "Steins;Gate: Oukoubakko no Poriomania",
		Sinopsis:    "Un episodio especial que muestra la vida cotidiana.",
		Type:        dto.Ova,
		Punctuation: 8.4,
		Image:       "https://cdn.myanimelist.net/images/anime/12/35643.jpg",
	}
}

// MockMovieAnime retorna un ejemplo de Película.
func MockMovieAnime() dto.AnimeStruct {
	return dto.AnimeStruct{
		ID:          "demon-slayer-movie",
		Title:       "Kimetsu no Yaiba: Mugen Ressha-hen",
		Sinopsis:    "Tanjiro y sus amigos abordan el Tren Infinito.",
		Type:        dto.Pelicula,
		Punctuation: 8.8,
		Image:       "https://cdn.myanimelist.net/images/anime/1704/106947.jpg",
	}
}

// MockSpecialAnime retorna un ejemplo de Especial.
func MockSpecialAnime() dto.AnimeStruct {
	return dto.AnimeStruct{
		ID:          "one-piece-special-3d",
		Title:       "One Piece 3D: Mugiwara Chase",
		Sinopsis:    "Los Sombreros de Paja persiguen a un pájaro que robó el sombrero de Luffy.",
		Type:        dto.Especial,
		Punctuation: 7.5,
		Image:       "https://cdn.myanimelist.net/images/anime/8/26313.jpg",
	}
}
