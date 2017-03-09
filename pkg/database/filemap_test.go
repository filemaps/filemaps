// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package database

import (
	"os"
	"testing"
	"time"
)

func TestFileMap(t *testing.T) {
	db := NewDB()
	path := "/tmp/filemaps-test.db"
	if err := db.openFile(path); err != nil {
		t.Error("Error in openFile", err)
		return
	}
	defer db.Close()
	if err := db.Init(); err != nil {
		t.Error("Error in Init", err)
		return
	}
	// clean database file after test
	defer os.Remove(path)

	// add FileMap
	fm := FileMap{
		Title:  "test1",
		Base:   "/tmp",
		File:   "project.filemap",
		Opened: time.Now(),
	}
	if err := db.AddFileMap(&fm); err != nil {
		t.Error("Error in AddFileMap", err)
		return
	}

	// update FileMap
	newTitle := "test2"
	fm.Title = newTitle
	fm.Base = "/tmp2"
	now := time.Now()
	fm.Opened = now
	if err := db.UpdateFileMap(fm); err != nil {
		t.Error("Error in UpdateFileMap", err)
		return
	}

	// read all FileMaps
	fileMaps, err := db.GetFileMaps(0)
	if err != nil {
		t.Error("Error in GetFileMaps", err)
		return
	}
	if len(fileMaps) != 1 {
		t.Error("expected", 1, "got", len(fileMaps))
		return
	}

	// read one FileMap
	fileMap, err := db.GetFileMap(fileMaps[0].ID)
	if err != nil {
		t.Error("Error in GetFileMap", err)
		return
	}
	if fileMap.Title != newTitle {
		t.Error("name expected", newTitle, "got", fileMap.Title)
	}

	// delete FileMap
	if err = db.DeleteFileMap(fileMap.ID); err != nil {
		t.Error("Error in DeleteFileMap", err)
	}

	fileMaps, err = db.GetFileMaps(0)
	if err != nil {
		t.Error("Error in GetFileMaps", err)
		return
	}
	if len(fileMaps) != 0 {
		t.Error("expected", 0, "got", len(fileMaps))
		return
	}
}
