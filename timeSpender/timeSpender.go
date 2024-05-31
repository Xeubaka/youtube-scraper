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
	dailyTimeSpender chan TimeSpended,
) {
	var day int
	var timeSpended TimeSpended
	timeSpendedDay := Day{Day: day}
	for video := range videoListWatcherChan {
		timeSpendedVideo := Video{
			ID:       video.ID,
			Duration: helper.ConvertStringToTime(video.ContentDetails.Duration),
		}
		timeSpended.TotalDuration += timeSpendedVideo.Duration
		if (timeSpendedVideo.Duration + timeSpendedDay.TotalTimeWatched) <= dailyTime[day%7] {
			timeSpended.TotalDurationWatched += timeSpendedVideo.Duration

			timeSpended, timeSpendedDay, timeSpendedVideo = watchVideo(
				timeSpended,
				timeSpendedDay,
				timeSpendedVideo,
			)
		} else if timeSpendedVideo.Duration <= dailyTime[(day+1)%7] {
			//Save current day on channel
			timeSpended.Days = append(timeSpended.Days, timeSpendedDay)

			//Reset day
			day, timeSpendedDay = newDay(day)

			//Insert video on the next day
			timeSpended, timeSpendedDay, timeSpendedVideo = watchVideo(
				timeSpended,
				timeSpendedDay,
				timeSpendedVideo,
			)
		} else {
			timeSpended.Days = append(timeSpended.Days, timeSpendedDay)
			day, timeSpendedDay = newDay(day)
		}
	}
	timeSpended.TotalDays = day - 2
	timeSpended.Days = append(timeSpended.Days[:0], timeSpended.Days[1:len(timeSpended.Days)-1]...)

	dailyTimeSpender <- timeSpended
	close(doneChannelTimeSpender)
}

func newDay(
	currentDay int,
) (
	int,
	Day,
) {
	return currentDay + 1, Day{Day: currentDay}
}

func watchVideo(
	timeSpended TimeSpended,
	timeSpendedDay Day,
	timeSpendedVideo Video,
) (
	TimeSpended,
	Day,
	Video,
) {
	timeSpendedDay.TotalTimeWatched += timeSpendedVideo.Duration
	timeSpended.TotalDurationWatched += timeSpendedVideo.Duration
	timeSpendedDay.Videos = append(timeSpendedDay.Videos, timeSpendedVideo)
	return timeSpended, timeSpendedDay, timeSpendedVideo
}
