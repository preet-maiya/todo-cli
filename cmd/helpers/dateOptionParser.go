package helpers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type DateOptionParsed struct {
	Type   string
	Sign   string
	Number int
}

func parseOption(dateOption string) DateOptionParsed {
	dateOptionRegex := "(?P<Sign>\\+|-)?(?P<Number>\\d+)(?P<Type>[dwmy])"
	var compRegEx = regexp.MustCompile(dateOptionRegex)
	match := compRegEx.FindStringSubmatch(dateOption)
	paramsMap := make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}

	n, err := strconv.Atoi(paramsMap["Number"])
	if err != nil {
		log.Errorf("Error converting number to int: %v", err)
		return DateOptionParsed{}
	}
	if paramsMap["Sign"] == "" {
		paramsMap["Sign"] = "+"
	}
	return DateOptionParsed{
		Type:   paramsMap["Type"],
		Sign:   paramsMap["Sign"],
		Number: n,
	}
}

func specialDateString(specialDate string) (time.Time, error) {
	now := time.Now()
	switch specialDate {
	case "today":
		return now, nil
	case "tomorrow":
		return now.AddDate(0, 0, 1), nil
	case "day-after", "day_after":
		return now.AddDate(0, 0, 1), nil
	case "yesterday":
		return now.AddDate(0, 0, -1), nil
	case "week-after", "week_after":
		return now.AddDate(0, 0, 7), nil
	case "week-before", "week_before":
		return now.AddDate(0, 0, -7), nil
	}
	return time.Time{}, errors.New("Unsupoorted special date format")
}

func nTimesDuration(durationString string) (time.Time, error) {
	now := time.Now()

	matched, err := regexp.MatchString("\\+|-?\\d+[dwmy]", durationString)
	if err != nil {
		log.Errorf("Error formatting the regex")
		return time.Time{}, err
	}
	if matched {
		parsedOption := parseOption(durationString)
		n := parsedOption.Number
		if parsedOption.Sign == "-" {
			n *= -1
		}
		switch parsedOption.Type {
		case "d":
			return now.AddDate(0, 0, n), nil
		case "w":
			return now.AddDate(0, 0, 7*n), nil
		case "m":
			return now.AddDate(0, n, 0), nil
		case "y":
			return now.AddDate(n, 0, 0), nil
		}
	}
	return time.Time{}, errors.New("Pattern not supported")
}

func dayOfWeek(day string) (time.Time, error) {
	today := time.Now().Weekday()
	days := 0
	switch day {
	case "Monday", "monday":
		days = int(time.Monday - today)
	case "Tuesday", "tuesday":
		days = int(time.Tuesday - today)
	case "Wednesday", "wednesday":
		days = int(time.Wednesday - today)
	case "Thursday", "thursday":
		days = int(time.Thursday - today)
	case "Friday", "friday":
		days = int(time.Friday - today)
	case "Saturday", "saturday":
		days = int(time.Saturday - today)
	case "Sunday", "sunday":
		days = int(time.Sunday - today)
	default:
		return time.Time{}, errors.New(fmt.Sprintf("Invalid weekday: %s", day))
	}

	if days <= 0 {
		days += 7
	}

	return time.Now().AddDate(0, 0, days), nil
}

func ParseDateOption(date string) string {
	if specialDate, err := specialDateString(date); err == nil {
		return specialDate.Format("2006-01-02")
	}

	if nTimesDuration, err := nTimesDuration(date); err == nil {
		return nTimesDuration.Format("2006-01-02")
	}

	if dayOfWeek, err := dayOfWeek(date); err == nil {
		return dayOfWeek.Format("2006-01-02")
	}

	if _, err := time.Parse("2006-01-02", date); err == nil {
		return date
	}

	// TODO: Return error
	return time.Time{}.Format("2006-01-02")
}
