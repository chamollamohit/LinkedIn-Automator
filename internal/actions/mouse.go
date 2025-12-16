package actions

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
)

// HumanMove moves the mouse to an element using Rod's built-in Hover
func HumanMove(page *rod.Page, element *rod.Element) {
	// 1. Move the mouse to the element
	// MustHover automatically calculates the center of the element
	// and sends real mouse movement events to get there.
	element.MustHover()

	// 2. Add a random "micro-pause"
	// Humans don't freeze instantly after moving. We wiggle/pause slightly.
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)+50))
}

// HumanClick combines moving and clicking
func HumanClick(page *rod.Page, element *rod.Element) {
	// 1. Move to the element
	HumanMove(page, element)

	// 2. Pause before clicking (The "Think Time")
	// This prevents the "Instant Click" bot detection.
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)+100))

	// 3. Click
	element.MustClick()
}
