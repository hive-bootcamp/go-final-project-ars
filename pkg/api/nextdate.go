package api

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	parsedTime, err := time.Parse(CommonDateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date: %w", err)
	}
	validSymbols := []string{"d", "y", "m", "w"}

	if repeat == "" {
		return "", errors.New("empty repeat")
	}
	if dstart == "" {
		return "", errors.New("empty date")
	}
	splittedRepeat := strings.Split(repeat, " ")
	repeatUnit := splittedRepeat[0]

	if !slices.Contains(validSymbols, repeatUnit) {
		return "", errors.New("invalid repeat unit")
	}

	switch repeatUnit {
	case "d":
		if len(splittedRepeat) != 2 {
			return "", errors.New("day interval missing")
		}
		repeatInterval, repeatErr := strconv.Atoi(splittedRepeat[1])
		if repeatErr != nil || repeatInterval > 400 {
			return "", errors.New("invalid day interval")
		}
		for {
			parsedTime = parsedTime.AddDate(0, 0, repeatInterval)
			if parsedTime.After(now) {
				break
			}
		}
	case "y":
		if len(splittedRepeat) != 1 {
			return "", errors.New("year repeat should not have interval")
		}
		for {
			parsedTime = parsedTime.AddDate(1, 0, 0)
			if parsedTime.After(now) {
				break
			}
		}
	case "w":
		if len(splittedRepeat) != 2 {
			return "", errors.New("week days missing")
		}
		daysStr := strings.Split(splittedRepeat[1], ",")
		var weekday [8]bool
		for _, d := range daysStr {
			n, err := strconv.Atoi(d)
			if err != nil || n < 1 || n > 7 {
				return "", errors.New("invalid weekday")
			}
			weekday[n] = true
		}
		for {
			parsedTime = parsedTime.AddDate(0, 0, 1)
			wd := int(parsedTime.Weekday())
			if wd == 0 {
				wd = 7
			}
			if weekday[wd] && parsedTime.After(now) {
				break
			}
		}
	case "m":
		if len(splittedRepeat) < 2 {
			return "", errors.New("month days missing")
		}
		dayStr := strings.Split(splittedRepeat[1], ",")
		var days []int
		for _, d := range dayStr {
			n, err := strconv.Atoi(d)
			if err != nil || n == 0 || n < -31 || n > 31 {
				return "", errors.New("invalid day of month")
			}
			days = append(days, n)
		}

		if len(splittedRepeat) > 2 {
			months := map[int]bool{}
			monthsStr := strings.Split(splittedRepeat[2], ",")
			for _, d := range monthsStr {
				n, err := strconv.Atoi(d)
				if err != nil || n < 1 || n > 12 {
					return "", errors.New("invalid month")
				}
				months[n] = true
			}

			for {
				parsedTime = parsedTime.AddDate(0, 0, 1)
				dayNum := parsedTime.Day()
				monthNum := int(parsedTime.Month())

				okDay := false
				for _, d := range days {
					if d > 0 && d == dayNum {
						okDay = true
						break
					}
					if d < 0 {
						lastDay := time.Date(parsedTime.Year(), parsedTime.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
						if -d == lastDay-dayNum+1 {
							okDay = true
							break
						}
					}
				}
				okMonth := len(months) == 0 || months[monthNum]

				if okDay && okMonth && parsedTime.After(now) {
					break
				}
			}
		}
	}

	return parsedTime.Format(CommonDateFormat), nil
}
