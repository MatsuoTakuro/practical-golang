package ch1

import (
	"fmt"
	"time"
)

func times() {
	timeTime()
	duration()
	println(nextMonthDay(time.Now()))
	println(nextMonthDay(time.Date(2021, 5, 31, 0, 0, 0, 0, time.Local)))
}

func timeTime() {
	now := time.Now()
	println(now, now.String())

	tz, _ := time.LoadLocation("America/Los_Angeles")
	println(tz)
	future := time.Date(2022, time.October, 21, 7, 28, 0, 0, tz)
	println(future, future.Format(time.RFC3339Nano))

	localPast := time.Date(1985, time.October, 26, 9, 0, 0, 0, time.Local)
	println(localPast)

	UTCPast := time.Date(1985, time.October, 26, 9, 0, 0, 0, time.UTC)
	println(UTCPast)
}

func duration() {
	fiveMinutes := 5 * time.Minute
	println(fiveMinutes)

	ten := 10
	tenSeconds := time.Duration(ten) * time.Second
	println(tenSeconds)

	past := time.Date(2022, time.May, 10, 0, 0, 0, 0, time.Local)
	dur := time.Since(past)
	println(dur)

	filepath := time.Now().Truncate(time.Hour).Format("20060102150405.json")
	println(filepath)

	fiveMinuteAfter := time.Now().Add(fiveMinutes)
	println(fiveMinuteAfter)

	fiveMinuteBefore := time.Now().Add(-fiveMinutes)
	println(fiveMinuteBefore)

	println("3 mins to stop started")
	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()
	<-timer.C
	println("3 mins stop completed.")

}

func nextMonthDay(t time.Time) time.Time {
	year1, month2, day := t.Date()
	firstDate := time.Date(year1, month2, 1, 0, 0, 0, 0, time.Local)
	year2, month2, _ := firstDate.AddDate(0, 1, 0).Date()                  // month2 is definitely the next month.
	nextMonthDate := time.Date(year2, month2, day, 0, 0, 0, 0, time.Local) // May not be the next month.
	if month2 != nextMonthDate.Month() {
		return firstDate.AddDate(0, 2, -1) // the end of the next month
	}
	return nextMonthDate
}

var num = 1

func println(any ...interface{}) {
	for _, a := range any {
		fmt.Printf("%d: %v\n", num, a)
		num++
	}
}
