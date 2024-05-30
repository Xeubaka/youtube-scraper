package timespender

import (
	"time"
	helper "youtube-scraper/helper"
	youtubehandler "youtube-scraper/youtubeHandler"
)

func TimeSpendCounter(
	dailyTime []time.Duration,
	doneChannelTimeSpender chan bool,
	videoListWatcherChan chan youtubehandler.YoutubeVideo,
	dailyVideoWatchChan chan map[int][]map[string]time.Duration,
	totalDurationChan chan time.Duration,
) {
	day := 1
	var interator int
	var dayTimeCounter time.Duration
	var totalDurationTime time.Duration
	videoDayMap := make(map[int][]map[string]time.Duration)
	for video := range videoListWatcherChan {
		videoMap := make(map[string]time.Duration)
		videoDuration := helper.ConvertStringToTime(video.ContentDetails.Duration)
		totalDurationTime += videoDuration
		if (videoDuration + dayTimeCounter) <= dailyTime[interator] {
			dayTimeCounter += videoDuration
			videoMap[video.ID] = videoDuration
			videoDayMap[day] = append(videoDayMap[day], videoMap)
		} else {
			dayTimeCounter, _ = time.ParseDuration("0m0s")
			day++
			dailyVideoWatchChan <- videoDayMap
		}

		interator++
		if interator == 7 {
			interator = 0
		}
	}
	totalDurationChan <- totalDurationTime
	close(doneChannelTimeSpender)
}
