// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/filemaps/filemaps-backend/pkg/config"
)

const (
	APIKeysFileName = "apikeys.json"
	APIKeysVersion  = 1
)

var (
	apiKeyManager *APIKeyManager // singleton instance
)

// APIKeyManagerV1 is first version of APIKeyManager struct.
type APIKeyManagerV1 struct {
	Version int                `json:"version"`
	APIKeys map[string]*APIKey `json:"apikeys"`
}

// APIKeyManager manages API keys.
// APIKeyManager works as singleton pattern.
type APIKeyManager APIKeyManagerV1

// CreateAPIKeyManager creates APIKeyManager singleton instance.
func CreateAPIKeyManager() (*APIKeyManager, error) {
	apiKeyManager = &APIKeyManager{
		Version: APIKeysVersion,
		APIKeys: make(map[string]*APIKey),
	}
	err := apiKeyManager.Read()
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
func (m *APIKeyManager) GetAPIKeys() []*APIKey {
	var keys []*APIKey
	for _, k := range m.APIKeys {
		keys = append(keys, k)
	}
	return keys
}

// GetAPIKey returns given API key or nil.
func (m *APIKeyManager) GetAPIKey(apiKey string) *APIKey {
	return m.APIKeys[apiKey]
}

// IsValidAPIKey returns true if given API key is valid.
func (m *APIKeyManager) IsValidAPIKey(apiKey string) bool {
	return m.GetAPIKey(apiKey) != nil
}

// CreateAPIKey generates new API key.
func (m *APIKeyManager) CreateAPIKey() string {
	k := NewAPIKey()
	m.APIKeys[k.APIKey] = &k
	return k.APIKey
}

// DeleteAPIKey deletes given API key.
func (m *APIKeyManager) DeleteAPIKey(apiKey string) {
	delete(m.APIKeys, apiKey)
}

// Write encodes Map.MapFileData to JSON file.
func (m *APIKeyManager) Write() error {
	return m.writeFile(m.getFilePath())
}

// getFilePath returns full path for API keys file
func (m *APIKeyManager) getFilePath() string {
	return filepath.Join(config.GetDir(), APIKeysFileName)
}

func (m *APIKeyManager) writeFile(path string) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0644)
	return err
}

// Read decodes JSON data from file.
func (m *APIKeyManager) Read() error {
	return m.readFile(m.getFilePath())
}

func (m *APIKeyManager) readFile(path string) error {
	fd, err := os.Open(path)
	if err != nil && os.IsNotExist(err) {
		log.Info(path + " does not exist, creating new")
		return m.writeFile(path)
	} else if err != nil {
		return err
	}
	defer fd.Close()

	err = m.ParseJSON(fd)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"path": path,
		}).Error("Could not read API keys JSON file")
	}
	return err
}

// ParseJSON parses API keys from Reader.
func (m *APIKeyManager) ParseJSON(r io.Reader) error {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	version, err := getJSONVersion(bs)
	if err != nil {
		return err
	}

	data, err := parseAPIKeysVersion(bs, version)
	if err != nil {
		return err
	}

	m.APIKeys = data.APIKeys
	return nil
}

// Versioning

func parseAPIKeysVersion(bs []byte, version float64) (*APIKeyManager, error) {
	if version == 1 {
		var data APIKeyManagerV1
		if err := json.Unmarshal(bs, &data); err != nil {
			return nil, err
		}
		return convertAPIKeyManagerV1(&data)
	}
	return nil, fmt.Errorf("Unsupported APIKeys JSON version %g", version)
}

func convertAPIKeyManagerV1(data *APIKeyManagerV1) (*APIKeyManager, error) {
	return (*APIKeyManager)(data), nil
}
