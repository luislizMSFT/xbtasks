package app

import "strings"

// stripRefPrefix removes the "refs/heads/" prefix from Azure DevOps branch references.
func stripRefPrefix(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}
