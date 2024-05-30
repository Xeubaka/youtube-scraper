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
	dailyVideoWatchChan := make(chan map[int][]map[string]time.Duration, maxResults)
	totalDurationChan := make(chan time.Duration, 1)
	go initializeTimeSpendCounter(
		dailyTime,
		doneChannelTimeSpender,
		channelVideoID,
		videoListWatcherChan,
		dailyVideoWatchChan,
		totalDurationChan,
	)

	//MenuOption
	handleOptionInteraction(
		option,
		doneChannelTop5,
		doneChannelTimeSpender,
		topWordsChan,
		topTitleWordsChan,
		topDescriptionWordsChan,
		totalDurationChan,
		dailyVideoWatchChan,
	)
}

func handleOptionInteraction(
	option int,
	doneChannelTop5 chan bool,
	doneChannelTimeSpender chan bool,
	topWordsChan chan []wordfinder.WordCount,
	topTitleWordsChan chan []wordfinder.WordCount,
	topDescriptionWordsChan chan []wordfinder.WordCount,
	totalDurationChan chan time.Duration,
	dailyVideoWatchChan chan map[int][]map[string]time.Duration,
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
					totalDurationChan,
					dailyVideoWatchChan,
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
	dailyVideoWatchChan chan map[int][]map[string]time.Duration,
	totalDurationChan chan time.Duration,
) {
	youtubehandler.FindAllVideos(maxResults, channelVideoID, videoListWatcherChan)
	timespender.TimeSpendCounter(
		dailyTime,
		doneChannelTimeSpender,
		videoListWatcherChan,
		dailyVideoWatchChan,
		totalDurationChan,
	)

}

func checkOnTop5Counter(
	doneChannelTop5 chan bool,
	topWordsChan chan []wordfinder.WordCount,
	topTitleWordsChan chan []wordfinder.WordCount,
	topDescriptionWordsChan chan []wordfinder.WordCount,
) {
	var text string
	for i := 0; ; i++ {
		commandline.CleanTerminal()
		select {
		case <-doneChannelTop5:
			top5 = commandline.GenerateTop5TextResponse(
				topWordsChan,
				topTitleWordsChan,
				topDescriptionWordsChan,
				maxResultsWords,
			)
			text = top5
		default:
			text = commandline.GetLoader(i)
		}
		commandline.PrintText(text)
		time.Sleep(1 * time.Second)
		if len(top5) > 0 || i > 10 {
			break
		}
	}
	commandline.PressEnterToContinue()
}

func checkOnTimeSpenderCounter(
	doneChannelTimeSpender chan bool,
	totalDurationChan chan time.Duration,
	dailyVideoWatchChan chan map[int][]map[string]time.Duration,
) {
	var text string
	for i := 0; ; i++ {
		commandline.CleanTerminal()
		select {
		case <-doneChannelTimeSpender:
			timeSpended = commandline.GenerateTimeSpendedTextResponse(
				totalDurationChan,
				dailyVideoWatchChan,
			)
			text = timeSpended
		default:
			text = commandline.GetLoader(i)
		}
		commandline.PrintText(text)
		time.Sleep(1 * time.Second)
		if len(timeSpended) > 0 || i > 10 {
			break
		}
	}
	commandline.PressEnterToContinue()
}
