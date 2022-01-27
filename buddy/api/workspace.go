package api

import (
	"net/http"
)

type WorkspaceService struct {
	client *Client
}

type WorkspaceOperationOptions struct {
	Domain         *string `json:"domain,omitempty"`
	EncryptionSalt *string `json:"encryption_salt,omitempty"`
	Name           *string `json:"name,omitempty"`
}

type Workspace struct {
	HtmlUrl    string `json:"html_url"`
	Id         int    `json:"id"`
	OwnerId    int    `json:"owner_id"`
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	Frozen     bool   `json:"frozen"`
	CreateDate string `json:"create_date"`
}

type Workspaces struct {
	Url        string       `json:"url"`
	HtmlUrl    string       `json:"html_url"`
	Workspaces []*Workspace `json:"workspaces"`
}

func (s *WorkspaceService) Create(opt *WorkspaceOperationOptions) (*Workspace, *http.Response, error) {
	var w *Workspace
	resp, err := s.client.Create(s.client.NewUrlPath("/workspaces"), &opt, &w)
	return w, resp, err
}

func (s *WorkspaceService) Delete(domain string) (*http.Response, error) {
	return s.client.Delete(s.client.NewUrlPath("workspaces/%s", domain))
}

func (s *WorkspaceService) Update(domain string, opt *WorkspaceOperationOptions) (*Workspace, *http.Response, error) {
	var w *Workspace
	resp, err := s.client.Update(s.client.NewUrlPath("/workspaces/%s", domain), &opt, &w)
	return w, resp, err
}

func (s *WorkspaceService) Get(domain string) (*Workspace, *http.Response, error) {
	var w *Workspace
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s", domain), &w, nil)
	return w, resp, err
}

func (s *WorkspaceService) GetList() (*Workspaces, *http.Response, error) {
	var all *Workspaces
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces"), &all, nil)
	return all, resp, err
}
