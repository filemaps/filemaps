// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package model

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

// ResourceV1 is the first version from Resource struct
type ResourceV1 struct {
	ResourceID ResourceID   `json:"id"`
	Type       ResourceType `json:"type"`
	Path       string       `json:"path"`
	Pos        Position     `json:"pos"`
	Style      Style        `json:"style"`
}

// Resource is alias to the latest Resource version
type Resource ResourceV1
