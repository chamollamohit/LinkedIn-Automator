package browser

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

// NewBrowser initializes a browser instance with stealth settings

func NewBrowser() *rod.Browser {
	// 1. Configure the Launcher

	path := launcher.New().
		Headless(false).
		MustLaunch()

	// 2. Connect Rod to that specific browser instance
	browser := rod.New().ControlURL(path).MustConnect()

	// 3. ENABLE VISUAL TRACE (To avoid Detection keep it off)
	browser.Trace(false)

	// 4. Apply Stealth
	stealth.MustPage(browser)

	return browser
}
