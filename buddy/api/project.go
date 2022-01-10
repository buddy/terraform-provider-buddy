package api

import (
	"net/http"
)

type ProjectService struct {
	client *Client
}

type Project struct {
	HtmlUrl        string  `json:"html_url"`
	Name           string  `json:"name"`
	DisplayName    string  `json:"display_name"`
	Status         string  `json:"status"`
	CreateDate     string  `json:"create_date"`
	CreatedBy      *Member `json:"created_by"`
	HttpRepository string  `json:"http_repository"`
	SshRepository  string  `json:"ssh_repository"`
	SshPublicKey   string  `json:"ssh_public_key"`
	KeyFingerprint string  `json:"key_fingerprint"`
	DefaultBranch  string  `json:"default_branch"`
}

type Projects struct {
	Url      string     `json:"url"`
	HtmlUrl  string     `json:"html_url"`
	Projects []*Project `json:"projects"`
}

type ProjectIntegration struct {
	HashId string `json:"hash_id"`
}

type ProjectCreateOptions struct {
	Name              *string             `json:"name,omitempty"`
	DisplayName       *string             `json:"display_name,omitempty"`
	ExternalProjectId *string             `json:"external_project_id,omitempty"`
	GitLabProjectId   *string             `json:"git_lab_project_id,omitempty"`
	CustomRepoUrl     *string             `json:"custom_repo_url,omitempty"`
	CustomRepoUser    *string             `json:"custom_repo_user,omitempty"`
	CustomRepoPass    *string             `json:"custom_repo_pass,omitempty"`
	Integration       *ProjectIntegration `json:"integration,omitempty"`
}

type ProjectUpdateOptions struct {
	DisplayName *string `json:"display_name,omitempty"`
}

type QueryProjectList struct {
	QueryPage
	Membership bool   `url:"membership,omitempty"`
	Status     string `url:"status,omitempty"`
}

func (s *ProjectService) Create(domain string, opt *ProjectCreateOptions) (*Project, *http.Response, error) {
	var p *Project
	resp, err := s.client.Create(s.client.NewUrlPath("/workspaces/%s/projects", domain), &opt, &p)
	return p, resp, err
}

func (s *ProjectService) Delete(domain string, projectName string) (*http.Response, error) {
	return s.client.Delete(s.client.NewUrlPath("/workspaces/%s/projects/%s", domain, projectName))
}

func (s *ProjectService) Update(domain string, projectName string, opt *ProjectUpdateOptions) (*Project, *http.Response, error) {
	var p *Project
	resp, err := s.client.Update(s.client.NewUrlPath("/workspaces/%s/projects/%s", domain, projectName), &opt, &p)
	return p, resp, err
}

func (s *ProjectService) Get(domain string, projectName string) (*Project, *http.Response, error) {
	var p *Project
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/projects/%s", domain, projectName), &p, nil)
	return p, resp, err
}

func (s *ProjectService) GetList(domain string, opt *QueryProjectList) (*Projects, *http.Response, error) {
	var all Projects
	page := 1
	perPage := 30
	for {
		var l *Projects
		opt.Page = page
		opt.PerPage = perPage
		resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/projects", domain), &l, opt)
		if err != nil {
			return nil, resp, err
		}
		if len(l.Projects) == 0 {
			break
		}
		all.Url = l.Url
		all.HtmlUrl = l.HtmlUrl
		all.Projects = append(all.Projects, l.Projects...)
		page += 1
	}
	return &all, nil, nil
}
