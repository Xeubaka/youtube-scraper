package main

import (
	"fmt"
	broadcaster "youtube-scraper/Broadcaster"
	commandLine "youtube-scraper/commandLine"
	flags "youtube-scraper/flags"
	youtubehandler "youtube-scraper/youtubeHandler"

	"github.com/joho/godotenv"
)

var f flags.Flags

func loadEnv() (envMap map[string]string) {
	envMap, err := godotenv.Read()
	if err != nil {
		fmt.Println("Error reading .env file")
	}
	return
}

func init() {
	env := loadEnv()
	f = flags.GetFlags()

	//Flag value has preference over .env
	if len(f.ApiKey) == 0 {
		f.ApiKey = fmt.Sprintf(env["API_KEY"])
	}

	youtubehandler.SetApiKey(f.ApiKey)
	commandLine.Setup()
}

func main() {
	broadcaster.RunApplication(f.SearchQuery, f.DailyTime, f.Option, f.ResponseType)
}
