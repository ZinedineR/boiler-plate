package random

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func generateRandomNumber() string {
	// Seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Generate a random 13-digit number
	randomNumber := 1000000000000 + r.Int63n(9999999999999-1000000000000+1)

	return fmt.Sprintf("%013d", randomNumber)
}

func GenerateCodes(word string) string {
	if len(word) < 2 {
		return "Word length must be at least 2 characters"
	}

	prefix := strings.ToUpper(word[0:1] + word[2:3])
	randomNumber := generateRandomNumber()
	code := fmt.Sprintf("%s%s", prefix, randomNumber)

	return code
}
