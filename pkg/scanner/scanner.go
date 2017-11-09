// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package scanner

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func Scan(path string, base string, exclude []string) []string {
	log.WithFields(log.Fields{
		"path":    path,
		"base":    base,
		"exclude": exclude,
	}).Info("Start")

	files := readDir(path, base, exclude)
	log.WithFields(log.Fields{
		"files": files,
	}).Info("Files found by scanning")
	return files
}

func readDir(path string, base string, exclude []string) []string {
	var found []string

	files, err := ioutil.ReadDir(path)
	if err != nil {

		log.WithFields(log.Fields{
			"path": path,
			"err":  err,
		}).Error("Error when reading dir")
	}

	var dirs []os.FileInfo
	for _, file := range files {
		filePath := filepath.Join(path, file.Name())

		if isExcluded(filePath, file.IsDir(), base, exclude) {
			continue
		}
		if file.IsDir() {
			dirs = append(dirs, file)
		} else {
			log.Info(filePath)
			found = append(found, filePath)
		}
	}

	for _, dir := range dirs {
		subFound := readDir(filepath.Join(path, dir.Name()), base, exclude)
		found = append(found, subFound...)
	}
	return found
}

func isExcluded(path string, isDir bool, base string, exclude []string) bool {
	relative, err := filepath.Rel(base, path)
	if err != nil {
		log.WithFields(log.Fields{
			"base": base,
			"path": path,
		}).Error("Could not get relative path")
		return false
	}
	for i := len(exclude) - 1; i >= 0; i-- {
		pattern := strings.Trim(exclude[i], " ")
		if len(pattern) == 0 {
			// empty line
			continue
		}

		if pattern[0] == '#' {
			// ignore comments
			continue
		}

		// patterns starting with '!' are inverted
		invert := false
		if pattern[0] == '!' {
			invert = true
			pattern = pattern[1:]
		}

		// patterns starting with '/'
		if pattern[0] == '/' {
			pattern = pattern[1:]
		}

		// patterns ending with '/' must match to directory
		mustBeDir := false
		if pattern[len(pattern)-1] == '/' {
			mustBeDir = true
			pattern = pattern[0 : len(pattern)-1]
		}

		name := relative
		// if patterns does not contain "/", compare with base name
		if !strings.ContainsAny(pattern, "/") {
			name = filepath.Base(relative)
		}

		matched, err := filepath.Match(pattern, name)
		if err != nil {
			log.WithFields(log.Fields{
				"pattern": pattern,
				"name":    name,
			}).Error("Invalid exclude pattern")
			continue
		}
		if matched && (!mustBeDir || isDir) {
			return !invert
		}
	}
	return false
}
