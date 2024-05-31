package timespender

import "time"

type TimeSpended struct {
	TotalDays            int           `json:totalDays`
	TotalDuration        time.Duration `json:totalDuration`
	TotalDurationWatched time.Duration `json:totalDurationWatched`
	Days                 []Day         `json:days`
}

type Day struct {
	Day              int           `json:day`
	TotalTimeWatched time.Duration `json:totalTimeWatched`
	Videos           []Video       `json:videos`
}

type Video struct {
	ID       string        `json:id`
	Duration time.Duration `json:duration`
}
