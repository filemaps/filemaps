// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

// ProxyMap is virtual proxy for Map struct
type ProxyMap struct {
	*Map
	IsRead  bool
	Changed bool
}

// NewProxyMap creates a new ProxyMap
func NewProxyMap(i MapInfo) *ProxyMap {
	p := &ProxyMap{
		Map:     NewMap(i),
		IsRead:  false,
		Changed: false,
	}
	return p
}

// Write encodes Map.MapFileData to JSON file.
func (p *ProxyMap) Write() error {
	return p.writeFile(p.getFilePath())
}

func (p *ProxyMap) writeFile(path string) error {
	data, err := json.Marshal(p.Map.MapFileData)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0644)
	return err
}

// Read decodes JSON data from file to Map.MapFileData.
func (p *ProxyMap) Read() error {
	if p.IsRead == true {
		// MapFileData already read
		return nil
	}
	path := p.getFilePath()
	err := p.readFile(path)
	if err == nil {
		p.IsRead = true
		// override title, use JSON title
		p.Title = p.Title2
	}
	return err
}

func (p *ProxyMap) readFile(path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	err = p.ParseJSON(fd)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"path": path,
		}).Error("Could not read FileMap JSON file")
	}
	return err
}

// ParseJSON parses MapFileData from Reader.
func (p *ProxyMap) ParseJSON(r io.Reader) error {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	version, err := getJSONVersion(bs)
	if err != nil {
		return err
	}

	data, err := parseFileMapVersion(bs, version)
	if err != nil {
		return err
	}

	p.Map.MapFileData = *data
	return nil
}

// SetTitle sets map title.
func (p *ProxyMap) SetTitle(title string) {
	p.Title = title
	p.Changed = true
}

// SetBase sets base path for map.
func (p *ProxyMap) SetBase(base string) {
	p.Base = base
	p.Changed = true
}

// SetFile sets filename for FileMap JSON file.
func (p *ProxyMap) SetFile(file string) {
	p.File = file
	p.Changed = true
}

// AddResource adds new resource to map and assigns ID for it.
// Returns new ID.
func (p *ProxyMap) AddResource(r *Resource) int {
	p.Read()
	r.ResourceID = p.getNewResourceID()
	p.Resources[r.ResourceID] = r
	p.Changed = true
	return r.ResourceID
}

// DeleteResource deletes resource from map.
func (p *ProxyMap) DeleteResource(resourceID int) {
	p.Read()
	delete(p.Resources, resourceID)
	p.Changed = true
}

// getNewResourceID returns unassigned ResourceID.
func (p *ProxyMap) getNewResourceID() int {
	max := 0
	for id := range p.Resources {
		if id > max {
			max = id
		}
	}
	return max + 1
}

// getFilePath returns path of FileMap file
func (p *ProxyMap) getFilePath() string {
	return filepath.Join(p.Base, p.File)
}

// getJSONVersion reads version from given JSON data
func getJSONVersion(bs []byte) (float64, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(bs, &data); err != nil {
		return -1, err
	}
	return data["version"].(float64), nil
}

// Versioning

func parseFileMapVersion(bs []byte, version float64) (*MapFileData, error) {
	if version == 1 {
		var data MapFileDataV1
		if err := json.Unmarshal(bs, &data); err != nil {
			return nil, err
		}
		return convertMapFileDataV1(&data)
	}
	return nil, fmt.Errorf("Unsupported FileMap JSON version %g", version)
}

func convertMapFileDataV1(data *MapFileDataV1) (*MapFileData, error) {
	return (*MapFileData)(data), nil
}
