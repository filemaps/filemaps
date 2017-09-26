// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	log "github.com/Sirupsen/logrus"

	"github.com/filemaps/filemaps/pkg/fileapp"
)

// ResourceID is unique in Map, identifies Resource
type ResourceID int

// ResourceType defines type of Resource
type ResourceType int

// Resource type enum
const (
	ResourceFile ResourceType = iota
	ResourceDir
)

// Converts ResourceType to string
func (r ResourceType) String() string {
	switch r {
	case ResourceFile:
		return "file"
	case ResourceDir:
		return "directory"
	default:
		return "unknown"
	}
}

// ResourcePosV1 is the first version from ResourcePos struct
type ResourcePosV1 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// ResourcePos defines resource position in 3D space.
// Alias to the latest ResourcePos version
type ResourcePos ResourcePosV1

// ResourceV1 is the first version from Resource struct
type ResourceV1 struct {
	ResourceID ResourceID   `json:"id"`
	Type       ResourceType `json:"type"`
	Path       string       `json:"path"`
	Pos        ResourcePos  `json:"pos"`
}

// Resource is alias to the latest Resource version
type Resource ResourceV1

func (r *Resource) Open() {
	log.Info("OPEN RESOURCE")
	fileapp.Open(r.Path)
}
