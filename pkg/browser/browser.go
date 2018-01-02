// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package browser

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
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

// Sorting technique using programmable sort criteria

// By is the type of a "less" function that defines the ordering of its Item arguments.
type By func(i1, i2 *Item) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(items []Item) {
	is := &itemSorter{
		items: items,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(is)
}

// itemSorter joins a By function and a slice of Items to be sorted.
type itemSorter struct {
	items []Item
	by    func(i1, i2 *Item) bool
}

// Len is part of sort.Interface.
func (s *itemSorter) Len() int {
	return len(s.items)
}

// Swao is part of sort.Interface.
func (s *itemSorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *itemSorter) Less(i, j int) bool {
	return s.by(&s.items[i], &s.items[j])
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

	// sort d.Contents by name
	name := func(i1, i2 *Item) bool {
		return strings.ToLower(i1.Name) < strings.ToLower(i2.Name)
	}
	By(name).Sort(d.Contents)

	return d, nil
}
