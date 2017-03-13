// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"time"

	"github.com/filemaps/filemaps-backend/pkg/database"
)

const (
	keyBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	apiKeyManager *APIKeyManager // singleton instance
)

// APIKeyManager manages API keys.
// APIKeyManager works as singleton pattern.
type APIKeyManager struct {
	APIKeys map[string]*database.APIKey `json:"apikeys"`
}

// CreateAPIKeyManager creates APIKeyManager singleton instance.
func CreateAPIKeyManager() (*APIKeyManager, error) {
	apiKeyManager = &APIKeyManager{
		APIKeys: make(map[string]*database.APIKey),
	}
	err := apiKeyManager.readDB()
	return apiKeyManager, err
}

// GetAPIKeyManager returns instance of APIKeyManager.
func GetAPIKeyManager() *APIKeyManager {
	if apiKeyManager == nil {
		log.Panic("APIKeyManager instance not created, has model.CreateAPIKeyManager() been called?")
	}
	return apiKeyManager
}

// GetAPIKeys returns array of database.APIKeys.
func (m *APIKeyManager) GetAPIKeys() []*database.APIKey {
	var keys []*database.APIKey
	for _, k := range m.APIKeys {
		keys = append(keys, k)
	}
	return keys
}

// GetAPIKey returns given API key or nil.
func (m *APIKeyManager) GetAPIKey(apiKey string) *database.APIKey {
	return m.APIKeys[apiKey]
}

// CreateAPIKey generates new API key.
func (m *APIKeyManager) CreateAPIKey() (string, error) {
	// add entry to db
	db := database.NewDB()
	if err := db.Open(); err != nil {
		return "", err
	}
	defer db.Close()

	k := generateAPIKey()
	if err := db.AddAPIKey(&k); err != nil {
		return k.APIKey, err
	}

	m.APIKeys[k.APIKey] = &k
	return k.APIKey, nil
}

// DeleteAPIKey deletes given API key.
func (m *APIKeyManager) DeleteAPIKey(apiKey string) error {
	db := database.NewDB()
	if err := db.Open(); err != nil {
		return err
	}
	defer db.Close()

	if err := db.DeleteAPIKey(apiKey); err != nil {
		return err
	}
	delete(m.APIKeys, apiKey)
	return nil
}

func (m *APIKeyManager) readDB() error {
	db := database.NewDB()
	if err := db.Open(); err != nil {
		return err
	}
	defer db.Close()

	if err := db.DeleteExpiredAPIKeys(); err != nil {
		return err
	}

	// read apikeys db
	keys, err := db.GetAPIKeys(0)
	if err != nil {
		return err
	}
	for _, k := range keys {
		m.APIKeys[k.APIKey] = &k
		log.WithFields(log.Fields{
			"key":     k.APIKey,
			"expires": k.Expires,
		}).Info("APIKey")
	}
	return nil
}

// generateAPIKey returns new APIKey with random API key and
// expiry date is one year from now.
func generateAPIKey() database.APIKey {
	expires := time.Now().AddDate(1, 0, 0)
	k := database.APIKey{
		APIKey:  randString(32),
		Expires: expires,
	}
	return k
}

// randString generates a random string with given length.
func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = keyBytes[rand.Intn(len(keyBytes))]
	}
	return string(b)
}
