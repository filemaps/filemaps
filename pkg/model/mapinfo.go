// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"time"
)

// MapInfo is info struct for FileMap.
type MapInfo struct {
	ID     int       `json:"id"`
	Title  string    `json:"title"`
	Base   string    `json:"base"`
	File   string    `json:"file"`
	Opened time.Time `json:"opened"`
}

// MapInfos is a collection of MapInfo pointers.
type MapInfos []MapInfo

// Implementation of sort.Interface for MapInfos.
func (slice MapInfos) Len() int {
	return len(slice)
}

func (slice MapInfos) Less(i, j int) bool {
	return slice[i].Opened.Before(slice[j].Opened)
}

func (slice MapInfos) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
