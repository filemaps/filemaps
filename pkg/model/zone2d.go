// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package model

// Zone2DV1 is the first version from Zone2D struct
type Zone2DV1 struct {
	ZoneID  int               `json:"id"`
	Label   string            `json:"label"`
	Path    []Position        `json:"path"`
	UIClass string            `json:"uiClass"`
	Style   map[string]string `json:"style"`
}

// Zone2D is alias to the latest Zone2D version
type Zone2D Zone2DV1

// NewZone2D creates a new Zone2D
func NewZone2D() *Zone2D {
	z := &Zone2D{
		Path: make([]Position, 0),
	}
	return z
}
