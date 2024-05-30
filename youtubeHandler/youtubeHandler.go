package youtubehandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	apiKey   string
	myClient = &http.Client{Timeout: 10 * time.Second}
)

const baseURL = "https://www.googleapis.com/youtube/v3/"

func SetApiKey(key string) {
	apiKey = key
}

func SearchAll(
	searchQuery string,
	total int,
	videoListChan chan YoutubeVideoList,
	channelVideoID chan string,
) {
	total /= 50
	pageToken := ""
	var videoList []YoutubeVideoList
	for i := 0; i < total; i++ {
		var newVideoList []YoutubeVideoList
		pageToken, newVideoList = SearchByPage(searchQuery, pageToken)
		videoList = append(videoList, newVideoList...)
	}
	go func(
		videoListChan chan YoutubeVideoList,
		videoList []YoutubeVideoList,
	) {
		for _, video := range videoList {
			videoListChan <- video
		}
		close(videoListChan)
	}(videoListChan, videoList)
	go func(
		channelVideoID chan string,
		videoList []YoutubeVideoList,
	) {
		for _, video := range videoList {
			channelVideoID <- video.ID.VideoID
		}
		close(channelVideoID)
	}(channelVideoID, videoList)
}

func SearchByPage(
	query string,
	pageToken string,
) (nextPageToken string, videoList []YoutubeVideoList) {
	params := map[string]string{
		"part":       "snippet",
		"maxResults": "50",
		"type":       "video",
		"q":          url.QueryEscape(query),
		"key":        apiKey,
		"pageToken":  pageToken,
	}

	r := execRequest("search", params)
	defer r.Body.Close()

	var youtubeSearchResponse YoutubeSearchResponse
	if err := json.NewDecoder(r.Body).Decode(&youtubeSearchResponse); err != nil {
		log.Fatal(err)
	}

	nextPageToken = youtubeSearchResponse.NextPageToken
	videoList = youtubeSearchResponse.Items
	return
}

func FindAllVideos(
	total int,
	channelVideoID chan string,
	videoListWatcherChan chan YoutubeVideo,
) {
	total /= 50
	var videoList []YoutubeVideo
	for i := 0; i < total; i++ {
		newVideoList := FindVideo(channelVideoID)
		videoList = append(videoList, newVideoList...)
	}

	go func(videoListWatcherChan chan YoutubeVideo, videoList []YoutubeVideo) {
		for _, video := range videoList {
			videoListWatcherChan <- video
		}
		close(videoListWatcherChan)
	}(videoListWatcherChan, videoList)
}

func FindVideo(channelVideoID chan string) []YoutubeVideo {
	params := map[string]string{
		"part": "contentDetails",
		"key":  apiKey,
	}

	i := 0
	for videoID := range channelVideoID {
		params["id"] += fmt.Sprintf("%s,", videoID)
		i++
		if i == 50 {
			break
		}
	}
	params["id"] = strings.TrimSuffix(params["id"], ",")

	r := execRequest("videos", params)
	defer r.Body.Close()

	var youtubeVideoResponse YoutubeVideoResponse
	if err := json.NewDecoder(r.Body).Decode(&youtubeVideoResponse); err != nil {
		log.Fatal(err)
	}

	return youtubeVideoResponse.Items
}

func createURL(
	method string,
	params map[string]string,
) (url string) {
	url = fmt.Sprintf("%s%s?", baseURL, method)

	for key, value := range params {
		if len(value) > 0 {
			url += fmt.Sprintf("%s=%s&", key, value)
		}
	}

	return
}

func execRequest(
	method string,
	params map[string]string,
) *http.Response {
	urlSearch := createURL(method, params)

	r, err := myClient.Get(urlSearch)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode > 300 {
		log.Fatal(r.Status)
	}
	return r
}
