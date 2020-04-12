package date

import "time"

var DEFAULT_FORMAT string = "2006-01-02"

func Now() time.Time {
	loc, _ := time.LoadLocation("Asia/Seoul")
	return time.Now().In(loc)
}
