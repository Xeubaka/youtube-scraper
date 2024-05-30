package commandline

import (
	"fmt"
	"slices"
	"time"
	helper "youtube-scraper/helper"
)

var warning string

func clearWarning() {
	warning = ""
}

func SearchMenu(
	searchQuery string,
) string {
	for len(searchQuery) == 0 {
		if len(warning) > 0 {
			CleanTerminal()
			PrintText(warning)
		}
		fmt.Println(SEARCH_QUESTION)

		searchQuery = ReadInput()

		searchQuery = helper.TrimEndLine(searchQuery)

		if len(searchQuery) == 0 {
			warning = fmt.Sprint(TERM_CANNOT_BE_EMPTY, searchQuery)
		}
	}
	clearWarning()
	return searchQuery
}

func OptionMenu(
	optionIn int,
	availableOptions []int,
) (option int) {
	if optionIn != 0 {
		return optionIn
	}

	var isOptionInvalid bool = true
	for isOptionInvalid {
		CleanTerminal()
		if len(warning) > 0 {
			PrintText(warning)
		}

		fmt.Println("Choose an option:")
		fmt.Println("1 - Top 5 founded words")
		fmt.Println("2 - Calculate time to watch all videos") //first 200 videos founded
		// fmt.Println("3 - Search again") //maybe?
		// fmt.Println("4 - Calculate for a different time") //maybe?
		fmt.Println("9 - Exit")

		text := ReadInput()

		text = helper.TrimWhiteSpaces(text)

		option, warning = helper.ConvertStringToInt(text)

		isOptionInvalid = !slices.Contains(availableOptions, option)
		if isOptionInvalid {
			warning = fmt.Sprintf(INVALID_OPTION_SELECTED, option)
		}
	}
	clearWarning()
	return
}

func TimeSpendMenu(
	dailyTime []time.Duration,
) []time.Duration {
	if len(dailyTime) > 0 {
		return dailyTime
	}

	var num int
	var err error
	var week = []time.Duration{0, 0, 0, 0, 0, 0, 0}

	for day := range week {
		for week[day] <= 0*time.Minute {
			CleanTerminal()
			if len(warning) > 0 {
				PrintText(warning)
			}

			fmt.Printf(TIME_SPENDER_QUESTION, time.Weekday(day))
			text := ReadInput()
			num, warning = helper.ConvertStringToInt(text)
			week[day], err = time.ParseDuration(fmt.Sprintf("%dm", num))

			if err != nil {
				panic(err)
			}
		}
	}
	clearWarning()
	return week
}
