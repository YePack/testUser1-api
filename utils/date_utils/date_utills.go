package date_utils

import "time"

const (
	apiDateLayout = "2006-01-02T15:04:05Z"
	apiDbLayout   = "2006-01-02 5:04:05"
)

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetNowDBFormat() string {
	return GetNow().Format(apiDbLayout)
}
