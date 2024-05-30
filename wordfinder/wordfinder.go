package wordfinder

import (
	"fmt"
	"sort"
	"strings"
	helper "youtube-scraper/helper"
	youtubehandler "youtube-scraper/youtubeHandler"
)

type WordCount struct {
	word  string
	count int
}

func GenerateTop5(
	videosList chan youtubehandler.YoutubeVideoList,
	doneChannelTop5 chan bool,
	topFiveChan chan string,
	topFiveTitleChan chan string,
	topFiveDescriptionChan chan string,
) {
	wordMap, wordMapTitle, wordMapDescription := countWords(videosList)
	wordCounter(wordMap, topFiveChan)
	wordCounter(wordMapTitle, topFiveTitleChan)
	wordCounter(wordMapDescription, topFiveDescriptionChan)

	close(doneChannelTop5)
}

func wordCounter(wordMap map[string]int, topFiveChan chan string) {
	var wordCounter []WordCount

	wordCounter = addWordCounts(wordCounter, wordMap)
	sort.Slice(wordCounter, func(i, j int) bool {
		return wordCounter[i].count > wordCounter[j].count
	})

	go getTop5(wordCounter, topFiveChan)
}

func countWords(
	videosList chan youtubehandler.YoutubeVideoList,
) (map[string]int, map[string]int, map[string]int) {
	wordMap := make(map[string]int, 0)
	wordMapTitle := make(map[string]int, 0)
	wordMapDescription := make(map[string]int, 0)
	for video := range videosList {
		wordMap, wordMapTitle, wordMapDescription = getWordsInSnippet(wordMap, wordMapTitle, wordMapDescription, video.Snippet)
	}

	return wordMap, wordMapTitle, wordMapDescription
}

func getWordsInSnippet(
	wordMap map[string]int,
	wordMapTitle map[string]int,
	wordMapDescription map[string]int,
	snippet youtubehandler.Snippet,
) (map[string]int, map[string]int, map[string]int) {
	title := helper.TrimNonWords(snippet.Title)
	description := helper.TrimNonWords(snippet.Description)
	description = helper.TrimLinks(description)

	var allWords []string
	var titleWords []string
	var descriptionWords []string

	titleWords = strings.Split(title, " ")
	descriptionWords = strings.Split(description, " ")
	allWords = append(allWords, titleWords...)
	allWords = append(allWords, descriptionWords...)

	wordMap = fillWordMap(allWords, wordMap)
	wordMapTitle = fillWordMap(titleWords, wordMapTitle)
	wordMapDescription = fillWordMap(descriptionWords, wordMapDescription)

	return wordMap, wordMapTitle, wordMapDescription
}

func fillWordMap(
	words []string,
	wordMap map[string]int,
) map[string]int {
	for _, word := range words {
		if len(word) > 0 {
			wordMap[strings.ToLower(word)] += 1
		}
	}
	return wordMap
}

func addWordCounts(
	wordCounter []WordCount,
	wordMap map[string]int,
) []WordCount {
	for word, counter := range wordMap {
		wordCounter = append(wordCounter, WordCount{word: word, count: counter})
	}
	return wordCounter
}

func getTop5(
	wordCounter []WordCount,
	topFiveChan chan string,
) {
	for i := 0; i < len(wordCounter) && i < 5; i++ {
		result := fmt.Sprintf("%d° '%s': %d time(s)", i+1, wordCounter[i].word, wordCounter[i].count)
		if x := len(result); x <= 22 {
			result += helper.FillWithSpace(22 - x)
		}
		topFiveChan <- result
	}
	close(topFiveChan)
}
