package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func capitalize(s string) string {
	r := []rune(s)
	if len(r) == 0 {
		return s
	}
	return strings.ToUpper(string(r[0])) + strings.ToLower(string(r[1:]))
}
func nextWordStart(content []string, i int) (byte, bool) {
	for j := i + 1; j < len(content); j++ {
		w := strings.TrimLeft(content[j], `'"([{`)
		if w != "" {
			return w[0], true
		}
	}
	return 0, false
}
func isvowel(v byte) bool {
	word := strings.ToLower(string(v))
	return word == "a" || word == "e" || word == "i" || word == "u" || word == "o" || word == "h"
}

func ProcessText(input string) string {
	text := strings.Split(input, "\n")
	var slice []string
	var result string

	for j := 0; j < len(text); j++ {
		if len(text[j]) > 0 {
			content := strings.Fields(text[j])

			for i := 0; i < len(content); i++ {
				if content[i] == "(hex)" {
					if i > 0 {
						hex, err := strconv.ParseUint(content[i-1], 16, 64)
						if err == nil {
							content[i-1] = fmt.Sprint(hex)
						}
					}
					content = append(content[:i], content[i+1:]...)
					i--
				} else if content[i] == "(bin)" {
					if i > 0 {
						bin, err := strconv.ParseUint(content[i-1], 2, 64)
						if err == nil {
							content[i-1] = fmt.Sprint(bin)
						}
					}
					content = append(content[:i], content[i+1:]...)
					i--
				} else if content[i] == "(up)" {
					if i > 0 {
						content[i-1] = strings.ToUpper(content[i-1])
					}
					content = append(content[:i], content[i+1:]...)
					i--
				} else if content[i] == "(low)" {
					if i > 0 {
						content[i-1] = strings.ToLower(content[i-1])
					}
					content = append(content[:i], content[i+1:]...)
					i--
				} else if content[i] == "(cap)" {
					if i > 0 {
						content[i-1] = capitalize(content[i-1])
					}
					content = append(content[:i], content[i+1:]...)
					i--
				} else if content[i] == "(low," {
					if i+1 < len(content) {
						num := strings.TrimRight(content[i+1], ")")
						nb, err := strconv.Atoi(num)
						if err == nil {
							for nb > 0 {
								if nb <= i {
									content[i-nb] = strings.ToLower(content[i-nb])
								} else {
									nb = i + 1
								}
								nb--
							}
							content = append(content[:i], content[i+2:]...)
							i--
						}
					}
				} else if content[i] == "(up," {
					if i+1 < len(content) {
						num := strings.TrimRight(content[i+1], ")")
						nb, err := strconv.Atoi(num)
						if err == nil {
							for nb > 0 {
								if nb <= i {
									content[i-nb] = strings.ToUpper(content[i-nb])
								} else {
									nb = i + 1
								}
								nb--
							}
							content = append(content[:i], content[i+2:]...)
							i--
						}
					}
				} else if content[i] == "(cap," {
					if i+1 < len(content) {
						num := strings.TrimRight(content[i+1], ")")
						nb, err := strconv.Atoi(num)
						if err == nil {
							for nb > 0 {
								if nb <= i {
									content[i-nb] = capitalize(content[i-nb])
								} else {
									nb = i + 1
								}
								nb--
							}
							content = append(content[:i], content[i+2:]...)
							i--
						}
					}
				}
			}
			for i := 0; i < len(content)-1; i++ {
				if content[i] == "a" || content[i] == "A" {
					if c, ok := nextWordStart(content, i); ok && isvowel(c) {
						if content[i] == "a" {
							content[i] = "an"
						} else {
							content[i] = "An"
						}
					}
				}
			}
			slice = append(slice, strings.Join(content, " "))
		}
	}

	result = strings.Join(slice, "\n")

	result = quote(result)

	result = formatPunctuation(result)

	return result
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("error : ur input is wrong")
		os.Exit(1)

	}
	input := os.Args[1]
	outputfile := os.Args[2]

	Bytecontent, err := os.ReadFile(input)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	output := ProcessText(string(Bytecontent))
	if err := os.WriteFile(outputfile, []byte(output), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
