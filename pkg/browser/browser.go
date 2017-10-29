// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package browser

import (
	"io/ioutil"
	"path/filepath"
)

// ItemType defines type of Item.
type ItemType int

const (
	// ItemFile represents regular file item
	ItemFile ItemType = iota
	// ItemDir represents directory item
	ItemDir
)

// Item is file or directory struct.
type Item struct {
	Name string   `json:"name"`
	Path string   `json:"path"`
	Size int      `json:"size"`
	Type ItemType `json:"type"`
}

// Dir is directory struct.
type Dir struct {
	Path     string `json:"path"`
	Parent   string `json:"parent"`
	Contents []Item `json:"contents"`
}

// ScanDir scans given directory and returns Dir struct.
func ScanDir(path string) (Dir, error) {
	d := Dir{
		Path:   path,
		Parent: filepath.Dir(path),
	}

	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return d, err
	}

	for _, info := range infos {
		t := ItemFile
		if info.IsDir() {
			t = ItemDir
		}
		i := Item{
			Name: info.Name(),
			Path: filepath.Join(path, info.Name()),
			Size: int(info.Size()),
			Type: t,
		}
		d.Contents = append(d.Contents, i)
	}
	return d, nil
}
