package app

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"dev.azure.com/xbox/xb-tasks/domain"
	"dev.azure.com/xbox/xb-tasks/internal/auth"
	"dev.azure.com/xbox/xb-tasks/internal/config"
	"dev.azure.com/xbox/xb-tasks/internal/db"
	"dev.azure.com/xbox/xb-tasks/pkg/ado"
)

// CommentService manages local task comments and selective push to ADO.
// Comments are private/local by default (D-25e) and only pushed when explicitly requested.
type CommentService struct {
	db        *db.DB
	tokenProv auth.TokenProvider
	cfg       *config.ConfigService
}

// NewCommentService creates a CommentService with database, token provider, and config.
func NewCommentService(database *db.DB, tokenProv auth.TokenProvider, cfg *config.ConfigService) *CommentService {
	return &CommentService{db: database, tokenProv: tokenProv, cfg: cfg}
}

// AddComment creates a new local comment on a task (always private by default).
func (s *CommentService) AddComment(taskID int, content string) (domain.TaskComment, error) {
	if content == "" {
		return domain.TaskComment{}, fmt.Errorf("comment content is required")
	}
	return s.db.CreateComment(taskID, content)
}

// ListComments returns all comments for a task, ordered by creation date.
func (s *CommentService) ListComments(taskID int) ([]domain.TaskComment, error) {
	return s.db.ListComments(taskID)
}

// UpdateComment updates the content of an existing comment.
func (s *CommentService) UpdateComment(id int, content string) error {
	if content == "" {
		return fmt.Errorf("comment content is required")
	}
	return s.db.UpdateComment(id, content)
}

// DeleteComment removes a comment by ID.
func (s *CommentService) DeleteComment(id int) error {
	return s.db.DeleteComment(id)
}

// PushCommentToADO pushes a local comment to the linked ADO work item (D-25f).
// After push, the comment is marked as public with the ADO comment ID.
func (s *CommentService) PushCommentToADO(commentID int) error {
	// Fetch the comment
	comment, err := s.getComment(commentID)
	if err != nil {
		return fmt.Errorf("comment %d not found: %w", commentID, err)
	}

	if comment.IsPublic {
		return fmt.Errorf("comment %d is already public (ADO comment %s)", commentID, comment.AdoCommentID)
	}

	// Get ADO link for the task
	adoID, err := s.getLinkedAdoID(comment.TaskID)
	if err != nil {
		return fmt.Errorf("task %d has no ADO link: %w", comment.TaskID, err)
	}

	// Create ADO client
	client, err := s.getClientForItem(adoID)
	if err != nil {
		return fmt.Errorf("get ADO client: %w", err)
	}

	// Push comment to ADO
	adoIDInt, _ := strconv.Atoi(adoID)
	adoComment, err := ado.AddComment(client, adoIDInt, comment.Content)
	if err != nil {
		return fmt.Errorf("push comment to ADO work item %s: %w", adoID, err)
	}

	// Mark comment as public
	adoCommentIDStr := strconv.Itoa(adoComment.ID)
	if err := s.db.MarkCommentPublic(commentID, adoCommentIDStr); err != nil {
		return fmt.Errorf("mark comment %d as public: %w", commentID, err)
	}

	log.Printf("[comments] pushed comment %d to ADO work item %s (ADO comment %s)", commentID, adoID, adoCommentIDStr)
	return nil
}

// FetchADOComments retrieves all comments from the linked ADO work item.
// If the task has no ADO link, an empty slice is returned (not an error).
func (s *CommentService) FetchADOComments(taskID int) ([]domain.ADOComment, error) {
	adoID, err := s.getLinkedAdoID(taskID)
	if err != nil {
		return []domain.ADOComment{}, nil
	}

	client, err := s.getClientForItem(adoID)
	if err != nil {
		return nil, fmt.Errorf("get ADO client: %w", err)
	}

	adoIDInt, _ := strconv.Atoi(adoID)
	comments, err := ado.GetComments(client, adoIDInt)
	if err != nil {
		return nil, fmt.Errorf("fetch comments from ADO work item %s: %w", adoID, err)
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedDate.After(comments[j].CreatedDate)
	})

	result := make([]domain.ADOComment, len(comments))
	for i, c := range comments {
		result[i] = domain.ADOComment{
			ID:          c.ID,
			Text:        c.Text,
			CreatedBy:   c.CreatedBy.DisplayName,
			CreatedDate: c.CreatedDate.Format(time.RFC3339),
		}
	}

	log.Printf("[comments] fetched %d comments from ADO work item %s for task %d", len(result), adoID, taskID)
	return result, nil
}

// ReplyToADOComment posts a new comment to the linked ADO work item.
func (s *CommentService) ReplyToADOComment(taskID int, content string) (*domain.ADOComment, error) {
	if content == "" {
		return nil, fmt.Errorf("comment content is required")
	}

	adoID, err := s.getLinkedAdoID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task %d has no ADO link: %w", taskID, err)
	}

	client, err := s.getClientForItem(adoID)
	if err != nil {
		return nil, fmt.Errorf("get ADO client: %w", err)
	}

	adoIDInt, _ := strconv.Atoi(adoID)
	adoComment, err := ado.AddComment(client, adoIDInt, content)
	if err != nil {
		return nil, fmt.Errorf("reply to ADO work item %s: %w", adoID, err)
	}

	result := &domain.ADOComment{
		ID:          adoComment.ID,
		Text:        adoComment.Text,
		CreatedBy:   adoComment.CreatedBy.DisplayName,
		CreatedDate: adoComment.CreatedDate.Format(time.RFC3339),
	}

	log.Printf("[comments] replied to ADO work item %s for task %d (ADO comment %d)", adoID, taskID, adoComment.ID)
	return result, nil
}

// --- internal helpers ---

// getComment fetches a comment by ID.
func (s *CommentService) getComment(id int) (domain.TaskComment, error) {
	var c domain.TaskComment
	var isPublic int
	err := s.db.QueryRow(`
		SELECT id, task_id, content, is_public, ado_comment_id, created_at, updated_at
		FROM task_comments WHERE id = ?`, id,
	).Scan(&c.ID, &c.TaskID, &c.Content, &isPublic, &c.AdoCommentID, &c.CreatedAt, &c.UpdatedAt)
	c.IsPublic = isPublic == 1
	return c, err
}

// getLinkedAdoID returns the first ADO ID linked to a task.
func (s *CommentService) getLinkedAdoID(taskID int) (string, error) {
	var adoID string
	err := s.db.QueryRow(
		`SELECT ado_id FROM task_ado_links WHERE task_id = ? LIMIT 1`, taskID,
	).Scan(&adoID)
	return adoID, err
}

// getClientForItem returns an ADO client for the org/project that owns the given work item.
func (s *CommentService) getClientForItem(adoID string) (*ado.Client, error) {
	id, err := strconv.Atoi(adoID)
	if err != nil {
		return nil, fmt.Errorf("invalid ADO ID %q: %w", adoID, err)
	}

	token, err := s.tokenProv.GetToken()
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}

	clients, err := ado.NewClients(token, config.GetOrgProjects())
	if err != nil {
		return nil, err
	}
	for _, c := range clients {
		_, err := ado.GetWorkItem(c, id)
		if err != nil {
			continue
		}
		return c, nil
	}
	return nil, fmt.Errorf("ADO item %s not found in any configured org/project", adoID)
}
