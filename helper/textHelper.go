package helper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func TrimWhiteSpaces(
	text string,
) string {
	text = strings.TrimSpace(text)
	regexpWhiteSpaces := regexp.MustCompile(`[\s\p{Zs}]{2,}`)

	return regexpWhiteSpaces.ReplaceAllString(text, " ")
}

func TrimNonWords(
	text string,
) string {
	regexpNonWords := regexp.MustCompile(`[^A-Za-zÀ-Üà-ü]`)
	return TrimWhiteSpaces(regexpNonWords.ReplaceAllString(text, " "))
}

func TrimLinks(
	text string,
) string {
	regexpLink := regexp.MustCompile(`(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w\.-]*)`)
	return TrimWhiteSpaces(regexpLink.ReplaceAllString(text, " "))
}

func TrimEndLine(
	text string,
) string {
	regexpEndLine := regexp.MustCompile(`\r|\n`)
	return regexpEndLine.ReplaceAllString(text, "")
}

func ConvertStringToInt(
	text string,
) (number int, printErr string) {
	number, err := strconv.Atoi(TrimEndLine(text))
	if err != nil {
		printErr = fmt.Sprintf("| %s is not an number. |\n", text)
	}
	return
}

func ConvertStringToTime(
	text string,
) (duration time.Duration) {
	regexpTime := regexp.MustCompile(`P\d{1,}DT|[PT]`)
	regexpString := regexpTime.ReplaceAllString(text, "")
	lowerString := strings.ToLower(regexpString)
	// https://en.wikipedia.org/wiki/ISO_8601#Durations: "PT0S" or "P0D" represents 0s
	if lowerString == "0d" {
		lowerString = "0s"
	}
	duration, err := time.ParseDuration(strings.ToLower(lowerString))
	if err != nil {
		panic(err)
	}

	return
}

func ConvertTimeDurationToString(
	duration time.Duration,
) string {
	s := duration.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "d0h") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "M0d") {
		s = s[:len(s)-2]
	}
	return s
}

func FillWithSpace(size int) (spaceString string) {
	for i := 0; i <= size; i++ {
		spaceString += " "
	}
	return
}
