package store

import "golang.org/x/oauth2"

type MapTokenStore struct {
	m map[string]*CachedTokenSource
}

type CachedTokenSource struct {
	RefreshToken string // For checking if token is refreshed
	TokenSource  oauth2.TokenSource
}

func NewMapTokenStore() *MapTokenStore {
	ms := &MapTokenStore{
		m: make(map[string]*CachedTokenSource),
	}

	return ms
}

func (m *MapTokenStore) GetTokenSource(identity string) (*CachedTokenSource, bool) {
	ts, ok := m.m[identity]
	return ts, ok
}

func (m *MapTokenStore) SetTokenSource(identity string, ct *CachedTokenSource) {
	m.m[identity] = ct
}

func (m *MapTokenStore) UnsetToken(identity string) {
	delete(m.m, identity)
}
