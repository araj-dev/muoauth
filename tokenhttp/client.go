package tokenhttp

import (
	"github.com/araj-dev/muoauth/store"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Base *http.Client
}

func NewClient(store *store.TokenStore) *Client {
	trans := clonedTransport(http.DefaultTransport)

	baseClient := &http.Client{
		Transport: &Transport{
			Base:  trans,
			Store: store,
		},
		Timeout: 10 * time.Second,
	}

	return &Client{
		Base: baseClient,
	}
}

func (c *Client) Do(id string, req *http.Request) (*http.Response, error) {
	newReq := addIDContext(id, req)
	return c.base().Do(newReq)
}

func (c *Client) Get(id string, url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	newReq := addIDContext(id, req)
	return c.base().Do(newReq)
}

func (c *Client) Post(id string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	newReq := addIDContext(id, req)
	return c.base().Do(newReq)
}

func (c *Client) Put(id string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	newReq := addIDContext(id, req)
	return c.base().Do(newReq)
}

func (c *Client) Patch(id string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}

	newReq := addIDContext(id, req)
	return c.base().Do(newReq)
}

func (c *Client) Delete(id string, url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	newReq := addIDContext(id, req)
	return c.base().Do(newReq)
}

func (c *Client) base() *http.Client {
	return c.Base
}

type CreateMeetingRequest struct {
	Type            int    `json:"type"`
	Topic           string `json:"topic"`
	Agenda          string `json:"agenda"`
	Duration        int64  `json:"duration"`
	StartTime       string `json:"start_time"`
	Timezone        string `json:"timezone"`
	DefaultPassword bool   `json:"default_password"`
}

func clonedTransport(rt http.RoundTripper) *http.Transport {
	t, ok := rt.(*http.Transport)
	if !ok {
		return nil
	}
	return t.Clone()
}

func addIDContext(id string, req *http.Request) *http.Request {
	return req.WithContext(SetID(req.Context(), id))
}
