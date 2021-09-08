package utils

import (
	"log"
	"time"
)

var timeFormats = [...]string{
	"Mon Jan 02 2006 03:04:05",
	//Sat Feb 05 1977 00:00:00 GMT+0000 (UTC)"
	"2006-01-02T15:04:05.000Z",
	//2018-05-26T22:13:14.000Z
}

var Epoch = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

func TimeParse(timeString string) time.Time {
	for _, format := range timeFormats {
		sample := timeString
		if len(sample) > len(format) {
			sample = sample[0:len(format)]
		}
		tm, err := time.Parse(format, sample)
		if err == nil {
			return tm
		}
	}
	return time.Unix(0, 0)
}

func RoundStart(from time.Time, interval time.Duration) time.Time {
	newTime := from.Round(interval)
	if newTime.Before(from) {
		newTime = newTime.Add(interval)
	}
	return newTime
}

func WaitUntil(to time.Time) {
	if to.Before(time.Now()) {
		return
	}
	time.Sleep(to.Sub(time.Now()))
}

func WaitRound(interval time.Duration, message string) {
	to := RoundStart(time.Now(), interval)
	if len(message) > 0 {
		log.Printf("Waiting until %s to %s", to.Format(time.RFC3339), message)
	}
	WaitUntil(to)
}

func SameDay(date1,date2 time.Time) bool{
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}