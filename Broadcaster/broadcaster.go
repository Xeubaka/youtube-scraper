package broadcaster

import (
	"fmt"
	"sort"
	"time"
	commandline "youtube-scraper/commandLine"
	helper "youtube-scraper/helper"
	timespender "youtube-scraper/timeSpender"
	"youtube-scraper/wordfinder"
	youtubehandler "youtube-scraper/youtubeHandler"
)

var (
	top5           string
	timeSpended    string
	dailyTime      []time.Duration
	defaultOptions = []int{1, 2, 9}
	maxResults     = 200
)

func RunApplication(searchQuery string, dailyTime []time.Duration, option int) {
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
	topFiveChan := make(chan string, 5)
	topFiveTitleChan := make(chan string, 5)
	topFiveDescriptionChan := make(chan string, 5)
	go initializeTop5Counter(
		videoListChan,
		doneChannelTop5,
		topFiveChan,
		topFiveTitleChan,
		topFiveDescriptionChan,
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
		topFiveChan,
		topFiveTitleChan,
		topFiveDescriptionChan,
		totalDurationChan,
		dailyVideoWatchChan,
	)
}

func handleOptionInteraction(
	option int,
	doneChannelTop5 chan bool,
	doneChannelTimeSpender chan bool,
	topFiveChan chan string,
	topFiveTitleChan chan string,
	topFiveDescriptionChan chan string,
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
				checkOnTop5Counter(doneChannelTop5, topFiveChan, topFiveTitleChan, topFiveDescriptionChan)
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
	topFiveChan chan string,
	topFiveTitleChan chan string,
	topFiveDescriptionChan chan string,
) {
	wordfinder.GenerateTop5(videoListChan, doneChannelTop5, topFiveChan, topFiveTitleChan, topFiveDescriptionChan)
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
	topFiveChan chan string,
	topFiveTitleChan chan string,
	topFiveDescriptionChan chan string,
) {
	var text string
	for i := 0; ; i++ {
		commandline.CleanTerminal()
		select {
		case <-doneChannelTop5:
			top5 = "| Top 5 words founded: |\n"
			for word := range topFiveChan {
				top5 += fmt.Sprintf("|%s|\n", word)
			}
			top5 += "|----------------------|\n"
			top5 += "| Top words in Title:  |\n"
			for word := range topFiveTitleChan {
				top5 += fmt.Sprintf("|%s|\n", word)
			}
			top5 += "|----------------------|\n"
			top5 += "|Top words Description:|\n"
			for word := range topFiveDescriptionChan {
				top5 += fmt.Sprintf("|%s|\n", word)
			}
			top5 += "|----------------------|\n"
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
			header := "| Total time spended:  |\n"
			timeSpended = header
			totalDuration := <-totalDurationChan
			body := fmt.Sprintf("| %v ", helper.ConvertTimeDurationToString(totalDuration))
			timeSpended += body
			if x := len(body); x <= (len(header) - 2) {
				timeSpended += helper.FillWithSpace(len(header) - x - 3)
			}
			timeSpended += "|\n"
			dailyVideoWatch := <-dailyVideoWatchChan

			keys := make([]int, 0, len(dailyVideoWatch))

			for k := range dailyVideoWatch {
				keys = append(keys, k)
			}

			sort.Ints(keys)
			totalDays := keys[len(keys)-1]
			header = "| Total days spended:  |\n"
			body = fmt.Sprintf("| %v ", totalDays)
			timeSpended += body
			if x := len(body); x <= (len(header) - 2) {
				timeSpended += helper.FillWithSpace(len(header) - x - 3)
			}
			timeSpended += "|\n"
			timeSpended += "|----------------------|\n"
			for _, day := range keys {
				videos := dailyVideoWatch[day]

				header := fmt.Sprintf("|------%s(%d)------|\n", time.Weekday(day%7), day)
				timeSpended += header
				var body string
				for _, video := range videos {
					for key, value := range video {
						body = fmt.Sprintf("|id:%s|%s", key, helper.ConvertTimeDurationToString(value))
						timeSpended += body
						if x := len(body); x <= (len(header) - 2) {
							timeSpended += helper.FillWithSpace(len(header) - x - 3)
						}
						timeSpended += "|\n"
					}
				}
			}
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
