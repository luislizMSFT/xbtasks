package provider

import "dev.azure.com/xbox/xb-tasks/domain"

// WorkItemProvider queries and mutates work items from an external provider (ADO, GitHub, etc).
type WorkItemProvider interface {
	QueryMine() ([]domain.WorkItem, error)
	GetByID(id int) (*domain.WorkItem, error)
	GetByIDs(ids []int) ([]domain.WorkItem, error)
	GetChildren(parentID int) ([]domain.WorkItem, error)
	Create(wiType, title, description string) (*domain.WorkItem, error)
	Update(id int, ops []PatchOp) (*domain.WorkItem, error)
}

// PRProvider queries pull requests from an external provider.
type PRProvider interface {
	ListMine(user string) ([]domain.PullRequest, error)
	ListReviews(user string) ([]domain.PullRequest, error)
	ListTeam() ([]domain.PullRequest, error)
}

// PipelineProvider queries CI/CD pipeline runs.
type PipelineProvider interface {
	ListRecentRuns(top int) ([]domain.Pipeline, error)
}

// CommentProvider manages work item comments.
type CommentProvider interface {
	GetComments(workItemID int) ([]domain.ADOComment, error)
	AddComment(workItemID int, text string) (*domain.ADOComment, error)
}

// PatchOp represents a field update operation.
type PatchOp struct {
	Op    string
	Path  string
	Value any
}
