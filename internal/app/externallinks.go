package app

import (
	"regexp"

	"dev.azure.com/xbox/xb-tasks/domain"
	"dev.azure.com/xbox/xb-tasks/internal/db"
)

// ExternalLinksService manages external URL links attached to tasks.
// Automatically detects URL patterns for ICM, Grafana, ADO, and Wiki links (D-25b).
type ExternalLinksService struct {
	db *db.DB
}

// NewExternalLinksService creates an ExternalLinksService.
func NewExternalLinksService(database *db.DB) *ExternalLinksService {
	return &ExternalLinksService{db: database}
}

// urlPatterns maps URL patterns to link type identifiers.
var urlPatterns = []struct {
	Pattern *regexp.Regexp
	Type    string
}{
	{regexp.MustCompile(`portal\.microsofticm\.com`), "icm"},
	{regexp.MustCompile(`grafana\.|\.grafana\.`), "grafana"},
	{regexp.MustCompile(`dev\.azure\.com/.+/_workitems`), "ado"},
	{regexp.MustCompile(`dev\.azure\.com/.+/_wiki`), "wiki"},
	{regexp.MustCompile(`\.visualstudio\.com/.+/_wiki`), "wiki"},
}

// DetectLinkType identifies the type of URL from known patterns.
// Returns "url" as the default fallback for unrecognized URLs.
func DetectLinkType(url string) string {
	for _, p := range urlPatterns {
		if p.Pattern.MatchString(url) {
			return p.Type
		}
	}
	return "url"
}

// AddLink creates a new external link for a task with auto-detected type.
func (s *ExternalLinksService) AddLink(taskID int, url, label string) (domain.TaskLink, error) {
	linkType := DetectLinkType(url)
	return s.db.CreateLink(taskID, url, label, linkType)
}

// ListLinks returns all external links for a task.
func (s *ExternalLinksService) ListLinks(taskID int) ([]domain.TaskLink, error) {
	return s.db.ListLinks(taskID)
}

// DeleteLink removes an external link by ID.
func (s *ExternalLinksService) DeleteLink(id int) error {
	return s.db.DeleteLink(id)
}
