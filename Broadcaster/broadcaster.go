package broadcaster

import (
	"time"
	commandline "youtube-scraper/commandLine"
	helper "youtube-scraper/helper"
	timespender "youtube-scraper/timeSpender"
	"youtube-scraper/wordfinder"
	youtubehandler "youtube-scraper/youtubeHandler"
)

var (
	top5            string
	timeSpended     string
	defaultOptions  = []int{1, 2, 9}
	maxResults      = 200
	maxResultsWords = 5
)

func RunApplication(
	searchQuery string,
	dailyTime []time.Duration,
	option int,
	responseType string,
) {
	//Search
	searchQuery = commandline.SearchMenu(searchQuery)
	videoListChan := make(chan youtubehandler.YoutubeVideoList)
	channelVideoID := make(chan string, maxResults)
	go requestSearch(
		searchQuery,
		videoListChan,
		channelVideoID,
	)

	//Top5
	doneChannelTop5 := make(chan bool)
	topWordsChan := make(chan []wordfinder.WordCount, 1)
	topTitleWordsChan := make(chan []wordfinder.WordCount, 1)
	topDescriptionWordsChan := make(chan []wordfinder.WordCount, 1)
	go initializeTop5Counter(
		videoListChan,
		doneChannelTop5,
		topWordsChan,
		topTitleWordsChan,
		topDescriptionWordsChan,
	)

	//DailyTimeSpent
	dailyTime = commandline.TimeSpendMenu(dailyTime)
	doneChannelTimeSpender := make(chan bool)
	videoListWatcherChan := make(chan youtubehandler.YoutubeVideo)
	dailyTimeSpender := make(chan timespender.TimeSpended, 1)
	go initializeTimeSpendCounter(
		dailyTime,
		doneChannelTimeSpender,
		channelVideoID,
		videoListWatcherChan,
		dailyTimeSpender,
	)

	//MenuOption
	handleOptionInteraction(
		option,
		doneChannelTop5,
		doneChannelTimeSpender,
		topWordsChan,
		topTitleWordsChan,
		topDescriptionWordsChan,
		dailyTimeSpender,
		responseType,
	)
}

func handleOptionInteraction(
	option int,
	doneChannelTop5 chan bool,
	doneChannelTimeSpender chan bool,
	topWordsChan chan []wordfinder.WordCount,
	topTitleWordsChan chan []wordfinder.WordCount,
	topDescriptionWordsChan chan []wordfinder.WordCount,
	dailyTimeSpender chan timespender.TimeSpended,
	responseType string,
) {
	for option != 9 {
		option = commandline.OptionMenu(option, defaultOptions)

		if option == 1 {
			if len(top5) > 0 {
				commandline.CleanTerminal()
				commandline.PrintText(top5)
				commandline.PressEnterToContinue()
			} else {
				checkOnTop5Counter(
					doneChannelTop5,
					topWordsChan,
					topTitleWordsChan,
					topDescriptionWordsChan,
					responseType,
				)
			}
			option = 0
		}

		if option == 2 {
			if len(timeSpended) > 0 {
				commandline.CleanTerminal()
				commandline.PrintText(timeSpended)
				commandline.PressEnterToContinue()
			} else {
				checkOnTimeSpenderCounter(
					doneChannelTimeSpender,
					dailyTimeSpender,
					responseType,
				)
			}
			option = 0
		}
	}
	commandline.EndProgram()
}

func requestSearch(
	searchQuery string,
	videoListChan chan youtubehandler.YoutubeVideoList,
	channelVideoID chan string,
) {
	youtubehandler.SearchAll(
		helper.TrimWhiteSpaces(searchQuery),
		maxResults,
		videoListChan,
		channelVideoID,
	)
}

func initializeTop5Counter(
	videoListChan chan youtubehandler.YoutubeVideoList,
	doneChannelTop5 chan bool,
	topWordsChan chan []wordfinder.WordCount,
	topTitleWordsChan chan []wordfinder.WordCount,
	topDescriptionWordsChan chan []wordfinder.WordCount,
) {
	wordfinder.GenerateTop5(
		videoListChan,
		doneChannelTop5,
		topWordsChan,
		topTitleWordsChan,
		topDescriptionWordsChan,
	)
}

func initializeTimeSpendCounter(
	dailyTime []time.Duration,
	doneChannelTimeSpender chan bool,
	channelVideoID chan string,
	videoListWatcherChan chan youtubehandler.YoutubeVideo,
	dailyTimeSpender chan timespender.TimeSpended,
) {
	youtubehandler.FindAllVideos(maxResults, channelVideoID, videoListWatcherChan)
	timespender.TimeSpendCounter(
		dailyTime,
		doneChannelTimeSpender,
		videoListWatcherChan,
		dailyTimeSpender,
	)

}

func checkOnTop5Counter(
	doneChannelTop5 chan bool,
	topWordsChan chan []wordfinder.WordCount,
	topTitleWordsChan chan []wordfinder.WordCount,
	topDescriptionWordsChan chan []wordfinder.WordCount,
	responseType string,
) {
	var text string
	for i := 0; ; i++ {
		commandline.CleanTerminal()
		select {
		case <-doneChannelTop5:
			if responseType == "json" {
				top5 = string(commandline.GenerateJsonTop5(
					topWordsChan,
					topTitleWordsChan,
					topDescriptionWordsChan,
					maxResultsWords,
				))
			} else {
				top5 = commandline.GenerateTop5TextResponse(
					topWordsChan,
					topTitleWordsChan,
					topDescriptionWordsChan,
					maxResultsWords,
				)
			}
			text = top5
		default:
			text = commandline.GetLoader(i)
		}
		if responseType != "json" {
			commandline.PrintText(text)
		}
		time.Sleep(1 * time.Second)
		if len(top5) > 0 || i > 10 {
			break
		}
	}
	commandline.PressEnterToContinue()
}

func checkOnTimeSpenderCounter(
	doneChannelTimeSpender chan bool,
	dailyTimeSpender chan timespender.TimeSpended,
	responseType string,
) {
	var text string
	for i := 0; ; i++ {
		commandline.CleanTerminal()
		select {
		case <-doneChannelTimeSpender:
			if responseType == "json" {
				timeSpended = string(commandline.GenerateJsonTimeSpended(
					dailyTimeSpender,
				))
			} else {
				timeSpended = commandline.GenerateTimeSpendedTextResponse(
					dailyTimeSpender,
				)
			}
			text = timeSpended
		default:
			text = commandline.GetLoader(i)
		}
		if responseType != "json" {
			commandline.PrintText(text)
		}
		time.Sleep(1 * time.Second)
		if len(timeSpended) > 0 || i > 10 {
			break
		}
	}
	commandline.PressEnterToContinue()
}
