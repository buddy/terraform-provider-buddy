package api

import "net/http"

type SourceService struct {
	client *Client
}

type SourceFileOperationOptions struct {
	Content *string `json:"content,omitempty"`
	Message *string `json:"message,omitempty"`
	Path    *string `json:"path,omitempty"`
	Branch  *string `json:"branch,omitempty"`
}

type SourceContent struct {
}

func (s *SourceService) CreateFile(domain string, projectName string, opt *SourceFileOperationOptions) (*http.Response, error) {
	var c *SourceContent
	resp, err := s.client.Create(s.client.NewUrlPath("/workspaces/%s/projects/%s/repository/contents", domain, projectName), &opt, &c)
	return resp, err
}
