package utils

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomString generates a random string of the given length
func GenerateRandomString(length int) string {
	seed := rand.NewSource(time.Now().UnixNano()) // Seed the random generator
	r := rand.New(seed)

	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[r.Intn(len(charset))] // Select a random character from the charset
	}

	return string(randomString)
}
