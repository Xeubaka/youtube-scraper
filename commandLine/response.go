package commandline

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"
	helper "youtube-scraper/helper"
	timespender "youtube-scraper/timeSpender"
	"youtube-scraper/wordfinder"
)

func GenerateJsonTop5(
	topWordsChan chan []wordfinder.WordCount,
	topTitleWordsChan chan []wordfinder.WordCount,
	topDescriptionWordsChan chan []wordfinder.WordCount,
	maxResults int,
) (text []byte) {
	text = marshalWordCount(topWordsChan, maxResults)
	return
}

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

func GenerateJsonTimeSpended(
	dailyTimeSpender chan timespender.TimeSpended,
) (text []byte) {
	text = marshalTimeSpended(dailyTimeSpender)
	return
}

func GenerateTimeSpendedTextResponse(
	dailyTimeSpender chan timespender.TimeSpended,
) (text string) {
	timeSpended := <-dailyTimeSpender

	text = mountBodyWithHeader(
		TOTAL_TIME_SPENDED_HEADER,
		fmt.Sprintf("| %v ", timeSpended.TotalDuration),
	)

	text += mountBodyWithHeader(
		TOTAL_DAYS_SPENDED_DEADER,
		fmt.Sprintf("| %v ", timeSpended.TotalDays),
	)
	text += LINE_WITH_BRACKETS + "\n"

	for _, timeSpendedDay := range timeSpended.Days {
		videos := timeSpendedDay.Videos

		header := fmt.Sprintf(
			DAY_OF_THE_WEEK_NUMBERED,
			time.Weekday(timeSpendedDay.Day%7),
			timeSpendedDay.Day,
		)
		text += header

		text += mountBody(
			header,
			fmt.Sprintf("|Time Watched: %v", timeSpendedDay.TotalTimeWatched),
		)

		for _, timeSpendedVideo := range videos {
			text += mountBody(
				header,
				fmt.Sprintf(
					"|id:%s|%s",
					timeSpendedVideo.ID,
					timeSpendedVideo.Duration,
				),
			)
		}
	}
	return
}

func marshalWordCount(
	wordCountChan chan []wordfinder.WordCount,
 maxResults int,
) (text []byte) {
	wordCount := <-wordCountChan
	text, err := json.MarshalIndent(wordCount[:maxResults], "", "    ")
	fmt.Println(string(text))
	if err != nil {
		log.Fatal(err)
	}
	return
}

func marshalTimeSpended(
	dailyTimeSpender chan timespender.TimeSpended,
) (text []byte) {
	timeSpended := <-dailyTimeSpender
	text, err := json.MarshalIndent(timeSpended, "", "    ")
	fmt.Println(string(text))
	if err != nil {
		log.Fatal(err)
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
