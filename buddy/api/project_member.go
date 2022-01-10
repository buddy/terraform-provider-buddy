package api

import (
	"net/http"
)

type ProjectMemberService struct {
	client *Client
}

type ProjectMember struct {
	Member
	PermissionSet *Permission `json:"permission_set"`
}

type ProjectMemberOperationOptions struct {
	Id            *int                           `json:"id,omitempty"`
	PermissionSet *ProjectMemberOperationOptions `json:"permission_set,omitempty"`
}

func (s *ProjectMemberService) CreateProjectMember(domain string, projectName string, opt *ProjectMemberOperationOptions) (*ProjectMember, *http.Response, error) {
	var pm *ProjectMember
	resp, err := s.client.Create(s.client.NewUrlPath("/workspaces/%s/projects/%s/members", domain, projectName), &opt, &pm)
	return pm, resp, err
}

func (s *ProjectMemberService) DeleteProjectMember(domain string, projectName string, memberId int) (*http.Response, error) {
	return s.client.Delete(s.client.NewUrlPath("/workspaces/%s/projects/%s/members/%d", domain, projectName, memberId))
}

func (s *ProjectMemberService) GetProjectMember(domain string, projectName string, memberId int) (*ProjectMember, *http.Response, error) {
	var pm *ProjectMember
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/projects/%s/members/%d", domain, projectName, memberId), &pm, nil)
	return pm, resp, err
}

func (s *ProjectMemberService) GetProjectMembers(domain string, projectName string) (*Members, *http.Response, error) {
	var all Members
	page := 1
	perPage := 30
	for {
		var l *Members
		resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/projects/%s/members", domain, projectName), &l, &QueryPage{
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

func (s *ProjectMemberService) UpdateProjectMember(domain string, projectName string, memberId int, opt *ProjectMemberOperationOptions) (*ProjectMember, *http.Response, error) {
	var pm *ProjectMember
	resp, err := s.client.Update(s.client.NewUrlPath("/workspaces/%s/projects/%s/members/%d", domain, projectName, memberId), &opt, &pm)
	return pm, resp, err
}
