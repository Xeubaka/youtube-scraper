package commandline

import (
	"fmt"
	"sort"
	"time"
	helper "youtube-scraper/helper"
	"youtube-scraper/wordfinder"
)

func GenerateTop5TextResponse(
	topWordsChan chan []wordfinder.WordCount,
	topTitleWordsChan chan []wordfinder.WordCount,
	topDescriptionWordsChan chan []wordfinder.WordCount,
	maxResults int,
) (text string) {
	text = mountParagraph(topWordsChan, TOP_5_RESPONSE_WORDS, maxResults)
	text += mountParagraph(topTitleWordsChan, TOP_5_RESPONSE_TITLE, maxResults)
	text += mountParagraph(topDescriptionWordsChan, TOP_5_RESPONSE_DESCRIPTION, maxResults)
	return
}

func GenerateTimeSpendedTextResponse(
	totalDurationChan chan time.Duration,
	dailyVideoWatchChan chan map[int][]map[string]time.Duration,
) (timeSpended string) {

	totalDuration := <-totalDurationChan
	timeSpended = mountBodyWithHeader(
		TOTAL_TIME_SPENDED_HEADER,
		fmt.Sprintf("| %v ", helper.ConvertTimeDurationToString(totalDuration)),
	)

	dailyVideoWatch := <-dailyVideoWatchChan
	orderedKeys, totalDays := orderMapforResponse(dailyVideoWatch)

	timeSpended += mountBodyWithHeader(
		TOTAL_DAYS_SPENDED_DEADER,
		fmt.Sprintf("| %v ", totalDays),
	)
	timeSpended += LINE_WITH_BRACKETS + "\n"

	for _, day := range orderedKeys {
		videos := dailyVideoWatch[day]

		header := fmt.Sprintf(DAY_OF_THE_WEEK_NUMBERED, time.Weekday(day%7), day)
		timeSpended += header
		for _, video := range videos {
			for key, value := range video {
				timeSpended += mountBody(
					header,
					fmt.Sprintf(
						"|id:%s|%s",
						key,
						helper.ConvertTimeDurationToString(value),
					),
				)
			}
		}
	}
	return
}

func mountParagraph(
	wordCountChan chan []wordfinder.WordCount,
	constTitle string,
	maxResults int,
) (text string) {
	wordCount := <-wordCountChan
	text = constTitle
	for k, word := range wordCount {
		text += mountBody(
			constTitle,
			fmt.Sprintf("|%d° '%s': %d time(s)", k+1, word.Word, word.Count),
		)

		if k+1 == maxResults {
			break
		}
	}
	text += LINE_WITH_BRACKETS + "\n"
	return
}

func mountBodyWithHeader(
	constHeader string,
	body string,
) (text string) {
	text = constHeader
	text += mountBody(
		constHeader,
		body,
	)
	return
}

func mountBody(
	constHeader string,
	body string,
) (text string) {
	text += body
	if x := len(body); x <= (len(constHeader) - 2) {
		text += helper.FillWithSpace(len(constHeader) - x - 3)
	}
	text += "|\n"
	return
}

func orderMapforResponse(
	dailyVideoWatch map[int][]map[string]time.Duration,
) (orderedKeys []int, totalDays int) {
	orderedKeys = make([]int, 0, len(dailyVideoWatch))
	for k := range dailyVideoWatch {
		orderedKeys = append(orderedKeys, k)
	}
	sort.Ints(orderedKeys)
	totalDays = orderedKeys[len(orderedKeys)-1]

	return
}
