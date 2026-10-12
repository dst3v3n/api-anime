package dto

type AnimeSummary struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
}

type AnimeStruct struct {
	AnimeSummary
	Synopsis    *string  `json:"synopsis,omitempty"`
	Type        *string  `json:"type,omitempty"`
	Punctuation *float64 `json:"punctuation,omitempty"`
}

type Anime struct {
	Animes     []AnimeStruct `json:"animes"`
	TotalPages uint          `json:"total_pages"`
	Source     string        `json:"source"`
}

type AnimeInfoSummary struct {
	Genres       []string `json:"genres"`
	Status       string   `json:"status"`
	FirstEpisode int      `json:"first_episode"`
	LastEpisode  int      `json:"last_episode"`
}

type AnimeInfo struct {
	AnimeStruct
	AnimeInfoSummary
	AnimeRelated *[]AnimeRelated `json:"related,omitempty"`
	NextEpisode  *string         `json:"next_episode,omitempty"`
	Source       string          `json:"source"`
}

type AnimeRelated struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
}
