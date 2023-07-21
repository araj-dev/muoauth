package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"time"
)

type CacheTokenStore interface {
	GetTokenSource(identity string) (*CachedTokenSource, bool)
	SetTokenSource(identity string, ct *CachedTokenSource)
	UnsetToken(identity string)
}

type PersistentTokenStore interface {
	FindRefreshTokenWithIndex(identity string, index int) (string, error)
	InsertRefreshToken(identity, refreshToken string) error
	InvokeRefreshToken(identity, refreshToken string) error
}

type TokenStore struct {
	cache      CacheTokenStore
	persistent PersistentTokenStore
	config     *oauth2.Config
}

func NewTokenStore(config *oauth2.Config, db *sql.DB) *TokenStore {
	cache := NewMapTokenStore()

	// TODO: crypto service flag
	p := &SQL{db: db, cs: nil}

	return &TokenStore{
		config:     config,
		cache:      cache,
		persistent: p,
	}
}

// GetTokenWithRetry gets token from cache or persistent store.
// If failed to get token from a cache, retry to get token from persistent store.
//
//	id: identity of token
//	maxIndex: max index of persistent store. 0 means once to get token from persistent store.
func (s *TokenStore) GetTokenWithRetry(id string, maxIndex int) (*oauth2.Token, error) {
	var t *oauth2.Token
	var err error

	// at first, try to get token from cache
	t, err = s.getTokenFromCache(id)
	if t != nil && err == nil {
		return t, nil
	}

	// if failed, retry to get token from persistent store
	for i := 0; i < maxIndex; i++ {
		var persistentErr error = nil
		t, persistentErr = s.getTokenFromPersistentWithIndex(id, i)
		if persistentErr != nil {
			continue
		}

		// if succeeded to get token from persistent store, save it to cache
		cacheTokenSource := &CachedTokenSource{
			TokenSource:  s.config.TokenSource(context.Background(), t),
			RefreshToken: t.RefreshToken,
		}
		s.cache.SetTokenSource(id, cacheTokenSource)
		return t, nil
	}

	return nil, errors.New("token not found")
}

func (s *TokenStore) getTokenFromCache(id string) (*oauth2.Token, error) {
	ts, ok := s.cache.GetTokenSource(id)
	if !ok {
		return nil, errors.New("token not found in cache")
	}
	t, err := ts.TokenSource.Token()
	if err != nil {
		s.cache.UnsetToken(id)
		return nil, err
	}
	if t.RefreshToken != ts.RefreshToken {
		ts.RefreshToken = t.RefreshToken
		err := s.persistent.InsertRefreshToken(id, t.RefreshToken)
		if err != nil {
			return nil, err
		}
	}
	log.Println("token from cache")
	return t, nil
}

func (s *TokenStore) getTokenFromPersistentWithIndex(id string, index int) (*oauth2.Token, error) {
	refreshToken, err := s.persistent.FindRefreshTokenWithIndex(id, index)
	if err != nil {
		return nil, err
	}

	ts := s.config.TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: refreshToken,
		Expiry:       time.Now(), // dummy expiry for instant refresh
	})

	t, tErr := ts.Token()
	if tErr != nil {
		fmt.Printf("error: %+v, ts: %+v ", tErr, ts)
		return nil, tErr
	}

	if t.RefreshToken != refreshToken {
		err := s.persistent.InsertRefreshToken(id, t.RefreshToken)
		if err != nil {
			return nil, err
		}
	}

	log.Println("token from persistent")
	return t, nil
}

func (s *TokenStore) invokePersistentRefreshTokens(id string, refreshTokens []string) error {
	for _, refreshToken := range refreshTokens {
		err := s.persistent.InvokeRefreshToken(id, refreshToken)
		if err != nil {
			return err
		}
	}
	return nil
}
