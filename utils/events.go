package utils

import (
	"math"
	"time"
)

func WeekNum(t time.Time) int {
	firstDay := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)

	adjustedDom := t.Day() + int(firstDay.Weekday())
	return int(math.Ceil(float64(adjustedDom) / 7.0))
}

func CheckEvent(eTimeStart time.Time, eTimeEnd time.Time,
	rTimeStart time.Time, rTimeEnd time.Time,
	repeatType string) bool {

	if repeatType == "" {
		return rTimeStart.Before(eTimeEnd) && rTimeEnd.After(eTimeStart)
	}

	eStart := time.Date(0, 1, 1,
		eTimeStart.Hour(), eTimeStart.Minute(), eTimeStart.Second(), 0, time.UTC)
	eEnd := time.Date(0, 1, eTimeEnd.Day()-eTimeStart.Day()+1,
		eTimeEnd.Hour(), eTimeEnd.Minute(), eTimeEnd.Second(), 0, time.UTC)
	rStart := time.Date(0, 1, 1,
		rTimeStart.Hour(), rTimeStart.Minute(), rTimeStart.Second(), 0, time.UTC)
	rEnd := time.Date(0, 1, rTimeEnd.Day()-rTimeStart.Day()+1,
		rTimeEnd.Hour(), rTimeEnd.Minute(), rTimeEnd.Second(), 0, time.UTC)

	if repeatType == "day" {
		if !(eStart.Before(rEnd) && eEnd.After(rStart)) {
			return false
		}
		return true
	}

	if repeatType == "workday" {
	}

	if repeatType == "week" {
		eStart = eStart.AddDate(0, 0, int(eTimeStart.Weekday()))
		eEnd = eEnd.AddDate(0, 0, int(eTimeEnd.Weekday()))
		rStart = rStart.AddDate(0, 0, int(rTimeStart.Weekday()))
		rEnd = rEnd.AddDate(0, 0, int(rTimeEnd.Weekday()))
		if !(eStart.Before(rEnd) && eEnd.After(rStart)) {
			return false
		}
		return true
	}
	if repeatType == "month" {
		eStart = eStart.AddDate(0, 0, WeekNum(eTimeStart)*7)
		eEnd = eEnd.AddDate(0, 0, WeekNum(eTimeEnd)*7)
		rStart = rStart.AddDate(0, 0, WeekNum(rTimeStart)*7)
		rEnd = rEnd.AddDate(0, 0, WeekNum(rTimeEnd)*7)
		if !(eStart.Before(rEnd) && eEnd.After(rStart)) {
			return false
		}
		return true
	}
	if repeatType == "year" {
		eStart = eStart.AddDate(0, int(eTimeStart.Month()), 0)
		eEnd = eEnd.AddDate(0, int(eTimeEnd.Month()), 0)
		rStart = rStart.AddDate(0, int(rTimeStart.Month()), 0)
		rEnd = rEnd.AddDate(0, int(rTimeStart.Month()), 0)
		if !(eStart.Before(rEnd) && eEnd.After(rStart)) {
			return false
		}
		return true
	}
	return false
}
