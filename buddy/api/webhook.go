package api

import (
	"net/http"
)

const (
	WebhookEventPush                = "PUSH"
	WebhookEventExecutionStarted    = "EXECUTION_STARTED"
	WebhookEventExecutionSuccessful = "EXECUTION_SUCCESSFUL"
	WebhookEventExecutionFailed     = "EXECUTION_FAILED"
	WebhookEventExecutionFinished   = "EXECUTION_FINISHED"
)

type WebhookService struct {
	client *Client
}

type Webhook struct {
	HtmlUrl   string            `json:"html_url"`
	Id        int               `json:"id"`
	TargetUrl string            `json:"target_url"`
	SecretKey string            `json:"secret_key"`
	Projects  []string          `json:"projects"`
	Events    []string          `json:"events"`
	Requests  []*WebhookRequest `json:"requests"`
}

type Webhooks struct {
	Url      string     `json:"url"`
	HtmlUrl  string     `json:"html_url"`
	Webhooks []*Webhook `json:"webhooks"`
}

type WebhookRequest struct {
	PostDate       string `json:"post_date"`
	ResponseStatus string `json:"response_status"`
	Body           string `json:"body"`
}

type WebhookOperationOptions struct {
	Events    []string `json:"events,omitempty"`
	Projects  []string `json:"projects,omitempty"`
	TargetUrl *string  `json:"target_url,omitempty"`
	SecretKey *string  `json:"secret_key,omitempty"`
}

func (s *WebhookService) Create(domain string, opt *WebhookOperationOptions) (*Webhook, *http.Response, error) {
	var w *Webhook
	resp, err := s.client.Create(s.client.NewUrlPath("/workspaces/%s/webhooks", domain), &opt, &w)
	return w, resp, err
}

func (s *WebhookService) Delete(domain string, webhookId int) (*http.Response, error) {
	return s.client.Delete(s.client.NewUrlPath("/workspaces/%s/webhooks/%d", domain, webhookId))
}

func (s *WebhookService) Update(domain string, webhookId int, opt *WebhookOperationOptions) (*Webhook, *http.Response, error) {
	var w *Webhook
	resp, err := s.client.Update(s.client.NewUrlPath("/workspaces/%s/webhooks/%d", domain, webhookId), &opt, &w)
	return w, resp, err
}

func (s *WebhookService) Get(domain string, webhookId int) (*Webhook, *http.Response, error) {
	var w *Webhook
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/webhooks/%d", domain, webhookId), &w, nil)
	return w, resp, err
}

func (s *WebhookService) GetList(domain string) (*Webhooks, *http.Response, error) {
	var all *Webhooks
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/webhooks", domain), &all, nil)
	return all, resp, err
}
