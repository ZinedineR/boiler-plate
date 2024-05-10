package text

import (
	"strconv"
	"strings"
)

// Abbreviate is a function that takes a sentence as a string and returns a string.
// The returned string is an abbreviation of the input sentence, where each word in the sentence is abbreviated to its first character.
// If the sentence is empty or contains only spaces, the function returns an empty string.
func Abbreviate(sentence string) string {
	words := strings.Split(sentence, " ")
	var result string
	for _, word := range words {
		if len(word) != 0 {
			result += string(word[0])
		}
	}
	return result
}

func StrToInt(idString string) (int64, error) {
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
