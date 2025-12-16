package auth

import (
	"encoding/json"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// SaveCookies extracts the current session cookies from the browser and persists them to a JSON file.
// This allows the bot to restart without needing to re-authenticate, significantly reducing detection risk.
func SaveCookies(browser *rod.Browser, filename string) error {
	// 1. Retrieve all cookies from the current browser context
	cookies, err := browser.GetCookies()
	if err != nil {
		return err
	}

	// 2. Serialize cookies to JSON format with indentation for readability
	data, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		return err
	}

	// 3. Write to disk (0644 permissions = readable by owner/group)
	return os.WriteFile(filename, data, 0644)
}

// LoadCookies reads a JSON file containing session cookies and injects them into the browser.
// This effectively "logs in" the user instantly by restoring their session state.
func LoadCookies(browser *rod.Browser, filename string) error {
	// 1. Read the cookie file
	data, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		// If file doesn't exist, we just return nil (proceed to manual login)
		return nil
	} else if err != nil {
		return err
	}

	// 2. Parse JSON into Rod's cookie structure
	var storedCookies []*proto.NetworkCookie
	if err := json.Unmarshal(data, &storedCookies); err != nil {
		return err
	}

	// 3. Convert stored cookies to NetworkCookieParam
	// Rod requires 'NetworkCookieParam' struct to SET cookies, even though it returns 'NetworkCookie' when GETTING them.
	var validCookies []*proto.NetworkCookieParam

	for _, c := range storedCookies {
		// Convert expiry time format
		expires := proto.TimeSinceEpoch(c.Expires)

		cookieParam := &proto.NetworkCookieParam{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   c.Domain,
			Path:     c.Path,
			Secure:   c.Secure,
			HTTPOnly: c.HTTPOnly,
			SameSite: c.SameSite,
			Expires:  expires,
		}
		validCookies = append(validCookies, cookieParam)
	}

	// 4. Inject the valid cookies into the browser context
	return browser.SetCookies(validCookies)
}
