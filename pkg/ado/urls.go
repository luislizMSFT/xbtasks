package ado

import "fmt"

// WorkItemWebURL returns the browser URL for a work item.
func WorkItemWebURL(org, project string, id int) string {
	return fmt.Sprintf("https://dev.azure.com/%s/%s/_workitems/edit/%d", org, project, id)
}

// OrgURL returns the browser URL for an ADO organization.
func OrgURL(org string) string {
	return fmt.Sprintf("https://dev.azure.com/%s", org)
}

// ProjectURL returns the browser URL for an ADO project.
func ProjectURL(org, project string) string {
	return fmt.Sprintf("https://dev.azure.com/%s/%s", org, project)
}
