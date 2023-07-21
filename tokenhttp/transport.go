package tokenhttp

import (
	"errors"
	"github.com/araj-dev/muoauth/store"
	"net/http"
)

var ErrNoToken = errors.New("no token found")

// Transport is an http.RoundTripper that makes OAuth 2.0 HTTP requests.
type Transport struct {
	Base http.RoundTripper

	Store *store.TokenStore
}

// RoundTrip authorizes and authenticates the request with an
// access token from Transport's Source.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBodyClosed := false
	if req.Body != nil {
		defer func() {
			if !reqBodyClosed {
				req.Body.Close()
			}
		}()
	}

	// get token ID from request context
	identity, ok := GetID(req.Context())
	if !ok {
		// normal request
		return t.base().RoundTrip(req)
	}

	token, err := t.Store.GetTokenWithRetry(identity, 2)
	if err != nil {
		return nil, ErrNoToken
	}

	req2 := cloneRequest(req) // per RoundTripper contract
	token.SetAuthHeader(req2)

	// req.Body is assumed to be closed by the base RoundTripper.
	reqBodyClosed = true
	return t.base().RoundTrip(req2)
}

func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
