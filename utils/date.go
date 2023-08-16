package utils

import "time"

func GetCurrentDayOfWeek() int {
	weekday := time.Now().Weekday()
	intWeekday := int(weekday)
	if intWeekday == 0 {
		intWeekday = 7
	}
	return intWeekday
}
