package dto

// VideoURL representa una URL directa de video con su resolución.
type VideoURL struct {
	URL        string `json:"url"`
	Resolution string `json:"resolution"`
}
