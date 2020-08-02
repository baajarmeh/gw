package secure

import (
	"math/rand"
	"time"
)

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var defaultLettersLen = len(defaultLetters)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomStr returns a random string with a fixed length
func RandomStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = defaultLetters[rand.Intn(defaultLettersLen)]
	}
	return string(b)
}

// RandomString returns a random string with a fixed length
func RandomString(n int, allowedChars ...[]rune) string {
	var letters []rune
	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
