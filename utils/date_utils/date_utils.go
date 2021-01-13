package date_utils

import "time"

const (
	apiDateLayout = "2006-01-02:02T15:04:05Z"
	apiDbLayout   = "2006-01-02 15:04:05"
)

// GetNow returns a Time type of the current time in UTC
func GetNow() time.Time {
	return time.Now().UTC()
}

// GetNowString returns a string type of the current time in UTC
// with a layouf of "2006-01-02:02T15:04:05Z"
func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}

// GetNowDBFormat returns a string type of the current time in UTC
// with a layouf of "2006-01-02 15:04:05"
func GetNowDBFormat() string {
	return GetNow().Format(apiDbLayout)
}
