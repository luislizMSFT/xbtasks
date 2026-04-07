package app

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"dev.azure.com/xbox/xb-tasks/domain"
	"dev.azure.com/xbox/xb-tasks/internal/config"
	"dev.azure.com/xbox/xb-tasks/internal/db"
	"dev.azure.com/xbox/xb-tasks/pkg/ado"
)

// PipelineService fetches Azure DevOps pipeline runs via az cli.
type PipelineService struct {
	db  *db.DB
	cfg *config.ConfigService
}

func NewPipelineService(database *db.DB, cfg *config.ConfigService) *PipelineService {
	return &PipelineService{db: database, cfg: cfg}
}

// appendOrgProject conditionally adds --organization and --project flags
// when config values are set. When empty, az cli uses its own defaults.
func (s *PipelineService) appendOrgProject(args []string) []string {
	if org := s.cfg.GetADOOrg(); org != "" {
		args = append(args, "--organization", "https://dev.azure.com/"+org)
	}
	if proj := s.cfg.GetADOProject(); proj != "" {
		args = append(args, "--project", proj)
	}
	return args
}

// azPipelineRun represents the JSON structure returned by az pipelines runs list/show.
type azPipelineRun struct {
	ID     int `json:"id"`
	Pipeline struct {
		Name string `json:"name"`
	} `json:"pipeline"`
	Name         string  `json:"name"`
	Status       string  `json:"status"`
	Result       string  `json:"result"`
	SourceBranch string  `json:"sourceBranch"`
	QueueTime    string  `json:"queueTime"`
	FinishTime   *string `json:"finishTime"`
	URL          string  `json:"url"`
}

func mapPipelineRun(raw azPipelineRun, org, project string) domain.Pipeline {
	name := raw.Pipeline.Name
	if name == "" {
		name = raw.Name
	}
	if name == "" {
		name = fmt.Sprintf("Pipeline #%d", raw.ID)
	}

	url := raw.URL
	if org != "" && project != "" {
		url = fmt.Sprintf("https://dev.azure.com/%s/%s/_build/results?buildId=%d", org, project, raw.ID)
	}

	p := domain.Pipeline{
		ID:           raw.ID,
		Name:         name,
		Status:       raw.Status,
		Result:       raw.Result,
		URL:          url,
		SourceBranch: stripRefPrefix(raw.SourceBranch),
	}

	if t, err := time.Parse(time.RFC3339, raw.QueueTime); err == nil {
		p.QueueTime = t
	}
	if raw.FinishTime != nil {
		if t, err := time.Parse(time.RFC3339, *raw.FinishTime); err == nil {
			p.FinishTime = &t
		}
	}

	return p
}

// ListRecentRuns returns the most recent pipeline runs for the project.
func (s *PipelineService) ListRecentRuns() ([]domain.Pipeline, error) {
	args := s.appendOrgProject([]string{
		"pipelines", "runs", "list",
		"--top", "20",
		"-o", "json",
	})
	output, err := ado.RunAzCli(args...)
	if err != nil {
		return nil, fmt.Errorf("list pipeline runs: %w", err)
	}

	var raw []azPipelineRun
	if err := json.Unmarshal(output, &raw); err != nil {
		return nil, fmt.Errorf("parse pipeline runs: %w", err)
	}

	org := s.cfg.GetADOOrg()
	project := s.cfg.GetADOProject()
	runs := make([]domain.Pipeline, 0, len(raw))
	for _, r := range raw {
		runs = append(runs, mapPipelineRun(r, org, project))
	}
	return runs, nil
}

// GetPipelineRun returns a single pipeline run by ID.
func (s *PipelineService) GetPipelineRun(runID int) (domain.Pipeline, error) {
	args := s.appendOrgProject([]string{
		"pipelines", "runs", "show",
		"--id", strconv.Itoa(runID),
		"-o", "json",
	})
	output, err := ado.RunAzCli(args...)
	if err != nil {
		return domain.Pipeline{}, fmt.Errorf("get pipeline run %d: %w", runID, err)
	}

	var raw azPipelineRun
	if err := json.Unmarshal(output, &raw); err != nil {
		return domain.Pipeline{}, fmt.Errorf("parse pipeline run: %w", err)
	}

	return mapPipelineRun(raw, s.cfg.GetADOOrg(), s.cfg.GetADOProject()), nil
}
