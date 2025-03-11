package utils

import "time"

func FormatToVientianeTime(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Vientiane")
	return t.In(loc)
}
