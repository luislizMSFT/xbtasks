package app

import (
	"fmt"

	"github.com/pkg/browser"
)

// BrowserService opens URLs in the system's default browser.
// Exposed to the frontend as a Wails binding.
type BrowserService struct{}

func NewBrowserService() *BrowserService {
	return &BrowserService{}
}

// OpenURL opens the given URL in the system's default browser.
func (s *BrowserService) OpenURL(url string) error {
	if url == "" {
		return fmt.Errorf("empty URL")
	}
	return browser.OpenURL(url)
}
