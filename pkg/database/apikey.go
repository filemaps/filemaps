// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package database

import (
	"fmt"
	"time"
)

// APIKey is a database struct for APIKey.
type APIKey struct {
	APIKey  string    `json:"apikey"`
	Expires time.Time `json:"expires"`
}

// CreateTableAPIKeys creates database table for API keys if it does not exist.
// API key works as primary key.
func (db *Database) CreateTableAPIKeys() error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS apikeys (
		apikey TEXT NOT NULL,
		expires INTEGER,
		PRIMARY KEY (apikey)
	)
	WITHOUT ROWID
	`
	_, err := db.h.Exec(sqlStmt)
	return err
}

// AddAPIKey inserts new API key to database.
func (db *Database) AddAPIKey(k *APIKey) error {
	stmt, err := db.h.Prepare("INSERT INTO apikeys (apikey, expires) VALUES (?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(k.APIKey, k.Expires.Unix())
	return err
}

// GetAPIKeys returns API keys from database, giving max limit rows.
// If limit is < 1, returns all rows.
func (db *Database) GetAPIKeys(limit int) ([]APIKey, error) {
	var keys []APIKey
	stmt := "SELECT apikey, expires FROM apikeys ORDER BY expires ASC"
	if limit > 0 {
		stmt += fmt.Sprintf(" LIMIT %d", limit)
	}
	rows, err := db.h.Query(stmt)
	if err != nil {
		return keys, err
	}
	defer rows.Close()

	for rows.Next() {
		k := APIKey{}
		var expires int64
		err = rows.Scan(&k.APIKey, &expires)
		k.Expires = time.Unix(expires, 0)
		if err != nil {
			return keys, err
		}
		keys = append(keys, k)
	}
	err = rows.Err()
	return keys, err
}

// GetAPIKey returns API key database row.
func (db *Database) GetAPIKey(apiKey string) (APIKey, error) {
	k := APIKey{}
	stmt, err := db.h.Prepare("SELECT apikey, expires FROM filemaps WHERE apikey = ?")
	if err != nil {
		return k, err
	}
	defer stmt.Close()

	var expires int64
	err = stmt.QueryRow(apiKey).Scan(&k.APIKey, &expires)
	k.Expires = time.Unix(expires, 0)
	return k, err
}

// UpdateAPIKey updates given API key.
func (db *Database) UpdateAPIKey(k APIKey) error {
	stmt, err := db.h.Prepare("UPDATE apikeys SET expires = ? WHERE apikey = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(k.Expires.Unix(), k.APIKey)
	return err
}

// DeleteAPIKey deletes API key from database.
func (db *Database) DeleteAPIKey(apiKey string) error {
	stmt, err := db.h.Prepare("DELETE FROM apikeys WHERE apikey = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(apiKey)
	return err
}

// DeleteExpiredAPIKeys delete all expired APIKeys.
func (db *Database) DeleteExpiredAPIKeys() error {
	stmt := fmt.Sprintf("DELETE FROM apikeys WHERE expires < %d", time.Now().Unix())
	_, err := db.h.Exec(stmt)
	return err
}
