package api

import (
	"net/http"
)

type MemberService struct {
	client *Client
}

type Member struct {
	HtmlUrl        string `json:"html_url"`
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	AvatarUrl      string `json:"avatar_url"`
	Admin          bool   `json:"admin"`
	WorkspaceOwner bool   `json:"workspace_owner"`
}

type Members struct {
	Url     string    `json:"url"`
	HtmlUrl string    `json:"html_url"`
	Members []*Member `json:"members"`
}

type MemberOperationOptions struct {
	Email *string `json:"email"`
}

type MemberAdminOperationOptions struct {
	Admin *bool `json:"admin"`
}

func (s *MemberService) Get(domain string, memberId int) (*Member, *http.Response, error) {
	var m *Member
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/members/%d", domain, memberId), &m, nil)
	return m, resp, err
}

func (s *MemberService) GetList(domain string) (*Members, *http.Response, error) {
	var all Members
	page := 1
	perPage := 30
	for {
		var l *Members
		resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/members", domain), &l, &QueryPage{
			Page:    page,
			PerPage: perPage,
		})
		if err != nil {
			return nil, resp, err
		}
		if len(l.Members) == 0 {
			break
		}
		all.Url = l.Url
		all.HtmlUrl = l.HtmlUrl
		all.Members = append(all.Members, l.Members...)
		page += 1
	}
	return &all, nil, nil
}

func (s *MemberService) Create(domain string, opt *MemberOperationOptions) (*Member, *http.Response, error) {
	var m *Member
	resp, err := s.client.Create(s.client.NewUrlPath("/workspaces/%s/members", domain), &opt, &m)
	return m, resp, err
}

func (s *MemberService) Delete(domain string, memberId int) (*http.Response, error) {
	return s.client.Delete(s.client.NewUrlPath("/workspaces/%s/members/%d", domain, memberId))
}

func (s *MemberService) UpdateAdmin(domain string, memberId int, opt *MemberAdminOperationOptions) (*Member, *http.Response, error) {
	var m *Member
	resp, err := s.client.Update(s.client.NewUrlPath("/workspaces/%s/members/%d", domain, memberId), &opt, &m)
	return m, resp, err
}
