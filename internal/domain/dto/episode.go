package dto

type RecentEpisode struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Episode uint   `json:"episode"`
	Image   string `json:"image"`
	Source  string `json:"source"`
}

type EpisodeStreamingProvider struct {
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	Episode  uint       `json:"episode"`
	Provider []Provider `json:"providers"`
	Source   string     `json:"source"`
}

type Provider struct {
	Server string `json:"server"`
	URL    string `json:"url"`
}
