// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package config

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/filemaps/filemaps-backend/pkg/filemaps"
)

const (
	// CfgFileName defines config file name
	CfgFileName = "config.json"
)

var (
	cfg *Configuration // singleton instance
)

// Configuration struct
type Configuration struct {
	Version              string `json:"version"`
	TextEditor           string `json:"textEditor"`
	TextEditorCustom1Cmd string `json:"textEditorCustom1Cmd"`
}

// CreateConfiguration creates Configuration singleton instance.
func CreateConfiguration() (*Configuration, error) {
	return readFile(getPath())
}

// GetConfiguration returns instance of Configuration.
func GetConfiguration() *Configuration {
	if cfg == nil {
		log.Panic("Configuration instance not created, has model.CreateConfiguration() been called?")
	}
	return cfg
}

// EnsureDir makes sure config path exists.
func EnsureDir() error {
	return os.MkdirAll(GetDir(), 0700)
}

func readFile(path string) (*Configuration, error) {
	fd, err := os.Open(path)
	if err != nil && os.IsNotExist(err) {
		log.Info(path + " does not exist, creating new")
		cfg = &Configuration{
			Version: filemaps.Version,
		}
		// write new config to disk
		err = cfg.Write()
		return cfg, nil
	} else if err != nil {
		return nil, err
	}
	defer fd.Close()
	return ParseJSON(fd)
}

// ParseJSON parses Configuration from given reader
func ParseJSON(r io.Reader) (*Configuration, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	cfg = &Configuration{}
	if err := json.Unmarshal(bs, cfg); err != nil {
		log.WithFields(log.Fields{
			"data": string(bs),
		}).Error(err)
		return nil, err
	}
	return cfg, nil
}

// Write writes config file
func (c *Configuration) Write() error {
	return c.writeFile(getPath())
}

func (c *Configuration) writeFile(path string) error {
	data, err := json.Marshal(cfg)

	// make sure cfg dir exists
	os.MkdirAll(filepath.Dir(path), 0700)

	err = ioutil.WriteFile(path, data, 0600)
	if err != nil {
		log.WithFields(log.Fields{
			"path": path,
		}).Error(err)
	} else {
		log.Info("config written to " + path)
	}
	return err
}

// GetDir returns directory path for config
func GetDir() string {
	switch runtime.GOOS {
	case "windows":
		if path := os.Getenv("LocalAppData"); path != "" {
			return filepath.Join(path, "FileMaps")
		}
		return filepath.Join(os.Getenv("AppData"), "FileMaps")

	// Mac OS X
	case "darwin":
		return os.Getenv("HOME") + "/Library/Application Support/FileMaps"

	// Others
	default:
		if path := os.Getenv("XDG_CONFIG_HOME"); path != "" {
			return filepath.Join(path, "filemaps")
		}
		return os.Getenv("HOME") + "/.config/filemaps"
	}
}

// returns full path of config file
func getPath() string {
	return filepath.Join(GetDir(), CfgFileName)
}
