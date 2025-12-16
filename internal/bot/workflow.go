package bot

import (
	"context"
	"fmt"
	"linkedin-automation/internal/actions"
	"linkedin-automation/pkg/storage"
	"linkedin-automation/pkg/utils"
	"log/slog"
	"os"
	"time"

	"github.com/go-rod/rod"
)

// Setup a logger that prints text to the terminal
var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func RunSearchAndConnect(page *rod.Page, keyword string) {
	// Structured Log: Key-Value pairs
	logger.Info("Starting search job", "keyword", keyword)

	searchURL := fmt.Sprintf("https://www.linkedin.com/search/results/people/?keywords=%s", keyword)
	page.MustNavigate(searchURL)

	// --- PAGINATION LOOP ---
	// Change this number to process more or fewer pages.
	maxPages := 3
	for p := 1; p <= maxPages; p++ {
		logger.Info("Processing page", "page_number", p)

		// 1. WAIT FOR RESULTS
		logger.Info("Waiting for results list to load...")

		_, err := page.Timeout(10 * time.Second).Race().
			Element(".reusable-search__result-container").MustHandle(func(e *rod.Element) {
			logger.Info("✅ Success: Search results loaded.")
		}).
			Element(".search-results-container").MustHandle(func(e *rod.Element) {
			logger.Info("✅ Success: Alternate container loaded.")
		}).
			Do()

		if err != nil {
			logger.Warn("Page load timeout - could not find result list", "error", err)

		}

		// 2. Scroll
		logger.Info("Scrolling to render elements")
		page.Mouse.Scroll(0, 500, 5)
		utils.RandomSleep(500, 1000)

		// 3. Find Targets
		logger.Info("Scanning for click targets")

		targets, _ := page.ElementsX("//span[text()='Connect']")
		if len(targets) == 0 {
			targets, _ = page.ElementsX("//span[text()='Follow']")
		}
		if len(targets) == 0 {
			targets, _ = page.ElementsX("//span[contains(text(), 'Connect')]")
		}

		logger.Info("Scan complete", "targets_found", len(targets))

		// 4. Process Targets
		// Change this number to send more invites per page.
		limitPerPage := 2
		count := 0

		for _, target := range targets {
			if count >= limitPerPage {
				break
			}
			if visible, _ := target.Visible(); !visible {
				continue
			}

			candidateName := "Unknown Candidate"
			if attr, err := target.Attribute("aria-label"); err == nil && attr != nil {
				candidateName = *attr
			}
			logger.Info("Interacting with candidate", "name", candidateName)

			logger.Info("Interacting with target", "target_index", count+1)
			actions.HumanClick(page, target)
			count++

			// --- PERSISTENCE LOGGING ---
			err := storage.LogAction(keyword, "Invited/Followed", candidateName)
			if err != nil {
				logger.Error("Failed to save to history.json", "error", err)
			} else {
				logger.Info("Action persisted to database")
			}

			// 5. SMART POPUP HANDLING
			logger.Info("Checking for 'Add Note' popup")
			utils.RandomSleep(1000, 1500)

			_, errPopup := page.Timeout(2 * time.Second).Race().
				Element("[aria-label='Send without a note']").MustHandle(func(e *rod.Element) {
				logger.Info("Popup detected: 'Send without note'")
				actions.HumanClick(page, e)
			}).
				ElementX("//div[@role='dialog']//button//span[text()='Send']").MustHandle(func(e *rod.Element) {
				logger.Info("Popup detected: Generic 'Send' button")
				parent := e.MustParent()
				actions.HumanClick(page, parent)
			}).
				Do()

			if errPopup == context.DeadlineExceeded {
				logger.Info("No popup detected, continuing workflow")
			} else if errPopup != nil {
				logger.Debug("Popup check ended with non-critical error", "details", errPopup)
			}

			logger.Info("Cooling down", "duration", "2s")
			utils.RandomSleep(2000, 3000)
		}

		// --- PAGINATION LOGIC ---
		logger.Info("Looking for pagination controls")

		nextBtn, err := page.Element("[aria-label='Next']")
		if err != nil {
			logger.Info("No 'Next' button found, stopping job")
			break
		}

		if disabled, _ := nextBtn.Attribute("disabled"); disabled != nil {
			logger.Info("Reached last page of results")
			break
		}

		logger.Info("Navigating to next page")
		actions.HumanClick(page, nextBtn)
		utils.RandomSleep(3000, 5000)
	}
}
