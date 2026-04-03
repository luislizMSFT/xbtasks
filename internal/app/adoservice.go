package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"dev.azure.com/xbox/xb-tasks/domain"
	"dev.azure.com/xbox/xb-tasks/internal/config"
	"dev.azure.com/xbox/xb-tasks/internal/db"
)

// ADOService fetches Azure DevOps data by shelling out to az cli.
type ADOService struct {
	db  *db.DB
	cfg *config.ConfigService
}

func NewADOService(database *db.DB, cfg *config.ConfigService) *ADOService {
	return &ADOService{db: database, cfg: cfg}
}

// runAzCli executes an az cli command and returns stdout bytes.
func (s *ADOService) runAzCli(args ...string) ([]byte, error) {
	cmd := exec.Command("az", args...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("az cli error: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("az cli not found or failed: %w", err)
	}
	return output, nil
}

// appendOrgProject conditionally adds --organization and --project flags
// when config values are set. When empty, az cli uses its own defaults.
func (s *ADOService) appendOrgProject(args []string) []string {
	if org := s.cfg.GetADOOrg(); org != "" {
		args = append(args, "--organization", "https://dev.azure.com/"+org)
	}
	if proj := s.cfg.GetADOProject(); proj != "" {
		args = append(args, "--project", proj)
	}
	return args
}

// CheckConnection verifies az cli is available and authenticated.
// Returns the authenticated user's display name, or error.
func (s *ADOService) CheckConnection() (string, error) {
	output, err := s.runAzCli("account", "show", "-o", "json")
	if err != nil {
		return "", fmt.Errorf("az cli not authenticated: %w", err)
	}
	var acct struct {
		User struct {
			Name string `json:"name"`
		} `json:"user"`
	}
	if err := json.Unmarshal(output, &acct); err != nil {
		return "", err
	}
	return acct.User.Name, nil
}

// azQueryResult represents one item returned by az boards query.
type azQueryResult struct {
	ID     int                    `json:"id"`
	URL    string                 `json:"url"`
	Fields map[string]interface{} `json:"fields"`
}

func parseWorkItem(r azQueryResult) domain.ADOWorkItem {
	str := func(key string) string {
		if v, ok := r.Fields[key]; ok {
			switch t := v.(type) {
			case string:
				return t
			case map[string]interface{}:
				if dn, ok := t["displayName"].(string); ok {
					return dn
				}
			}
		}
		return ""
	}
	priority := 2
	if v, ok := r.Fields["Microsoft.VSTS.Common.Priority"]; ok {
		switch t := v.(type) {
		case float64:
			priority = int(t)
		case json.Number:
			if n, err := t.Int64(); err == nil {
				priority = int(n)
			}
		}
	}

	return domain.ADOWorkItem{
		AdoID:       strconv.Itoa(r.ID),
		Title:       str("System.Title"),
		State:       str("System.State"),
		Type:        str("System.WorkItemType"),
		AssignedTo:  str("System.AssignedTo"),
		Priority:    priority,
		AreaPath:    str("System.AreaPath"),
		Description: str("System.Description"),
		URL:         r.URL,
	}
}

// ListMyWorkItems queries ADO for work items assigned to the current user.
func (s *ADOService) ListMyWorkItems() ([]domain.ADOWorkItem, error) {
	wiql := `SELECT [System.Id],[System.Title],[System.State],[System.WorkItemType],` +
		`[System.AssignedTo],[Microsoft.VSTS.Common.Priority],[System.AreaPath],[System.Description] ` +
		`FROM WorkItems WHERE [System.AssignedTo] = @Me ` +
		`ORDER BY [Microsoft.VSTS.Common.Priority] ASC, [System.ChangedDate] DESC`

	args := s.appendOrgProject([]string{
		"boards", "query",
		"--wiql", wiql,
		"-o", "json",
	})
	output, err := s.runAzCli(args...)
	if err != nil {
		return nil, fmt.Errorf("list my work items: %w", err)
	}

	var results []azQueryResult
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, fmt.Errorf("parse az output: %w", err)
	}

	items := make([]domain.ADOWorkItem, 0, len(results))
	for _, r := range results {
		items = append(items, parseWorkItem(r))
	}
	return items, nil
}

// GetWorkItem fetches a single work item from ADO by its ID.
func (s *ADOService) GetWorkItem(adoID string) (domain.ADOWorkItem, error) {
	args := s.appendOrgProject([]string{
		"boards", "work-item", "show",
		"--id", adoID,
		"-o", "json",
	})
	output, err := s.runAzCli(args...)
	if err != nil {
		return domain.ADOWorkItem{}, fmt.Errorf("get work item %s: %w", adoID, err)
	}

	var result azQueryResult
	if err := json.Unmarshal(output, &result); err != nil {
		return domain.ADOWorkItem{}, fmt.Errorf("parse az output: %w", err)
	}
	return parseWorkItem(result), nil
}

// SyncWorkItems fetches work items from ADO and upserts them into SQLite.
func (s *ADOService) SyncWorkItems() error {
	items, err := s.ListMyWorkItems()
	if err != nil {
		return fmt.Errorf("sync work items: %w", err)
	}

	for _, item := range items {
		if err := s.db.UpsertADOWorkItem(item); err != nil {
			log.Printf("upsert work item %s: %v", item.AdoID, err)
		}
	}
	return nil
}

// GetCachedWorkItems returns work items from the SQLite cache without calling az cli.
func (s *ADOService) GetCachedWorkItems() ([]domain.ADOWorkItem, error) {
	return s.db.ListADOWorkItems()
}
