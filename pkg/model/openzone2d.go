// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package model

// ResourceType defines type of Resource
type OpenZoneType int

// Resource type enum
const (
	OpenUp OpenZoneType = iota
	OpenRight
	OpenDown
	OpenLeft
)

// Converts ResourceType to string
func (o OpenZoneType) String() string {
	switch o {
	case OpenUp:
		return "up"
	case OpenRight:
		return "right"
	case OpenDown:
		return "down"
	case OpenLeft:
		return "left"
	default:
		return "unknown"
	}
}

// Zone2DV1 is the first version from Zone2D struct
type OpenZone2DV1 struct {
	Zone2DV1
	Type  OpenZoneType `json:"type"`
	Pos   Position2D   `json:"pos"`
	Width float64      `json:"width"`
}

// Resource is alias to the latest Resource version
type OpenZone2D OpenZone2DV1

// NewZone creates a new Zone
func NewOpenZone2D() *OpenZone2D {
	z := &OpenZone2D{
		Zone2DV1: Zone2DV1{
			Path: make([]Position, 0),
		},
		Pos: Position2D{X: 0, Y: 0},
	}
	return z
}

func NewNewZone2D() *OpenZone2D {
	z := NewOpenZone2D()
	z.ZoneID = 1
	z.Label = "New"
	z.Pos = Position2D{X: 0, Y: 0}
	z.Width = 500
	z.UIClass = "new"
	return z
}

func (z *OpenZone2D) posIsIn(pos Position) bool {
	switch z.Type {
	case OpenUp:
		return (pos.X >= z.Pos.X && pos.X <= z.Pos.X+z.Width && pos.Y >= z.Pos.Y)
	case OpenRight:
		return (pos.X >= z.Pos.X && pos.Y >= z.Pos.Y && pos.Y <= z.Pos.Y+z.Width)
	case OpenDown:
		return (pos.X >= z.Pos.X && pos.X <= z.Pos.X+z.Width && pos.Y <= z.Pos.Y)
	case OpenLeft:
		return (pos.X <= z.Pos.X && pos.Y >= z.Pos.Y && pos.Y <= z.Pos.Y+z.Width)
	}
	return false
}
