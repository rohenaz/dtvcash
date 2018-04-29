package profile

const (
	TimeRange1Hour   = "1h"
	TimeRange24Hours = "24h"
	TimeRange7Days   = "7d"
	TimeRangeAll     = "all"
)

var AllTimeRanges = []string{
	TimeRange1Hour,
	TimeRange24Hours,
	TimeRange7Days,
	TimeRangeAll,
}

func StringIsTimeRange(str string) bool {
	for _, timeRangeConst := range AllTimeRanges {
		if str == timeRangeConst {
			return true
		}
	}
	return false
}
