package actions

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

func repeatByDays(now time.Time, date *time.Time, repeat string) error {

	repeatRule := strings.Split(repeat, " ")

	if len(repeatRule) != 2 {
		return errors.New("не указан интервал в днях")
	}

	daysCount, err := strconv.Atoi(repeatRule[1])

	if daysCount > 400 {
		return errors.New("превышен максимально допустимый интервал")
	}

	if err != nil {
		return err
	}

	for {
		*date = date.AddDate(0, 0, daysCount)

		if date.After(now) {
			break
		}
	}

	return nil
}

func repeatByWeek(now time.Time, date *time.Time, repeat string) error {

	repeatRule := strings.Split(repeat, " ")

	if len(repeatRule) != 2 {
		return errors.New("не указан интервал в днях")
	}

	weekDaysString := strings.Split(repeatRule[1], ",")

	dates := make([]time.Time, 0, len(weekDaysString))

	for _, dayString := range weekDaysString {

		day, err := strconv.Atoi(dayString)

		if err != nil {
			return err
		}
		if day < 1 || day > 7 {
			return errors.New("некорректный день недели")
		}

		dates = append(dates, nextDayByWeekNumber(now, *date, day))
	}

	sortSliceDates(dates, true)
	*date = dates[0]

	return nil
}

func repeatByMonthDay(now time.Time, date *time.Time, repeat string) error {

	repeatRule := strings.Split(repeat, " ")

	if len(repeatRule) < 2 {
		return errors.New("не указан интервал в днях")
	}
	days := []int{}
	months := []int{}

	dayRules := strings.Split(repeatRule[1], ",")

	if len(dayRules) == 0 {
		return errors.New("не указан интервал в днях")
	}

	for _, v := range dayRules {
		dayNum, err := strconv.Atoi(v)

		if err != nil {
			return err
		}

		days = append(days, dayNum)
	}

	if len(repeatRule) == 3 {
		monthsRules := strings.Split(repeatRule[2], ",")

		for _, v := range monthsRules {
			monthNum, err := strconv.Atoi(v)

			if err != nil {
				return err
			}

			months = append(months, monthNum)
		}
	}

	dates := make([]time.Time, 0, len(days)*len(months))

	for _, day := range days {

		if len(months) > 0 {
			for _, month := range months {
				dates = append(dates, nextDayByMonthNumber(now, *date, day, month))
			}
			continue
		}

		dates = append(dates, nextDayByMonthNumber(now, *date, day, 0))
	}

	sortSliceDates(dates, true)
	*date = dates[0]

	return nil
}

func nextDayByMonthNumber(now, date time.Time, dayNumber, monthNumber int) time.Time {

	for {
		date = date.AddDate(0, 0, 1)

		if date.Before(now) {
			continue
		}
		fixedDayNumber := dayNumber
		if dayNumber < 0 {
			fixedDayNumber = dayInMonth(date) + dayNumber
		}

		if date.Day() != fixedDayNumber {
			continue
		}

		if monthNumber != 0 && int(date.Month()) != monthNumber {
			continue
		}

		return date
	}
}

func nextDayByWeekNumber(now time.Time, date time.Time, weekDay int) time.Time {
	startDate := date

	for int(date.Weekday()) != weekDay%7 || now.After(date) || startDate.Equal(date) {
		date = date.AddDate(0, 0, 1)
	}

	return date
}

func sortSliceDates(arr []time.Time, ascending bool) {

	sort.Slice(arr, func(i, j int) bool {
		if arr[i].After(arr[j]) {
			return !ascending
		}

		return ascending
	})
}

func dayInMonth(date time.Time) int {
	return time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
}
