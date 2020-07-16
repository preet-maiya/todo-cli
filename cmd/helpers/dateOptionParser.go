package helpers

import (
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

func ParseDateOption(date string) string {
	now := time.Now()
	switch date {
	case "today":
		return now.Format("2006-01-02")
	case "tomorrow":
		return now.AddDate(0, 0, 1).Format("2006-01-02")
	case "day-after", "day_after":
		return now.AddDate(0, 0, 1).Format("2006-01-02")
	case "yesterday":
		return now.AddDate(0, 0, -1).Format("2006-01-02")
	case "week-after", "week_after":
		return now.AddDate(0, 0, 7).Format("2006-01-02")
	case "week-before", "week_before":
		return now.AddDate(0, 0, -7).Format("2006-01-02")
	}

	matched, err := regexp.MatchString("\\+|-?\\d+[dwmy]", date)
	if err != nil {
		log.Errorf("Error formatting the regex")
		return time.Time{}.Format("2006-01-02")
	}
	if matched {
		parsedOption := parseOption(date)
		n := parsedOption.Number
		if parsedOption.Sign == "-" {
			n *= -1
		}
		switch parsedOption.Type {
		case "d":
			return now.AddDate(0, 0, n).Format("2006-01-02")
		case "w":
			return now.AddDate(0, 0, 7*n).Format("2006-01-02")
		case "m":
			return now.AddDate(0, n, 0).Format("2006-01-02")
		case "y":
			return now.AddDate(n, 0, 0).Format("2006-01-02")
		}
	}

	return date
}
