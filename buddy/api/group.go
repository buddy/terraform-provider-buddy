package api

import (
	"net/http"
)

type GroupService struct {
	client *Client
}

type Group struct {
	Url         string `json:"url"`
	HtmlUrl     string `json:"html_url"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Groups struct {
	Url     string   `json:"url"`
	HtmlUrl string   `json:"html_url"`
	Groups  []*Group `json:"groups"`
}

type GroupOperationOptions struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type GroupMemberOperationOptions struct {
	Id *int `json:"id"`
}

func (s *GroupService) Get(domain string, groupId int) (*Group, *http.Response, error) {
	var g *Group
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/groups/%d", domain, groupId), &g, nil)
	return g, resp, err
}

func (s *GroupService) GetList(domain string) (Groups, *http.Response, error) {
	var l Groups
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/groups", domain), &l, nil)
	return l, resp, err
}

func (s *GroupService) Create(domain string, opt *GroupOperationOptions) (*Group, *http.Response, error) {
	var g *Group
	resp, err := s.client.Create(s.client.NewUrlPath("/workspaces/%s/groups", domain), &opt, &g)
	return g, resp, err
}

func (s *GroupService) Update(domain string, groupId int, opt *GroupOperationOptions) (*Group, *http.Response, error) {
	var g *Group
	resp, err := s.client.Update(s.client.NewUrlPath("/workspaces/%s/groups/%d", domain, groupId), &opt, &g)
	return g, resp, err
}

func (s *GroupService) Delete(domain string, groupId int) (*http.Response, error) {
	return s.client.Delete(s.client.NewUrlPath("/workspaces/%s/groups/%d", domain, groupId))
}

func (s *GroupService) AddGroupMember(domain string, groupId int, opt *GroupMemberOperationOptions) (*Member, *http.Response, error) {
	var m *Member
	resp, err := s.client.Create(s.client.NewUrlPath("/workspaces/%s/groups/%d/members", domain, groupId), &opt, &m)
	return m, resp, err
}

func (s *GroupService) DeleteGroupMember(domain string, groupId int, memberId int) (*http.Response, error) {
	return s.client.Delete(s.client.NewUrlPath("/workspaces/%s/groups/%d/members/%d", domain, groupId, memberId))
}

func (s *GroupService) GetGroupMember(domain string, groupId int, memberId int) (*Member, *http.Response, error) {
	var m *Member
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/groups/%d/members/%d", domain, groupId, memberId), &m, nil)
	return m, resp, err
}

func (s *GroupService) GetGroupMembers(domain string, groupId int) (*Members, *http.Response, error) {
	var l *Members
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/groups/%d/members", domain, groupId), &l, nil)
	return l, resp, err
}
