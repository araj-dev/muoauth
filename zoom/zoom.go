package zoom

import (
	"database/sql"
	"github.com/araj-dev/muoauth/gen/zoomapi"
	"github.com/araj-dev/muoauth/store"
	"github.com/araj-dev/muoauth/tokenhttp"
	"golang.org/x/oauth2"
)

var baseURL = "https://api.zoom.us/v2"

type Client struct {
	tc *tokenhttp.Client
}

func New(oc *oauth2.Config, db *sql.DB) (*Client, error) {
	s := store.NewTokenStore(oc, db)
	c := tokenhttp.NewClient(s)
	return &Client{c}, nil
}

// Ptr is a helper function to return a pointer to a value
// generated golang struct has a pointer field for optional parameters.
func Ptr[T any](v T) *T {
	return &v
}

func (c *Client) Me(tokenID string) (*zoomapi.UserResponse, error) {
	req, err := zoomapi.NewUserRequest(baseURL, "me", nil)
	if err != nil {
		return nil, err
	}
	res, err := c.tc.Do(tokenID, req)
	if err != nil {
		return nil, err
	}
	return zoomapi.ParseUserResponse(res)
}

func (c *Client) CreateMeeting(tokenID string, body zoomapi.MeetingCreateJSONRequestBody) (*zoomapi.MeetingCreateResponse, error) {
	req, err := zoomapi.NewMeetingCreateRequest(baseURL, "me", body)
	if err != nil {
		return nil, err
	}
	res, err := c.tc.Do(tokenID, req)
	if err != nil {
		return nil, err
	}
	return zoomapi.ParseMeetingCreateResponse(res)
}
