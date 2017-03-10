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

// Configuration struct
type Configuration struct {
	Version string
}

// New returns new Configuration
func New() Configuration {
	var cfg Configuration
	cfg.Version = filemaps.Version
	return cfg
}

// Read parses config file and returns Configuration,
// if file not found, returns a new Configuration
func Read() (Configuration, error) {
	return readFile(getPath())
}

func readFile(path string) (Configuration, error) {
	fd, err := os.Open(path)
	if err != nil && os.IsNotExist(err) {
		log.Info(path + " does not exist, creating new")
		cfg := New()
		// write new config to disk
		err = Write(cfg)
		return cfg, err
	} else if err != nil {
		return Configuration{}, err
	}
	defer fd.Close()
	return ParseJSON(fd)
}

// ParseJSON parses Configuration from given reader
func ParseJSON(r io.Reader) (Configuration, error) {
	var cfg Configuration
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		log.Error(err)
		return Configuration{}, err
	}
	if err := json.Unmarshal(bs, &cfg); err != nil {
		log.WithFields(log.Fields{
			"data": bs,
		}).Error(err)
		return Configuration{}, err
	}
	return cfg, nil
}

// Write writes given Configuration to config file
func Write(cfg Configuration) error {
	return writeFile(cfg, getPath())
}

func writeFile(cfg Configuration, path string) error {
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
