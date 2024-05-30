package flags

import "time"

type Flags struct {
	SearchQuery   string
	DailyTimeFlag string
	DailyTime     []time.Duration
	Option        int
	ApiKey        string
}
