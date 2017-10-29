// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package model

const (
	currentMapFileDataVersion = 1
)

// MapFileDataV1 is version 1 from MapFileData struct.
type MapFileDataV1 struct {
	Version int `json:"version"`
	// Title2 is a copy from MapInfo.Title
	// Title2 is stored to file so it is permanent and shareable
	Title2    string      `json:"title2"`
	Resources []*Resource `json:"resources"`
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
			Resources: make([]*Resource, 0),
		},
	}
	return m
}
