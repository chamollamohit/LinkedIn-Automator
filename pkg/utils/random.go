package utils

import (
	"math/rand"
	"time"
)

// RandomSleep sleeps for a random duration between min and max milliseconds
func RandomSleep(minMs, maxMs int) {
	// Initialize the random seed based on current time
	// (Otherwise it generates the same number every time)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	duration := r.Intn(maxMs-minMs) + minMs
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

// RandomInt returns a number between min and max
func RandomInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}
