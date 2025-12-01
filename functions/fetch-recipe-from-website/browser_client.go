package handler

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// fetchRecipeFromURLWithBrowser uses a headless browser to fetch the page
// This is a fallback when HTTP client is blocked by bot protection
func fetchRecipeFromURLWithBrowser(urlStr string) (*Recipe, error) {
	// Launch browser with stealth options to avoid detection
	l := launcher.New().
		Headless(true).
		Set("disable-blink-features", "AutomationControlled").
		Set("excludeSwitches", "enable-automation").
		NoSandbox(true).
		Set("disable-dev-shm-usage", "true")

	// Launch browser and get control URL
	controlURL, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	// Connect to browser with timeout
	browser := rod.New().
		ControlURL(controlURL).
		Timeout(20 * time.Second)

	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}
	defer browser.MustClose()

	// Create page
	page := browser.MustPage("")
	defer page.MustClose()

	// Set realistic browser properties to avoid detection before navigating
	// Set extra headers (key-value pairs)
	page.MustSetExtraHeaders(
		"User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language", "en-US,en;q=0.9",
	)
	page.MustSetViewport(1920, 1080, 1, false)

	// Navigate to the URL
	page.MustNavigate(urlStr)

	// Wait for page to load
	page.MustWaitLoad()

	// Wait a bit for any JavaScript to execute and render content
	time.Sleep(1 * time.Second)

	// Get the HTML content
	html, err := page.HTML()
	if err != nil {
		return nil, fmt.Errorf("failed to get page HTML: %w", err)
	}

	// Extract recipe from HTML
	return extractRecipeFromHTML(html)
}

