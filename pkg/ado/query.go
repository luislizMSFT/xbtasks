package ado

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// QueryMyWorkItems queries ADO for work items assigned to the current user
// that are not in a terminal state, ordered by most recently changed.
func QueryMyWorkItems(c *Client) ([]WorkItem, error) {
	wiql := `SELECT [System.Id] FROM WorkItems WHERE [System.AssignedTo] = @Me AND [System.State] NOT IN ('Closed','Removed','Completed','Cut') ORDER BY [System.ChangedDate] DESC`

	url := c.apiURL("wit/wiql?api-version=7.0")

	body := map[string]string{"query": wiql}
	resp, err := c.doJSON("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("WIQL query: %w", err)
	}

	var wiqlResp WIQLResponse
	if err := decodeResponse(resp, &wiqlResp); err != nil {
		return nil, fmt.Errorf("parse WIQL response: %w", err)
	}

	if len(wiqlResp.WorkItems) == 0 {
		return []WorkItem{}, nil
	}

	ids := make([]int, len(wiqlResp.WorkItems))
	for i, ref := range wiqlResp.WorkItems {
		ids[i] = ref.ID
	}

	return GetWorkItemsByIDs(c, ids)
}

// GetWorkItemsByIDs batch-fetches work items by their IDs (max 200 per request).
func GetWorkItemsByIDs(c *Client, ids []int) ([]WorkItem, error) {
	var allItems []WorkItem

	// Batch in groups of 200
	for i := 0; i < len(ids); i += 200 {
		end := i + 200
		if end > len(ids) {
			end = len(ids)
		}
		batch := ids[i:end]

		idStrs := make([]string, len(batch))
		for j, id := range batch {
			idStrs[j] = strconv.Itoa(id)
		}

		url := c.apiURL(fmt.Sprintf("wit/workitems?ids=%s&$expand=none&api-version=7.0", strings.Join(idStrs, ",")))

		resp, err := c.doRequest("GET", url, nil, "")
		if err != nil {
			return nil, fmt.Errorf("batch fetch work items: %w", err)
		}

		var result struct {
			Value []struct {
				ID     int            `json:"id"`
				URL    string         `json:"url"`
				Fields map[string]any `json:"fields"`
			} `json:"value"`
		}
		if err := decodeResponse(resp, &result); err != nil {
			return nil, fmt.Errorf("parse batch response: %w", err)
		}

		for _, raw := range result.Value {
			wi := parseWorkItemFromAPI(raw.Fields)
			wi.ID = raw.ID
			wi.URL = raw.URL
			wi.Org = c.org
			wi.Project = c.project
			allItems = append(allItems, wi)
		}
	}

	return allItems, nil
}

// GetWorkItem fetches a single work item by ID with all fields.
func GetWorkItem(c *Client, id int) (*WorkItem, error) {
	url := c.apiURL(fmt.Sprintf("wit/workitems/%d?api-version=7.0", id))

	resp, err := c.doRequest("GET", url, nil, "")
	if err != nil {
		return nil, fmt.Errorf("get work item %d: %w", id, err)
	}

	var result struct {
		ID     int            `json:"id"`
		URL    string         `json:"url"`
		Fields map[string]any `json:"fields"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("parse work item response: %w", err)
	}

	wi := parseWorkItemFromAPI(result.Fields)
	wi.ID = result.ID
	wi.URL = result.URL
	wi.Org = c.org
	wi.Project = c.project
	return &wi, nil
}

// GetWorkItemChildren fetches child work items of a given parent.
func GetWorkItemChildren(c *Client, parentID int) ([]WorkItem, error) {
	wiql := fmt.Sprintf(`SELECT [System.Id] FROM WorkItems WHERE [System.Parent] = %d ORDER BY [System.ChangedDate] DESC`, parentID)

	url := c.apiURL("wit/wiql?api-version=7.0")
	body := map[string]string{"query": wiql}

	resp, err := c.doJSON("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("WIQL children query: %w", err)
	}

	var wiqlResp WIQLResponse
	if err := decodeResponse(resp, &wiqlResp); err != nil {
		return nil, fmt.Errorf("parse WIQL children response: %w", err)
	}

	if len(wiqlResp.WorkItems) == 0 {
		return []WorkItem{}, nil
	}

	ids := make([]int, len(wiqlResp.WorkItems))
	for i, ref := range wiqlResp.WorkItems {
		ids[i] = ref.ID
	}

	return GetWorkItemsByIDs(c, ids)
}

// parseWorkItemFromAPI maps ADO REST API field names to the WorkItem struct.
func parseWorkItemFromAPI(fields map[string]any) WorkItem {
	wi := WorkItem{
		Title:       stringField(fields, "System.Title"),
		State:       stringField(fields, "System.State"),
		Type:        stringField(fields, "System.WorkItemType"),
		AreaPath:    stringField(fields, "System.AreaPath"),
		Description: stringField(fields, "System.Description"),
		Priority:    intField(fields, "Microsoft.VSTS.Common.Priority"),
		ParentID:    intField(fields, "System.Parent"),
	}

	// AssignedTo can be a string or an object with displayName
	if v, ok := fields["System.AssignedTo"]; ok {
		switch t := v.(type) {
		case string:
			wi.AssignedTo = t
		case map[string]any:
			if dn, ok := t["displayName"].(string); ok {
				wi.AssignedTo = dn
			}
		}
	}

	// ChangedDate
	if v, ok := fields["System.ChangedDate"].(string); ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			wi.ChangedDate = t
		}
	}

	return wi
}

// stringField extracts a string value from the ADO fields map.
func stringField(fields map[string]any, key string) string {
	if v, ok := fields[key].(string); ok {
		return v
	}
	return ""
}

// intField extracts an int value from the ADO fields map.
func intField(fields map[string]any, key string) int {
	if v, ok := fields[key]; ok {
		switch t := v.(type) {
		case float64:
			return int(t)
		}
	}
	return 0
}
