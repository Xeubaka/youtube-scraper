# Youtube-Scrapper CLI application

I made a simple scrapper onto the youtube Data Api
It has 2 main functions:
 - Guiven a searchQuery, it will return the top 5 words on titles and descriptions.
 - Guiven a SearchQuery and a daily quota of minutes in a week, it will calculate how much time you'll take to watch the first 200 videos.

# How to use 
> Its recommeded to use your own Google Api Key, instructions [here](https://developers.google.com/youtube/registering_an_application), but if you dont have don't worry i provide 3 functional keys on the ```.env.example```. 
-----
> This application use a ```.env``` file for setting the ```API_KEY```, but be aware because they have a [daily quota](https://developers.google.com/analytics/devguides/reporting/mcf/v3/limits-quotas) for requests.

## Downlaod executable

1. Search on the right side menu for the tag "Releases"
2. Download the file accordingly to yours OS
3. Execute it
> You can execute it from the terminal with args, run ``{EXECUTABLE_FILE} -h`` to see all available flags

## Download repository and run
1. Install Golang on you machine, instructions [here](https://go.dev/doc/install)

2. Them download this repository in your machine and run the commands:
    ```
    go mod tidy
    go build

3. It will create a executable file like: ```youtube-scraper.exe``` or just ```youtube-scraper``` if you are using linux 
4. (Optional) You can also execute with ```go run app.go``` (it is a bit slower since its gonna compile first, then run it) 


