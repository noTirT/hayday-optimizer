package base

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseDurationString(toParse string) (time.Duration, error) {
	re := regexp.MustCompile(`(?:(\d+) h)?\s*(?:(\d+) min)?`)

	matches := re.FindStringSubmatch(toParse)

	if matches == nil {
		return 0, ErrInvalidTimeString
	}

	var hours, minutes int
	var err error

	if matches[1] != "" {
		hours, err = strconv.Atoi(matches[1])
		if err != nil {
			return 0, ErrInvalidStringParse
		}
	}

	if matches[2] != "" {
		minutes, err = strconv.Atoi(matches[2])
		if err != nil {
			return 0, ErrInvalidStringParse
		}
	}

	totalDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
	return totalDuration, nil
}

func CapializeWordsOfString(text string) string {
	words := strings.Fields(text)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}
