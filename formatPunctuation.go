package main

import (
	"regexp"
)

func formatPunctuation(text string) string {

	reGroup := regexp.MustCompile(`([.,!?:;])[ \t]+([.,!?:;])`)
	for reGroup.MatchString(text) {
		text = reGroup.ReplaceAllString(text, "$1$2")
	}

	reBefore := regexp.MustCompile(`[ \t]+([.,!?:;])`)
	text = reBefore.ReplaceAllString(text, "$1")

	reAfter := regexp.MustCompile(`([.,!?:;]+)([\p{L}\p{N}])`)
	text = reAfter.ReplaceAllString(text, "$1 $2")

	return text
}
