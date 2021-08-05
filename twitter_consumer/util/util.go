package util

import (
	"log"
	"regexp"
)

func RemoveSpecialChar(input string) string {
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	return re.ReplaceAllString(input, " ")
}
