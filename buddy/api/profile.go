package api

import "net/http"

type Profile struct {
	Url           string `json:"url"`
	HtmlUrl       string `json:"html_url"`
	Id            int    `json:"id"`
	Name          string `json:"name"`
	AvatarUrl     string `json:"avatar_url"`
	WorkspacesUrl string `json:"workspaces_url"`
}

type ProfileOperationOptions struct {
	Name *string `json:"name"`
}

type ProfileService struct {
	client *Client
}

func (s *ProfileService) Get() (*Profile, *http.Response, error) {
	var p *Profile
	resp, err := s.client.Get(s.client.NewUrlPath("/user"), &p, nil)
	return p, resp, err
}

func (s *ProfileService) Update(opt *ProfileOperationOptions) (*Profile, *http.Response, error) {
	var p *Profile
	resp, err := s.client.Update(s.client.NewUrlPath("/user"), &opt, &p)
	return p, resp, err
}
