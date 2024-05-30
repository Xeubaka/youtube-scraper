package flags

import (
	"flag"
	"log"
	"strings"
	"time"
	commandline "youtube-scraper/commandLine"
	"youtube-scraper/helper"
)

var f Flags

func loadFlags() {
	flag.StringVar(&f.SearchQuery, "searchQuery", "", "Query to search videos on youtube API")
	flag.StringVar(&f.DailyTimeFlag, "weekDailyTime", "", "Week Daily time, comma-separeted like:\n 15,120,30,150,20,40,90\n*In case of less than 7 inputs it will assume 0 minutes to that day")
	flag.IntVar(&f.Option, "option", 0, "Available options:\n  1 - Top 5 word\n  2 - Time spended")
	flag.StringVar(&f.ApiKey, "apiKey", "", "Send your own Google API KEY")

	flag.Parse()

	if len(f.DailyTimeFlag) > 0 {
		f.DailyTime = []time.Duration{0, 0, 0, 0, 0, 0, 0}
		str := helper.TrimWhiteSpaces(f.DailyTimeFlag)
		strMap := strings.Split(str, ",")
		for i := 0; i < 7; i++ {
			strInt, err := helper.ConvertStringToInt(strMap[i])
			if len(err) > 0 {
				commandline.PrintText(err)
				log.Fatal("Ending program...")
			}
			f.DailyTime[i] = time.Duration(strInt) * time.Minute
		}
	}
}

func GetFlags() Flags {
	loadFlags()
	return f
}
