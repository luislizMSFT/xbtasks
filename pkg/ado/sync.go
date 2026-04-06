package ado

import "strconv"

// StatusToADO maps local task statuses to ADO work item states, keyed by work item type.
// "default" is the fallback for types not explicitly listed.
var StatusToADO = map[string]map[string]string{
	"default": {
		"todo":        "New",
		"in_progress": "Active",
		"in_review":   "Resolved",
		"done":        "Closed",
		"blocked":     "Active",    // ADO has no "Blocked" state
		"cancelled":   "Removed",
	},
	"Task": {
		"done": "Completed", // Tasks use "Completed" not "Closed"
	},
	"Bug": {
		"done": "Closed",
	},
}

// ADOToStatus maps ADO work item states to local task statuses.
var ADOToStatus = map[string]string{
	"New":       "todo",
	"Proposed":  "todo",
	"Active":    "in_progress",
	"Resolved":  "in_review",
	"Closed":    "done",
	"Completed": "done",
	"Removed":   "cancelled",
}

// MapStatusToADO converts a local status to an ADO state, respecting work item type overrides.
func MapStatusToADO(localStatus, wiType string) string {
	// Check type-specific mapping first
	if typeMap, ok := StatusToADO[wiType]; ok {
		if adoState, ok := typeMap[localStatus]; ok {
			return adoState
		}
	}
	// Fall back to default
	if adoState, ok := StatusToADO["default"][localStatus]; ok {
		return adoState
	}
	return "New" // safe default
}

// MapADOToStatus converts an ADO state to a local status.
func MapADOToStatus(adoState string) string {
	if status, ok := ADOToStatus[adoState]; ok {
		return status
	}
	return "todo" // safe default
}

// MapPriorityToLocal maps ADO integer priority (1-4) to local priority string (P0-P3).
func MapPriorityToLocal(adoPriority int) string {
	switch adoPriority {
	case 1:
		return "P0"
	case 2:
		return "P1"
	case 3:
		return "P2"
	case 4:
		return "P3"
	default:
		return "P2"
	}
}

// MapPriorityToADO maps a local priority string (P0-P3) to ADO integer priority (1-4).
func MapPriorityToADO(localPriority string) int {
	switch localPriority {
	case "P0":
		return 1
	case "P1":
		return 2
	case "P2":
		return 3
	case "P3":
		return 4
	default:
		return 3
	}
}

// GenerateSyncDiff compares a local and remote WorkItem and returns the differences.
func GenerateSyncDiff(local, remote WorkItem) SyncDiff {
	diff := SyncDiff{
		TaskID:    local.ID,
		AdoID:     strconv.Itoa(remote.ID),
		Direction: "outbound",
	}

	if local.Title != remote.Title {
		diff.Changes = append(diff.Changes, FieldDiff{
			Field:    "title",
			Local:    local.Title,
			Remote:   remote.Title,
			Proposed: local.Title,
		})
	}

	// Compare states using mapping
	localMapped := MapStatusToADO(local.State, remote.Type)
	if localMapped != remote.State {
		diff.Changes = append(diff.Changes, FieldDiff{
			Field:    "state",
			Local:    local.State,
			Remote:   remote.State,
			Proposed: localMapped,
		})
	}

	if local.Description != remote.Description {
		diff.Changes = append(diff.Changes, FieldDiff{
			Field:    "description",
			Local:    local.Description,
			Remote:   remote.Description,
			Proposed: local.Description,
		})
	}

	return diff
}
