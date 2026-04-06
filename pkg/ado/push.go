package ado

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// CreateWorkItem creates a new work item of the given type in ADO.
// Uses JSON Patch format (application/json-patch+json) as required by the ADO API.
func CreateWorkItem(c *Client, wiType, title, description string) (*WorkItem, error) {
	ops := []PatchOperation{
		{Op: "add", Path: "/fields/System.Title", Value: title},
		{Op: "add", Path: "/fields/System.Description", Value: description},
	}

	url := c.apiURL(fmt.Sprintf("wit/workitems/$%s?api-version=7.0", wiType))

	opsJSON, err := json.Marshal(ops)
	if err != nil {
		return nil, fmt.Errorf("marshal patch ops: %w", err)
	}

	resp, err := c.doRequest("POST", url, bytes.NewReader(opsJSON), "application/json-patch+json")
	if err != nil {
		return nil, fmt.Errorf("create work item: %w", err)
	}

	var result struct {
		ID     int            `json:"id"`
		URL    string         `json:"url"`
		Fields map[string]any `json:"fields"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("parse create response: %w", err)
	}

	wi := parseWorkItemFromAPI(result.Fields)
	wi.ID = result.ID
	wi.URL = result.URL
	wi.Org = c.org
	wi.Project = c.project
	return &wi, nil
}

// UpdateWorkItem updates a work item by ID using JSON Patch operations.
func UpdateWorkItem(c *Client, id int, ops []PatchOperation) (*WorkItem, error) {
	url := c.apiURL(fmt.Sprintf("wit/workitems/%d?api-version=7.0", id))

	resp, err := c.doPatch(url, ops)
	if err != nil {
		return nil, fmt.Errorf("update work item %d: %w", id, err)
	}

	var result struct {
		ID     int            `json:"id"`
		URL    string         `json:"url"`
		Fields map[string]any `json:"fields"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("parse update response: %w", err)
	}

	wi := parseWorkItemFromAPI(result.Fields)
	wi.ID = result.ID
	wi.URL = result.URL
	wi.Org = c.org
	wi.Project = c.project
	return &wi, nil
}

// fieldPathMap maps friendly field names to ADO system field paths.
var fieldPathMap = map[string]string{
	"title":       "/fields/System.Title",
	"state":       "/fields/System.State",
	"description": "/fields/System.Description",
	"areapath":    "/fields/System.AreaPath",
	"priority":    "/fields/Microsoft.VSTS.Common.Priority",
	"assignedto":  "/fields/System.AssignedTo",
}

// UpdateWorkItemFields is a convenience wrapper that builds PatchOperations from a
// map of friendly field names to values, then calls UpdateWorkItem.
func UpdateWorkItemFields(c *Client, id int, fields map[string]string) (*WorkItem, error) {
	ops := make([]PatchOperation, 0, len(fields))
	for key, value := range fields {
		path, ok := fieldPathMap[key]
		if !ok {
			return nil, fmt.Errorf("unknown field name: %q", key)
		}
		ops = append(ops, PatchOperation{
			Op:    "replace",
			Path:  path,
			Value: value,
		})
	}
	return UpdateWorkItem(c, id, ops)
}
