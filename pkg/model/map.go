// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

const (
	currentMapFileDataVersion = 1
)

// MapFileDataV1 is version 1 from MapFileData struct.
type MapFileDataV1 struct {
	Version int `json:"version"`
	// Title2 is a copy from MapInfo.Title
	// Title2 is stored to file so it is permanent and shareable
	Title2    string            `json:"title2"`
	Resources map[int]*Resource `json:"resources"` // ResourceID -> Resource
}

// MapFileData struct
type MapFileData MapFileDataV1

// Map struct
type Map struct {
	MapInfo
	MapFileData
}

// NewMap creates a new Map struct
func NewMap(i MapInfo) *Map {
	m := &Map{
		MapInfo: i,
		MapFileData: MapFileData{
			Version:   currentMapFileDataVersion,
			Title2:    i.Title,
			Resources: make(map[int]*Resource),
		},
	}
	return m
}
