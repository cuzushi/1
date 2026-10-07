package main

import (
	"strings"
)

func quote(text string) string {
	var join string
	var app []string
	text1 := strings.Split(text, "\n")
	for j := 0; j < len(text1); j++ {

		words := strings.Fields(text1[j])
		inquote := false
		for i := 0; i < len(words); i++ {
			if words[i] == "'" {
				if !inquote {
					if i+1 < len(words) {
						words[i+1] = "'" + words[i+1]
						words = append(words[:i], words[i+1:]...)
						i--
					}
					inquote = true
				} else {
					if i > 0 {
						words[i-1] = words[i-1] + "'"
						words = append(words[:i], words[i+1:]...)
						i--
						inquote = false
					}
				}
			}

		}
		join = strings.Join(words, " ")
		app = append(app, join)
	}
	result := strings.Join(app, "\n")
	return result
}
