package api

import "net/http"

const (
	PipelineOnClick    = "CLICK"
	PipelineOnEvent    = "EVENT"
	PipelineOnSchedule = "SCHEDULE"

	PipelineEventTypePush      = "PUSH"
	PipelineEventTypeCreateRef = "CREATE_REF"
	PipelineEventTypeDeleteRef = "DELETE_REF"

	PipelineTriggerConditionOnChange        = "ON_CHANGE"
	PipelineTriggerConditionOnChangeAtPath  = "ON_CHANGE_AT_PATH"
	PipelineTriggerConditionVarIs           = "VAR_IS"
	PipelineTriggerConditionVarIsNot        = "VAR_IS_NOT"
	PipelineTriggerConditionVarContains     = "VAR_CONTAINS"
	PipelineTriggerConditionVarNotContains  = "VAR_NOT_CONTAINS"
	PipelineTriggerConditionDateTime        = "DATETIME"
	PipelineTriggerConditionSuccessPipeline = "SUCCESS_PIPELINE"

	PipelinePriorityHigh   = "HIGH"
	PipelinePriorityNormal = "NORMAL"
	PipelinePriorityLow    = "LOW"
)

type Pipeline struct {
	Url                       string                      `json:"url"`
	HtmlUrl                   string                      `json:"html_url"`
	Id                        int                         `json:"id"`
	Name                      string                      `json:"name"`
	On                        string                      `json:"on"`
	Refs                      []string                    `json:"refs"`
	Events                    []*PipelineEvent            `json:"events"`
	TriggerConditions         []*PipelineTriggerCondition `json:"trigger_conditions"`
	ExecutionMessageTemplate  string                      `json:"execution_message_template"`
	LastExecutionStatus       string                      `json:"last_execution_status"`
	LastExecutionRevision     string                      `json:"last_execution_revision"`
	CreateDate                string                      `json:"create_date"`
	Priority                  string                      `json:"priority"`
	AlwaysFromScratch         bool                        `json:"always_from_scratch"`
	FailOnPrepareEnvWarning   bool                        `json:"fail_on_prepare_env_warning"`
	FetchAllRefs              bool                        `json:"fetch_all_refs"`
	AutoClearCache            bool                        `json:"auto_clear_cache"`
	NoSkipToMostRecent        bool                        `json:"no_skip_to_most_recent"`
	DoNotCreateCommitStatus   bool                        `json:"do_not_create_commit_status"`
	IgnoreFailOnProjectStatus bool                        `json:"ignore_fail_on_project_status"`
	StartDate                 string                      `json:"start_date"`
	Delay                     int                         `json:"delay"`
	CloneDepth                int                         `json:"clone_depth"`
	Cron                      string                      `json:"cron"`
	Paused                    bool                        `json:"paused"`
	Worker                    string                      `json:"worker"`
	TargetSiteUrl             string                      `json:"target_site_url"`
	Tags                      []string                    `json:"tags"`
	Project                   *Project                    `json:"project"`
	Creator                   *Member                     `json:"creator"`
}

type Pipelines struct {
	Url       string      `json:"url"`
	HtmlUrl   string      `json:"html_url"`
	Pipelines []*Pipeline `json:"pipelines"`
}

type PipelineEvent struct {
	Type string   `json:"type"`
	Refs []string `json:"refs"`
}

type PipelineTriggerCondition struct {
	TriggerCondition      string   `json:"trigger_condition"`
	TriggerConditionPaths []string `json:"trigger_condition_paths"`
	TriggerVariableKey    string   `json:"trigger_variable_key"`
	TriggerVariableValue  string   `json:"trigger_variable_value"`
	TriggerHours          []int    `json:"trigger_hours"`
	TriggerDays           []int    `json:"trigger_days"`
	ZoneId                string   `json:"zone_id"`
	TriggerProjectName    string   `json:"trigger_project_name"`
	TriggerPipelineName   string   `json:"trigger_pipeline_name"`
}

type PipelineService struct {
	client *Client
}

type PipelineOperationOptions struct {
	Name                      *string                      `json:"name,omitempty"`
	On                        *string                      `json:"on,omitempty"`
	Refs                      *[]string                    `json:"refs,omitempty"`
	Tags                      *[]string                    `json:"tags,omitempty"`
	Events                    *[]*PipelineEvent            `json:"events,omitempty"`
	TriggerConditions         *[]*PipelineTriggerCondition `json:"trigger_conditions,omitempty"`
	AlwaysFromScratch         *bool                        `json:"always_from_scratch,omitempty"`
	Priority                  *string                      `json:"priority,omitempty"`
	FailOnPrepareEnvWarning   *bool                        `json:"fail_on_prepare_env_warning,omitempty"`
	FetchAllRefs              *bool                        `json:"fetch_all_refs,omitempty"`
	AutoClearCache            *bool                        `json:"auto_clear_cache,omitempty"`
	NoSkipToMostRecent        *bool                        `json:"no_skip_to_most_recent,omitempty"`
	DoNotCreateCommitStatus   *bool                        `json:"do_not_create_commit_status,omitempty"`
	StartDate                 *string                      `json:"start_date,omitempty"`
	Delay                     *int                         `json:"delay,omitempty"`
	CloneDepth                *int                         `json:"clone_depth,omitempty"`
	Cron                      *string                      `json:"cron,omitempty"`
	Paused                    *bool                        `json:"paused,omitempty"`
	IgnoreFailOnProjectStatus *bool                        `json:"ignore_fail_on_project_status,omitempty"`
	ExecutionMessageTemplate  *string                      `json:"execution_message_template,omitempty"`
	Worker                    *string                      `json:"worker,omitempty"`
	TargetSiteUrl             *string                      `json:"target_site_url,omitempty"`
}

func (s *PipelineService) Create(domain string, projectName string, opt *PipelineOperationOptions) (*Pipeline, *http.Response, error) {
	var p *Pipeline
	resp, err := s.client.Create(s.client.NewUrlPath("/workspaces/%s/projects/%s/pipelines", domain, projectName), &opt, &p)
	return p, resp, err
}

func (s *PipelineService) Delete(domain string, projectName string, pipelineId int) (*http.Response, error) {
	return s.client.Delete(s.client.NewUrlPath("/workspaces/%s/projects/%s/pipelines/%d", domain, projectName, pipelineId))
}

func (s *PipelineService) Update(domain string, projectName string, pipelineId int, opt *PipelineOperationOptions) (*Pipeline, *http.Response, error) {
	var p *Pipeline
	resp, err := s.client.Update(s.client.NewUrlPath("/workspaces/%s/projects/%s/pipelines/%d", domain, projectName, pipelineId), &opt, &p)
	return p, resp, err
}

func (s *PipelineService) Get(domain string, projectName string, pipelineId int) (*Pipeline, *http.Response, error) {
	var p *Pipeline
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/projects/%s/pipelines/%d", domain, projectName, pipelineId), &p, nil)
	return p, resp, err
}

func (s *PipelineService) GetList(domain string, projectName string) (*Pipelines, *http.Response, error) {
	var all Pipelines
	page := 1
	perPage := 30
	for {
		var l *Pipelines
		resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/projects/%s/pipelines", domain, projectName), &l, &QueryPage{
			Page:    page,
			PerPage: perPage,
		})
		if err != nil {
			return nil, resp, err
		}
		if len(l.Pipelines) == 0 {
			break
		}
		all.Url = l.Url
		all.HtmlUrl = l.HtmlUrl
		all.Pipelines = append(all.Pipelines, l.Pipelines...)
		page += 1
	}
	return &all, nil, nil
}
