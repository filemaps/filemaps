// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	log "github.com/Sirupsen/logrus"
	"github.com/filemaps/filemaps-backend/pkg/database"
	"sort"
)

var (
	instance *MapManager // singleton instance
)

// MapManager manages Maps, reads and stores them.
// MapManager works as singleton pattern.
type MapManager struct {
	Maps map[int]*ProxyMap `json:"maps"` // MapID -> Map
}

// CreateMapManager creates MapManager singleton instance.
func CreateMapManager() (*MapManager, error) {
	instance = &MapManager{
		Maps: make(map[int]*ProxyMap),
	}
	err := instance.readDB()
	return instance, err
}

// GetMapManager returns instance of MapManager.
func GetMapManager() *MapManager {
	if instance == nil {
		log.Panic("MapManager instance not created, has model.CreateMapManager() been called?")
	}
	return instance
}

// GetMaps returns database.FileMaps.
func (mm *MapManager) GetMaps() database.FileMaps {
	var maps database.FileMaps
	for _, pm := range mm.Maps {
		maps = append(maps, pm.Map.FileMap)
	}
	sort.Sort(maps)
	return maps
}

func (mm *MapManager) GetMap(id int) *Map {
	m := mm.Maps[id]
	if m != nil {
		m.Read()
	}
	return m.Map
}

// AddMap adds new Map and assigns new ID for it.
func (mm *MapManager) AddMap(fm *database.FileMap) error {
	// add entry to db and get id
	db := database.NewDB()
	if err := db.Open(); err != nil {
		return err
	}
	defer db.Close()

	if err := db.AddFileMap(fm); err != nil {
		return err
	}

	pm := NewProxyMap(fm)
	mm.Maps[pm.ID] = pm
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
		pm := NewProxyMap(&fm)
		mm.Maps[pm.ID] = pm
		log.WithFields(log.Fields{
			"ID":    pm.ID,
			"Title": pm.Title,
		}).Info(pm.Title)
	}
	return nil
}
