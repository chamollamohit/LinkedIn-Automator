package auth

import (
	"fmt"
	"linkedin-automation/internal/actions"
	"linkedin-automation/pkg/utils"
	"time"

	"github.com/go-rod/rod"
)

func Login(page *rod.Page, email, password string) {
	fmt.Println("Navigating to Login Page...")
	page.MustNavigate("https://www.linkedin.com/login")

	// Wait for the email box to be visible
	page.MustWaitLoad()

	// 1. Find Email Input and Type
	fmt.Println("Typing Email...")
	emailInput := page.MustElement("#username")
	actions.HumanType(emailInput, email)

	// 2. Random Sleep (Simulate user checking their phone or notes)
	utils.RandomSleep(1000, 2000)

	// 3. Find Password Input and Type
	fmt.Println("Typing Password...")
	passInput := page.MustElement("#password")
	actions.HumanType(passInput, password)

	// 4. Random Sleep before clicking
	utils.RandomSleep(500, 1500)

	// 5. Click Sign In
	fmt.Println("Clicking Sign In...")
	page.MustElement("button[type='submit']").MustClick()

	// 6. WAIT for the feed to load
	// This is crucial. We wait until we see the search bar (global-nav-search)
	// which proves we are inside the dashboard.
	fmt.Println("Waiting for Dashboard...")

	// We use Race() to handle 2 possibilities:
	// A: Successful Login (We see the search bar)
	// B: Security Challenge (Captcha)
	page.Race().Element(".global-nav__search").MustHandle(func(e *rod.Element) {
		fmt.Println("Login Successful! Dashboard detected.")
	}).Element("#captcha-internal").MustHandle(func(e *rod.Element) {
		fmt.Println("⚠️  SECURITY CHECK DETECTED! ⚠️")
		fmt.Println("Please solve the captcha manually in the browser window.")
		fmt.Println("I will wait 60 seconds...")
		time.Sleep(60 * time.Second)
	}).MustDo()
}
