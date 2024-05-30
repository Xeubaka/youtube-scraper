package youtubehandler

type YoutubeSearchResponse struct {
	NextPageToken string             `json:"nextPageToken"`
	PrevPageToken string             `json:"prevPageToken"`
	Items         []YoutubeVideoList `json:items`
}

type YoutubeVideoResponse struct {
	Items []YoutubeVideo `json:items`
}

type YoutubeVideoList struct {
	ID VideoID `json:"id"`

	Snippet Snippet `json:"snippet"`
}

type YoutubeVideo struct {
	ID             string         `json:"id"`
	ContentDetails ContentDetails `json:"contentDetails"`
}

type VideoID struct {
	VideoID string `json:"videoId"`
}

type Snippet struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ContentDetails struct {
	Duration string `json:"duration"`
}
