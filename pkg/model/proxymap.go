// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/filemaps/filemaps/pkg/fileapp"
)

// ProxyMap is virtual proxy for Map struct
type ProxyMap struct {
	*Map
	IsRead  bool
	Changed bool
	// resourceIdx is resource index for internal usage
	resourceIdx map[ResourceID]int // ResourceID -> pos in Resources array
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
		if p.Title != p.Title2 {
			log.WithFields(log.Fields{
				"id":     p.ID,
				"title":  p.Title,
				"title2": p.Title2,
			}).Error("Map.Title (db) and Map.Title2 (.filemap) mismatch")
		}
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
	p.refreshResourceIdx()
	return nil
}

func (p *ProxyMap) OpenResource(r *Resource) {
	path := filepath.Join(p.Base, r.Path)
	fileapp.Open(path)
}

// GetResource returns Resource by ResourceID.
func (p *ProxyMap) GetResource(id ResourceID) *Resource {
	return p.Resources[p.resourceIdx[id]]
}

// GetResourceByPath returns Resource having given path or nil.
func (p *ProxyMap) GetResourceByPath(path string) *Resource {
	for _, res := range p.Resources {
		if path == res.Path {
			return res
		}
	}
	return nil
}

// AddResource adds new resource to map and assigns ID for it.
// Returns new ID.
func (p *ProxyMap) AddResource(r *Resource) ResourceID {
	p.Read()
	r.ResourceID = p.getNewResourceID()
	p.Resources = append(p.Resources, r)
	// update resource index
	p.resourceIdx[r.ResourceID] = len(p.Resources) - 1
	p.Changed = true
	return r.ResourceID
}

// DeleteResource deletes resource from map.
func (p *ProxyMap) DeleteResource(resourceID ResourceID) {
	p.Read()
	i := p.resourceIdx[resourceID]
	// swap element with the last one
	p.Resources[len(p.Resources)-1], p.Resources[i] = p.Resources[i], p.Resources[len(p.Resources)-1]
	// delete the last element
	p.Resources = p.Resources[:len(p.Resources)-1]
	p.refreshResourceIdx()
	p.Changed = true
}

// refreshResourceIdx refreshes resource index in var resourceIdx.
// ResourceID -> Resources array pos
func (p *ProxyMap) refreshResourceIdx() {
	p.resourceIdx = make(map[ResourceID]int)
	for i := 0; i < len(p.Resources); i++ {
		p.resourceIdx[p.Resources[i].ResourceID] = i
	}
}

// getNewResourceID returns unassigned ResourceID.
func (p *ProxyMap) getNewResourceID() ResourceID {
	var max ResourceID
	max = 0
	for id := range p.resourceIdx {
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
