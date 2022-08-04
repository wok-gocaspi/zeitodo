package service

import (
	"fmt"
	"math"
	"time"
)

func Startendtime() {

	start := time.Date(2022, time.August, 12, 23, 23, 0, 0, time.UTC)
	endtime := start.Add(8 * time.Hour)

	fmt.Println("the start time is  :", start, "\nlater:", endtime)

}

func CalcworkingDay(startTime time.Time, endTime time.Time) int {

	startOffset := weekday(startTime)
	startTime = startTime.AddDate(0, 0, -startOffset)
	endOffset := weekday(endTime)
	endTime = endTime.AddDate(0, 0, -endOffset)

	//Calculete weeks and days

	dif := endTime.Sub(startTime)
	weeks := int(math.Round((dif.Hours() / 24) / 7))
	days := -math.Min(float64(startOffset), 5) + math.Min(float64(endOffset), 5)
	{
		//Calcul total days

		return int(weeks)*5 + int(days)
	}
}
func weekday(d time.Time) int {
	wd := d.Weekday()
	if wd == time.Sunday {
		return 6
	}
	return int(wd) - 1
}
