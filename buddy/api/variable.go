package api

import (
	"net/http"
)

const (
	VariableTypeVar    = "VAR"
	VariableTypeSshKey = "SSH_KEY"

	VariableSshKeyFilePlaceContainer = "CONTAINER"
	VariableSshKeyFilePlaceNone      = "NONE"
)

type VariableService struct {
	client *Client
}

type Variable struct {
	Id             int    `json:"id"`
	Key            string `json:"key"`
	Value          string `json:"value"`
	Type           string `json:"type"`
	Encrypted      bool   `json:"encrypted"`
	Settable       bool   `json:"settable"`
	Description    string `json:"description"`
	FileName       string `json:"file_name"`
	FilePath       string `json:"file_path"`
	FileChmod      string `json:"file_chmod"`
	FilePlace      string `json:"file_place"`
	PublicValue    string `json:"public_value"`
	KeyFingerprint string `json:"key_fingerprint"`
	Checksum       string `json:"checksum"`
}

type Variables struct {
	HtmlUrl   string      `json:"html_url"`
	Variables []*Variable `json:"variables"`
}

type VariableOperationOptions struct {
	Key         *string           `json:"key,omitempty"`
	Value       *string           `json:"value,omitempty"`
	Type        *string           `json:"type,omitempty"`
	Description *string           `json:"description,omitempty"`
	Settable    *bool             `json:"settable,omitempty"`
	Encrypted   *bool             `json:"encrypted,omitempty"`
	Project     *VariableProject  `json:"project,omitempty"`
	Pipeline    *VariablePipeline `json:"pipeline,omitempty"`
	Action      *VariableAction   `json:"action,omitempty"`
	FilePlace   *string           `json:"file_place,omitempty"`
	FileName    *string           `json:"file_name,omitempty"`
	FilePath    *string           `json:"file_path,omitempty"`
	FileChmod   *string           `json:"file_chmod,omitempty"`
}

type VariableGetListQuery struct {
	ProjectName string `url:"projectName,omitempty"`
	PipelineId  int    `url:"pipelineId,omitempty"`
	ActionId    int    `url:"actionId,omitempty"`
}

type VariableProject struct {
	Name string `json:"name"`
}

type VariablePipeline struct {
	Id int `json:"id"`
}

type VariableAction struct {
	Id int `json:"id"`
}

func (s *VariableService) Create(domain string, opt *VariableOperationOptions) (*Variable, *http.Response, error) {
	var v *Variable
	resp, err := s.client.Create(s.client.NewUrlPath("/workspaces/%s/variables", domain), &opt, &v)
	return v, resp, err
}

func (s *VariableService) Delete(domain string, variableId int) (*http.Response, error) {
	return s.client.Delete(s.client.NewUrlPath("/workspaces/%s/variables/%d", domain, variableId))
}

func (s *VariableService) Update(domain string, variableId int, opt *VariableOperationOptions) (*Variable, *http.Response, error) {
	var v *Variable
	resp, err := s.client.Update(s.client.NewUrlPath("/workspaces/%s/variables/%d", domain, variableId), &opt, &v)
	return v, resp, err
}

func (s *VariableService) Get(domain string, variableId int) (*Variable, *http.Response, error) {
	var v *Variable
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/variables/%d", domain, variableId), &v, nil)
	return v, resp, err
}

func (s *VariableService) GetList(domain string, opt *VariableGetListQuery) (*Variables, *http.Response, error) {
	var all *Variables
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/variables", domain), &all, &opt)
	return all, resp, err
}
