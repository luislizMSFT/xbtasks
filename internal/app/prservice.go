package app

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"dev.azure.com/xbox/xb-tasks/domain"
	"dev.azure.com/xbox/xb-tasks/internal/config"
	"dev.azure.com/xbox/xb-tasks/internal/db"
	"dev.azure.com/xbox/xb-tasks/pkg/ado"
)

// PRService fetches pull requests from Azure DevOps via az cli.
type PRService struct {
	db       *db.DB
	cfg      *config.ConfigService
	meEmail  string // cached identity for --creator/--reviewer
}

func NewPRService(database *db.DB, cfg *config.ConfigService) *PRService {
	return &PRService{db: database, cfg: cfg}
}

// resolveMe gets the current user's email/UPN for PR queries.
// az repos pr list doesn't support @me, so we resolve it once.
func (s *PRService) resolveMe() (string, error) {
	if s.meEmail != "" {
		return s.meEmail, nil
	}
	output, err := ado.RunAzCli("account", "show", "-o", "json")
	if err != nil {
		return "", fmt.Errorf("resolve identity: %w", err)
	}
	var acct struct {
		User struct {
			Name string `json:"name"`
		} `json:"user"`
	}
	if err := json.Unmarshal(output, &acct); err != nil {
		return "", err
	}
	// Try ad signed-in-user show for UPN
	upnOutput, err := ado.RunAzCli("ad", "signed-in-user", "show", "--query", "userPrincipalName", "-o", "tsv")
	if err == nil && len(strings.TrimSpace(string(upnOutput))) > 0 {
		s.meEmail = strings.TrimSpace(string(upnOutput))
	} else {
		s.meEmail = acct.User.Name
	}
	return s.meEmail, nil
}

// appendOrgProject conditionally adds --organization and --project flags
// when config values are set. When empty, az cli uses its own defaults.
func (s *PRService) appendOrgProject(args []string) []string {
	if org := s.cfg.GetADOOrg(); org != "" {
		args = append(args, "--organization", "https://dev.azure.com/"+org)
	}
	if proj := s.cfg.GetADOProject(); proj != "" {
		args = append(args, "--project", proj)
	}
	return args
}

// azPR represents the JSON structure returned by az repos pr list.
type azPR struct {
	PullRequestID int    `json:"pullRequestId"`
	Title         string `json:"title"`
	Status        string `json:"status"`
	Repository    struct {
		Name string `json:"name"`
	} `json:"repository"`
	SourceRefName string `json:"sourceRefName"`
	TargetRefName string `json:"targetRefName"`
	CreatedBy     struct {
		DisplayName string `json:"displayName"`
		UniqueName  string `json:"uniqueName"`
	} `json:"createdBy"`
	Reviewers []struct {
		DisplayName string `json:"displayName"`
		UniqueName  string `json:"uniqueName"`
		Vote        int    `json:"vote"`
	} `json:"reviewers"`
	URL          string  `json:"url"`
	CreationDate string  `json:"creationDate"`
	ClosedDate   *string `json:"closedDate"`
	IsDraft      bool    `json:"isDraft"`
}

func buildPRWebURL(org, project, repo string, prID int) string {
	if org == "" || project == "" {
		return ""
	}
	return fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s/pullrequest/%d", org, project, repo, prID)
}

func (s *PRService) mapPR(raw azPR) domain.PullRequest {
	status := raw.Status
	if raw.IsDraft {
		status = "draft"
	}

	reviewersJSON, _ := json.Marshal(raw.Reviewers)

	var totalVotes int
	for _, r := range raw.Reviewers {
		totalVotes += r.Vote
	}

	prURL := buildPRWebURL(s.cfg.GetADOOrg(), s.cfg.GetADOProject(), raw.Repository.Name, raw.PullRequestID)
	if prURL == "" {
		prURL = raw.URL
	}

	pr := domain.PullRequest{
		Title:        raw.Title,
		PRURL:        prURL,
		PRNumber:     raw.PullRequestID,
		Repo:         raw.Repository.Name,
		Status:       status,
		Reviewers:    string(reviewersJSON),
		SourceBranch: stripRefPrefix(raw.SourceRefName),
		TargetBranch: stripRefPrefix(raw.TargetRefName),
		Votes:        totalVotes,
		CreatedBy:    raw.CreatedBy.DisplayName,
	}

	if t, err := time.Parse(time.RFC3339, raw.CreationDate); err == nil {
		pr.CreatedAt = t
	}
	pr.UpdatedAt = pr.CreatedAt

	if raw.ClosedDate != nil {
		if t, err := time.Parse(time.RFC3339, *raw.ClosedDate); err == nil {
			pr.MergedAt = &t
		}
	}

	return pr
}

func (s *PRService) fetchPRs(extraArgs ...string) ([]domain.PullRequest, error) {
	args := s.appendOrgProject([]string{
		"repos", "pr", "list",
	})
	args = append(args, extraArgs...)
	args = append(args, "-o", "json")

	output, err := ado.RunAzCli(args...)
	if err != nil {
		return nil, err
	}

	var raw []azPR
	if err := json.Unmarshal(output, &raw); err != nil {
		return nil, fmt.Errorf("parse az pr output: %w", err)
	}

	prs := make([]domain.PullRequest, 0, len(raw))
	for _, r := range raw {
		prs = append(prs, s.mapPR(r))
	}
	return prs, nil
}

// ListMyPRs returns PRs created by the current user.
func (s *PRService) ListMyPRs() ([]domain.PullRequest, error) {
	me, err := s.resolveMe()
	if err != nil {
		return nil, fmt.Errorf("list my PRs: %w", err)
	}
	prs, err := s.fetchPRs("--creator", me, "--status", "all", "--top", "50")
	if err != nil {
		return nil, fmt.Errorf("list my PRs: %w", err)
	}
	return prs, nil
}

// ListReviewPRs returns active PRs where the current user is a reviewer.
func (s *PRService) ListReviewPRs() ([]domain.PullRequest, error) {
	me, err := s.resolveMe()
	if err != nil {
		return nil, fmt.Errorf("list review PRs: %w", err)
	}
	prs, err := s.fetchPRs("--reviewer", me, "--status", "active", "--top", "50")
	if err != nil {
		return nil, fmt.Errorf("list review PRs: %w", err)
	}
	return prs, nil
}

// ListTeamPRs returns all active PRs for the project.
func (s *PRService) ListTeamPRs() ([]domain.PullRequest, error) {
	prs, err := s.fetchPRs("--status", "active", "--top", "100")
	if err != nil {
		return nil, fmt.Errorf("list team PRs: %w", err)
	}
	return prs, nil
}

// SyncPRs fetches my PRs and review PRs, deduplicates, and upserts to SQLite.
func (s *PRService) SyncPRs() error {
	myPRs, err := s.ListMyPRs()
	if err != nil {
		return fmt.Errorf("sync PRs (my): %w", err)
	}

	reviewPRs, err := s.ListReviewPRs()
	if err != nil {
		return fmt.Errorf("sync PRs (review): %w", err)
	}

	// Deduplicate by PR number + repo
	seen := make(map[string]bool)
	var all []domain.PullRequest
	for _, pr := range myPRs {
		key := fmt.Sprintf("%d:%s", pr.PRNumber, pr.Repo)
		if !seen[key] {
			seen[key] = true
			all = append(all, pr)
		}
	}
	for _, pr := range reviewPRs {
		key := fmt.Sprintf("%d:%s", pr.PRNumber, pr.Repo)
		if !seen[key] {
			seen[key] = true
			all = append(all, pr)
		}
	}

	for _, pr := range all {
		if err := s.db.UpsertPullRequest(pr); err != nil {
			log.Printf("upsert PR %d: %v", pr.PRNumber, err)
		}
	}
	return nil
}

// GetCachedPRs returns pull requests from the SQLite cache.
func (s *PRService) GetCachedPRs() ([]domain.PullRequest, error) {
	return s.db.ListPullRequests()
}
