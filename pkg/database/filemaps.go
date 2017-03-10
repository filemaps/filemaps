// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package database

// FileMaps is a collection of FileMap pointers.
type FileMaps []*FileMap

// Implementation of sort.Interface for FileMaps.
func (slice FileMaps) Len() int {
	return len(slice)
}

func (slice FileMaps) Less(i, j int) bool {
	return slice[i].Opened.Before(slice[j].Opened)
}

func (slice FileMaps) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
