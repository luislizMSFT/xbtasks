package ado

import "fmt"

// commentsResponse represents the ADO comments list response.
type commentsResponse struct {
	Comments []Comment `json:"comments"`
}

// GetComments fetches all comments for a work item.
func GetComments(c *Client, workItemID int) ([]Comment, error) {
	url := c.apiURL(fmt.Sprintf("wit/workitems/%d/comments?api-version=7.0-preview.3", workItemID))

	resp, err := c.doRequest("GET", url, nil, "")
	if err != nil {
		return nil, fmt.Errorf("get comments for work item %d: %w", workItemID, err)
	}

	var result commentsResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("parse comments response: %w", err)
	}

	return result.Comments, nil
}

// AddComment adds a new comment to a work item.
func AddComment(c *Client, workItemID int, text string) (*Comment, error) {
	url := c.apiURL(fmt.Sprintf("wit/workitems/%d/comments?api-version=7.0-preview.3", workItemID))

	body := map[string]string{"text": text}
	resp, err := c.doJSON("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("add comment to work item %d: %w", workItemID, err)
	}

	var comment Comment
	if err := decodeResponse(resp, &comment); err != nil {
		return nil, fmt.Errorf("parse add comment response: %w", err)
	}

	return &comment, nil
}
