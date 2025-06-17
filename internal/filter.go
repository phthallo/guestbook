package internal

import (
	"fmt"
	"github.com/phthallo/go-censorword"

)

var CensorWhiteList = []string { // this is not an extensive whitelist lol but just the most common miscensors i've seen
	"ho",
	"hell",
	"wtf",
	"sob",
	"mad",
	"gin",
}

func Filter(name string, message string) (string, string) {
	var detector = gocensorword.NewDetector(
		gocensorword.WithCensorReplaceChar("*"),
		gocensorword.WithSanitizeSpecialCharacters(false),
		gocensorword.WithCustomWhiteList(CensorWhiteList),
		gocensorword.WithReplaceCheckPattern("(?i) %s"),
	)

	filteredName, err := detector.CensorWord(name)
	if (err != nil) {
		fmt.Errorf("Error filtering name: %v", err)
	}

	filteredMessage, err := detector.CensorWord(message)
	if (err != nil) {
		fmt.Errorf("Error filtering message: %v", err)
	}

	return filteredName, filteredMessage
}