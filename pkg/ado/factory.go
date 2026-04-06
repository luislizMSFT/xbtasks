package ado

import (
	"fmt"

	"dev.azure.com/xbox/xb-tasks/domain"
)

// FlattenOrgProjects converts config-style org/project pairs (one org → many projects)
// into individual OrgProject pairs suitable for client creation.
func FlattenOrgProjects(orgs []domain.OrgProject) []OrgProject {
	var ops []OrgProject
	for _, o := range orgs {
		for _, p := range o.Projects {
			ops = append(ops, OrgProject{Org: o.Org, Project: p})
		}
	}
	return ops
}

// NewClients creates a Client for every configured org+project combination.
func NewClients(token string, orgs []domain.OrgProject) ([]*Client, error) {
	ops := FlattenOrgProjects(orgs)
	if len(ops) == 0 {
		return nil, fmt.Errorf("no ADO org/project pairs configured — go to Settings to add org/project pairs")
	}
	return NewMultiOrgClients(ops, token), nil
}

// NewDefaultClient creates a Client for the first configured org+project.
func NewDefaultClient(token string, orgs []domain.OrgProject) (*Client, error) {
	if len(orgs) == 0 {
		return nil, fmt.Errorf("no ADO orgs configured — go to Settings to add org/project pairs")
	}
	op := orgs[0]
	if len(op.Projects) == 0 {
		return nil, fmt.Errorf("org %s has no projects configured", op.Org)
	}
	return NewClient(op.Org, op.Projects[0], token), nil
}
