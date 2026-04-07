package app

import (
	"fmt"
	"log"
	"strconv"

	"dev.azure.com/xbox/xb-tasks/domain"
	"dev.azure.com/xbox/xb-tasks/internal/auth"
	"dev.azure.com/xbox/xb-tasks/internal/config"
	"dev.azure.com/xbox/xb-tasks/internal/db"
	"dev.azure.com/xbox/xb-tasks/pkg/ado"
)

// ADOService fetches Azure DevOps data via REST client with Bearer token.
type ADOService struct {
	db        *db.DB
	cfg       *config.ConfigService
	tokenProv auth.TokenProvider
}

// NewADOService creates an ADOService using a TokenProvider for auth.
func NewADOService(database *db.DB, cfg *config.ConfigService, tokenProv auth.TokenProvider) *ADOService {
	return &ADOService{db: database, cfg: cfg, tokenProv: tokenProv}
}

// getClients returns REST clients for all configured org/project pairs.
func (s *ADOService) getClients() ([]*ado.Client, error) {
	token, err := s.tokenProv.GetToken()
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}
	return ado.NewClients(token, config.GetOrgProjects())
}

// CheckConnection verifies the token provider is authenticated.
// Returns a status string describing the auth method.
func (s *ADOService) CheckConnection() (string, error) {
	_, err := s.tokenProv.GetToken()
	if err != nil {
		return "", fmt.Errorf("not authenticated: %w", err)
	}
	return fmt.Sprintf("Authenticated via %s", s.tokenProv.Name()), nil
}

// adoWorkItemToDomain converts a pkg/ado.WorkItem to domain.WorkItem.
func adoWorkItemToDomain(w ado.WorkItem) domain.WorkItem {
	return domain.WorkItem{
		AdoID:       fmt.Sprintf("%d", w.ID),
		Title:       w.Title,
		State:       w.State,
		Type:        w.Type,
		AssignedTo:  w.AssignedTo,
		Priority:    w.Priority,
		AreaPath:    w.AreaPath,
		Description: w.Description,
		URL:         w.URL,
		Org:         w.Org,
		Project:     w.Project,
		ParentID:    w.ParentID,
		ChangedDate: w.ChangedDate,
	}
}

// ListMyWorkItems queries ADO for work items assigned to the current user
// across all configured org/project pairs.
func (s *ADOService) ListMyWorkItems() ([]domain.WorkItem, error) {
	clients, err := s.getClients()
	if err != nil {
		return nil, err
	}

	type result struct {
		items []domain.WorkItem
		err   error
	}
	ch := make(chan result, len(clients))
	for _, c := range clients {
		go func(c *ado.Client) {
			items, err := ado.QueryMyWorkItems(c)
			if err != nil {
				log.Printf("[ado] query %s/%s failed: %v", c.Org(), c.Project(), err)
				ch <- result{nil, err}
				return
			}
			mapped := make([]domain.WorkItem, len(items))
			for i, item := range items {
				mapped[i] = adoWorkItemToDomain(item)
			}
			ch <- result{mapped, nil}
		}(c)
	}

	var all []domain.WorkItem
	for range clients {
		r := <-ch
		if r.err == nil {
			all = append(all, r.items...)
		}
	}
	return all, nil
}

// GetWorkItem fetches a single work item from ADO by its string ID.
// Tries each configured org/project client until found.
func (s *ADOService) GetWorkItem(adoID string) (domain.WorkItem, error) {
	id, err := strconv.Atoi(adoID)
	if err != nil {
		return domain.WorkItem{}, fmt.Errorf("invalid ADO ID %q: %w", adoID, err)
	}
	clients, err := s.getClients()
	if err != nil {
		return domain.WorkItem{}, err
	}
	for _, c := range clients {
		wi, err := ado.GetWorkItem(c, id)
		if err != nil {
			continue // try next org/project
		}
		return adoWorkItemToDomain(*wi), nil
	}
	return domain.WorkItem{}, fmt.Errorf("work item %s not found in any configured org/project", adoID)
}

// SyncWorkItems fetches work items from all orgs and upserts them into SQLite.
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

// GetCachedWorkItems returns work items from the SQLite cache without calling ADO.
func (s *ADOService) GetCachedWorkItems() ([]domain.WorkItem, error) {
	return s.db.ListADOWorkItems()
}

// GetCachedWorkItem returns a single work item from the SQLite cache by ADO ID.
func (s *ADOService) GetCachedWorkItem(adoID string) (domain.WorkItem, error) {
	return s.db.GetADOWorkItem(adoID)
}

// SavedQuery represents an ADO saved WIQL query for the frontend.
type SavedQuery struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsFolder bool   `json:"isFolder"`
}

// GetSavedQueries returns saved WIQL queries from all configured org/project pairs.
func (s *ADOService) GetSavedQueries() ([]SavedQuery, error) {
	clients, err := s.getClients()
	if err != nil {
		return nil, err
	}
	var all []SavedQuery
	for _, c := range clients {
		queries, err := ado.GetSavedQueries(c)
		if err != nil {
			log.Printf("[ado] get saved queries %s/%s: %v", c.Org(), c.Project(), err)
			continue
		}
		for _, q := range queries {
			all = append(all, SavedQuery{ID: q.ID, Name: q.Name, Path: q.Path, IsFolder: q.IsFolder})
		}
	}
	return all, nil
}

// RunSavedQuery executes a saved WIQL query and returns the matching work items.
// Tries each client until one succeeds (query belongs to a specific org/project).
func (s *ADOService) RunSavedQuery(queryID string) ([]domain.WorkItem, error) {
	clients, err := s.getClients()
	if err != nil {
		return nil, err
	}
	for _, c := range clients {
		items, err := ado.RunSavedQuery(c, queryID)
		if err != nil {
			continue
		}
		var result []domain.WorkItem
		for _, item := range items {
			result = append(result, adoWorkItemToDomain(item))
		}
		return result, nil
	}
	return nil, fmt.Errorf("query %s not found in any configured org/project", queryID)
}

// GetWorkItemTree fetches assigned items and their parent hierarchy for ADO browser.
// Returns a flat list with ParentID relationships for frontend tree rendering.
func (s *ADOService) GetWorkItemTree() ([]domain.WorkItem, error) {
	items, err := s.ListMyWorkItems()
	if err != nil {
		return nil, fmt.Errorf("get work item tree: %w", err)
	}

	clients, err := s.getClients()
	if err != nil {
		return nil, err
	}

	// Collect parent IDs that need fetching (not already in the list)
	existing := make(map[string]bool)
	for _, item := range items {
		existing[item.AdoID] = true
	}

	// Fetch parents up to 3 levels deep
	toFetch := make(map[int]bool)
	for _, item := range items {
		if item.ParentID > 0 && !existing[fmt.Sprintf("%d", item.ParentID)] {
			toFetch[item.ParentID] = true
		}
	}

	for level := 0; level < 3 && len(toFetch) > 0; level++ {
		ids := make([]int, 0, len(toFetch))
		for id := range toFetch {
			ids = append(ids, id)
		}

		var fetched []ado.WorkItem
		for _, c := range clients {
			batch, err := ado.GetWorkItemsByIDs(c, ids)
			if err != nil {
				log.Printf("[ado] batch fetch parents from %s/%s failed: %v", c.Org(), c.Project(), err)
				continue
			}
			fetched = append(fetched, batch...)
		}

		toFetch = make(map[int]bool)
		for _, wi := range fetched {
			adoID := fmt.Sprintf("%d", wi.ID)
			if !existing[adoID] {
				items = append(items, adoWorkItemToDomain(wi))
				existing[adoID] = true
			}
			if wi.ParentID > 0 && !existing[fmt.Sprintf("%d", wi.ParentID)] {
				toFetch[wi.ParentID] = true
			}
		}
	}

	return items, nil
}
