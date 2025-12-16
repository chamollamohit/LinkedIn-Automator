package main

import (
	"fmt"
	"linkedin-automation/internal/actions"
	"linkedin-automation/internal/auth"
	"linkedin-automation/internal/bot"
	"linkedin-automation/internal/browser"
	"linkedin-automation/pkg/config"
	"os"
	"time"
)

// isBusinessHours checks if the current time falls within standard working hours (9 AM - 6 PM, Mon-Fri).
// STEALTH STRATEGY: This implements the "Activity Scheduling" requirement.
// By operating only during business hours, we simulate realistic human work schedules and reduce the probability of being flagged as a bot.
func isBusinessHours() bool {
	now := time.Now()
	// Check strictly for 9 AM to 6 PM window
	if now.Hour() < 9 || now.Hour() >= 18 {
		return false
	}
	// Exclude weekends (Saturday and Sunday)
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return false
	}
	return true
}

func main() {
	// ---------------------------------------------------------
	// 1. SAFETY & STEALTH CHECKS
	// ---------------------------------------------------------
	// Strictly enforce operation windows to mimic human behavior.
	if !isBusinessHours() {
		fmt.Println("⚠️  Stealth Mode: Outside business hours (9AM-6PM). Exiting to avoid detection.")
		// We exit immediately to protect the account from suspicious off-hours activity.
		os.Exit(0)
	}

	// ---------------------------------------------------------
	// 2. CONFIGURATION & SETUP
	// ---------------------------------------------------------
	// Load sensitive credentials from environment variables (Modular Config).
	cfg := config.LoadConfig()
	fmt.Printf("Target Email: %s\n", cfg.LinkedInEmail)

	// Initialize the browser engine with Stealth settings (Fingerprint masking).
	browserInstance := browser.NewBrowser()
	// Ensure browser cleanup happens even if the program crashes.
	defer browserInstance.MustClose()

	// ---------------------------------------------------------
	// 3. AUTHENTICATION (STATE PERSISTENCE)
	// ---------------------------------------------------------
	fmt.Println("Checking session...")

	// Attempt to restore a previous session to avoid repeated logins.
	// This satisfies the "Persist session cookies" requirement.
	auth.LoadCookies(browserInstance, "cookies.json")

	// Navigate to the main feed to verify session validity.
	page := browserInstance.MustPage("https://www.linkedin.com/feed/")
	page.MustWaitLoad()

	// Intelligent Login Flow:
	// Only attempt login if the feed didn't load (user is redirected to login page).
	if actions.IsLoginRequired(page) {
		fmt.Println("Session expired or missing. Initiating login sequence...")
		auth.Login(page, cfg.LinkedInEmail, cfg.LinkedInPassword)

		// Save the new session state for future runs.
		auth.SaveCookies(browserInstance, "cookies.json")
	} else {
		fmt.Println("✅ Logged in successfully via Cookies.")
	}

	// ---------------------------------------------------------
	// 4. CORE AUTOMATION WORKFLOW
	// ---------------------------------------------------------
	// Execute the primary "Search & Connect" business logic.
	// We target "IT Recruiter" as the keyword for this demonstration.
	fmt.Println("Starting Automation Job...")
	bot.RunSearchAndConnect(page, "IT Recruiter")

	// Graceful shutdown with a delay to allow pending network requests to finish.
	fmt.Println("Job Complete. Closing in 10s...")
	time.Sleep(10 * time.Second)
}
