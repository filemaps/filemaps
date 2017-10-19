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
	"sort"
	"time"

	"github.com/filemaps/filemaps/pkg/config"
)

const (
	// MapManagerVersion defines current MapManager version.
	MapManagerVersion = 1
	// MapsFileName defines JSON filename
	MapsFileName = "maps.json"
)

var (
	mapManager *MapManager // singleton instance
)

// MapManagerV1 is first version of MapManager struct.
type MapManagerV1 struct {
	Version   int               `json:"version"`
	MapInfos  MapInfos          `json:"maps"`
	proxyMaps map[int]*ProxyMap // Map ID -> proxyMap
}

// MapManager manages Maps, reads and stores them.
// MapManager works as singleton pattern.
type MapManager MapManagerV1

// CreateMapManager creates MapManager singleton instance.
func CreateMapManager() (*MapManager, error) {
	mapManager = &MapManager{
		Version:   MapManagerVersion,
		MapInfos:  make([]MapInfo, 0),
		proxyMaps: make(map[int]*ProxyMap),
	}
	err := mapManager.Read()
	return mapManager, err
}

// GetMapManager returns instance of MapManager.
func GetMapManager() *MapManager {
	if mapManager == nil {
		log.Panic("MapManager instance not created, has model.CreateMapManager() been called?")
	}
	return mapManager
}

// GetMaps returns MapInfos.
func (mm *MapManager) GetMaps() MapInfos {
	sort.Sort(mm.MapInfos)
	return mm.MapInfos
}

// AddMap adds new Map and assigns new ID for it.
func (mm *MapManager) AddMap(mi MapInfo) (*ProxyMap, error) {
	mi.ID = mm.getNewMapID()
	pm := NewProxyMap(mi)
	mm.MapInfos = append(mm.MapInfos, mi)
	mm.proxyMaps[mi.ID] = pm
	return pm, nil
}

// ImportMap imports new Map from filemap JSON file.
func (mm *MapManager) ImportMap(path string) (*ProxyMap, error) {
	base := filepath.Dir(path)
	file := filepath.Base(path)

	if pm := mm.findMapByFile(base, file); pm != nil {
		// given path already exists
		return pm, nil
	}

	// import new map
	info := MapInfo{
		Base:   base,
		File:   file,
		Opened: time.Now(),
	}
	// read title from file
	pm := NewProxyMap(info)
	if err := pm.Read(); err != nil {
		return nil, err
	}
	info.Title = pm.Title
	return mm.AddMap(info)
}

func (mm *MapManager) GetProxyMap(mapID int) *ProxyMap {
	// check if it is already in proxyMaps
	if mm.proxyMaps[mapID] != nil {
		return mm.proxyMaps[mapID]
	}

	// check if it is in mapInfos
	for _, mi := range mm.MapInfos {
		if mi.ID == mapID {
			// found, store it to proxyMaps
			pm := NewProxyMap(mi)
			mm.proxyMaps[mi.ID] = pm
			return pm
		}
	}

	// not found
	return nil
}

func (mm *MapManager) findMapByFile(base string, file string) *ProxyMap {
	for _, mi := range mm.MapInfos {
		if mi.Base == base && mi.File == file {
			return mm.GetProxyMap(mi.ID)
		}
	}
	return nil
}

// DeleteMap deletes Map with given ID.
func (mm *MapManager) DeleteMap(mapID int) bool {
	delete(mm.proxyMaps, mapID)

	for i, mi := range mm.MapInfos {
		if mi.ID == mapID {
			// delete from slice by swapping to last element
			mm.MapInfos[i] = mm.MapInfos[len(mm.MapInfos)-1]
			mm.MapInfos = mm.MapInfos[:len(mm.MapInfos)-1]
			return true
		}
	}
	return false
}

// getNewMapID returns unassigned MapID.
func (mm *MapManager) getNewMapID() int {
	max := 0
	for _, mi := range mm.MapInfos {
		if mi.ID > max {
			max = mi.ID
		}
	}
	return max + 1
}

// Write encodes Map.MapFileData to JSON file.
func (m *MapManager) Write() error {
	return m.writeFile(m.getFilePath())
}

// getFilePath returns full path for maps storage file
func (m *MapManager) getFilePath() string {
	return filepath.Join(config.GetDir(), MapsFileName)
}

func (m *MapManager) writeFile(path string) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0644)
	return err
}

// Read decodes JSON data from file.
func (m *MapManager) Read() error {
	return m.readFile(m.getFilePath())
}

func (m *MapManager) readFile(path string) error {
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
		}).Error("Could not read maps JSON file")
	}
	return err
}

// ParseJSON parses API keys from Reader.
func (m *MapManager) ParseJSON(r io.Reader) error {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	version, err := getJSONVersion(bs)
	if err != nil {
		return err
	}

	data, err := parseMapsVersion(bs, version)
	if err != nil {
		return err
	}

	m.MapInfos = data.MapInfos
	m.proxyMaps = make(map[int]*ProxyMap)

	return nil
}

// Versioning

func parseMapsVersion(bs []byte, version float64) (*MapManager, error) {
	if version == 1 {
		var data MapManagerV1
		if err := json.Unmarshal(bs, &data); err != nil {
			return nil, err
		}
		return convertMapManagerV1(&data)
	}
	return nil, fmt.Errorf("Unsupported maps JSON version %g", version)
}

func convertMapManagerV1(data *MapManagerV1) (*MapManager, error) {
	return (*MapManager)(data), nil
}
