package store

import (
	"database/sql"
)

type SQL struct {
	db *sql.DB
	cs CryptoService
}

type CryptoService interface {
	EncryptFixed(plain string) (string, error)
	Decrypt(cipher string) (string, error)
}

func (p *SQL) FindRefreshTokenWithIndex(identity string, index int) (string, error) {
	query := "SELECT refresh_token FROM google WHERE user_id = ? AND valid = 1 ORDER BY created_at DESC LIMIT 1 OFFSET ?"
	row := p.db.QueryRow(query, identity, index)
	var refreshToken string
	err := row.Scan(&refreshToken)
	if err != nil {
		return "", err
	}
	if p.cs != nil {
		decrypted, err := p.cs.Decrypt(refreshToken)
		if err != nil {
			return "", err
		}
		return decrypted, nil
	}
	return refreshToken, nil
}

func (p *SQL) InsertRefreshToken(identity string, refreshToken string) error {
	var err error
	rt := refreshToken
	if p.cs != nil {
		rt, err = p.cs.EncryptFixed(refreshToken)
		if err != nil {
			return err
		}
	}
	query := "INSERT INTO google (user_id, refresh_token, valid) VALUES (?, ?, ?)"
	_, err = p.db.Exec(query, identity, rt, 1)
	if err != nil {
		return err
	}
	return nil
}

func (p *SQL) InvokeRefreshToken(id, refreshToken string) error {
	var err error
	rt := refreshToken
	if p.cs != nil {
		rt, err = p.cs.EncryptFixed(refreshToken)
		if err != nil {
			return err
		}
	}

	query := "UPDATE google SET valid = 0 WHERE user_id = ? AND refresh_token = ?"
	_, err = p.db.Exec(query, id, rt)
	if err != nil {
		return err
	}
	return nil
}
