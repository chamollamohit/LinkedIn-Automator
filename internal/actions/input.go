package actions

import (
	"linkedin-automation/pkg/utils"

	"github.com/go-rod/rod"
)

// HumanType types text into an element one character at a time with random delays
func HumanType(element *rod.Element, text string) {
	// Loop through every character in the string
	for _, char := range text {
		// 1. Type the single character
		element.Input(string(char))

		// 2. Random sleep between 50ms and 150ms
		// This mimics average human typing speed
		utils.RandomSleep(50, 150)

		// 3. Occasionally pause longer (simulating "thinking")
		// 10% chance to wait extra time
		if utils.RandomInt(0, 100) < 10 {
			utils.RandomSleep(300, 600)
		}
	}
}

// IsLoginRequired checks if the current URL suggests we are logged out
func IsLoginRequired(page *rod.Page) bool {
	info := page.MustInfo()
	url := info.URL

	// If the URL is NOT the feed, we probably need to login.
	return url != "https://www.linkedin.com/feed/" && url != "https://www.linkedin.com/feed"
}
