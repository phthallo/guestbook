package internal

import (
	"fmt"
	"github.com/pcpratheesh/go-censorword"
)

func Filter(name string, message string) (string, string) {
	var detector = gocensorword.NewDetector(
		gocensorword.WithCensorReplaceChar("*"),
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