// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	log "github.com/Sirupsen/logrus"
	"github.com/filemaps/filemaps-backend/pkg/database"
	"path/filepath"
	"sort"
	"time"
)

var (
	mapManager *MapManager // singleton instance
)

// MapManager manages Maps, reads and stores them.
// MapManager works as singleton pattern.
type MapManager struct {
	Maps map[int]*ProxyMap `json:"maps"` // MapID -> Map
}

// CreateMapManager creates MapManager singleton instance.
func CreateMapManager() (*MapManager, error) {
	mapManager = &MapManager{
		Maps: make(map[int]*ProxyMap),
	}
	err := mapManager.readDB()
	return mapManager, err
}

// GetMapManager returns instance of MapManager.
func GetMapManager() *MapManager {
	if mapManager == nil {
		log.Panic("MapManager instance not created, has model.CreateMapManager() been called?")
	}
	return mapManager
}

// GetMaps returns database.FileMaps.
func (mm *MapManager) GetMaps() database.FileMaps {
	var maps database.FileMaps
	for _, pm := range mm.Maps {
		maps = append(maps, &pm.Map.FileMap)
	}
	sort.Sort(maps)
	return maps
}

// GetMap returns given map.
func (mm *MapManager) GetMap(id int) *Map {
	m := mm.Maps[id]
	if m != nil {
		m.Read()
		return m.Map
	}
	return nil
}

// AddMap adds new Map and assigns new ID for it.
func (mm *MapManager) AddMap(fm database.FileMap) (*ProxyMap, error) {
	// add entry to db and get id
	db := database.NewDB()
	if err := db.Open(); err != nil {
		return nil, err
	}
	defer db.Close()

	if err := db.AddFileMap(&fm); err != nil {
		return nil, err
	}

	pm := NewProxyMap(fm)
	mm.Maps[pm.ID] = pm
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
	fm := database.FileMap{
		Base:   base,
		File:   file,
		Opened: time.Now(),
	}
	// read title from file
	pm := NewProxyMap(fm)
	if err := pm.Read(); err != nil {
		return nil, err
	}
	fm.Title = pm.Title
	return mm.AddMap(fm)
}

func (mm *MapManager) findMapByFile(base string, file string) *ProxyMap {
	for _, pm := range mm.Maps {
		if pm.Base == base && pm.File == file {
			return pm
		}
	}
	return nil
}

// DeleteMap deletes Map with given ID.
func (mm *MapManager) DeleteMap(mapID int) error {
	db := database.NewDB()
	if err := db.Open(); err != nil {
		return err
	}
	defer db.Close()

	if err := db.DeleteFileMap(mapID); err != nil {
		return err
	}
	delete(mm.Maps, mapID)
	return nil
}

func (mm *MapManager) readDB() error {
	db := database.NewDB()
	if err := db.Open(); err != nil {
		return err
	}
	defer db.Close()

	// read filemaps db
	fms, err := db.GetFileMaps(0)
	if err != nil {
		return err
	}
	for _, fm := range fms {
		pm := NewProxyMap(fm)
		mm.Maps[pm.ID] = pm
		log.WithFields(log.Fields{
			"ID":    pm.ID,
			"Title": pm.Title,
		}).Info("filemap read from db")
	}
	return nil
}
