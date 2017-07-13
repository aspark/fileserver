package utils

import (
	"math/rand"
	"strings"
)

//
func IsNullOrWhiteSpace(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

var letters = "zxcvbnmasdfghjklqwertyuiopZXCVBNMASDFGHJKLQWERTYUIOP"

//
func Random(num int) string {
	var bytes = make([]byte, num)
	// rand.Read(bytes)

	var max = len(letters)
	for i := range bytes {
		bytes[i] = letters[rand.Intn(max)]
	}

	return string(bytes)
}
