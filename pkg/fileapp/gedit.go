// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fileapp

import (
	log "github.com/Sirupsen/logrus"
	"os/exec"
)

const (
	geditCmd = "gedit"
)

func init() {
	// register file app if command exists
	_, err := exec.Command("which", geditCmd).Output()
	if err == nil {
		register(NewGedit())
	}
}

// Gedit implements FileApp interface
type Gedit struct {
}

func NewGedit() *Gedit {
	return &Gedit{}
}

func (a *Gedit) getInfo() FileAppInfo {
	return FileAppInfo{
		ID:   "gedit",
		Name: "gedit",
	}
}

func (a *Gedit) open(path string) int {
	log.WithFields(log.Fields{
		"path": path,
	}).Info("Gedit: open")

	out, err := exec.Command(geditCmd, path).Output()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("gedit open error")
		return -1
	}
	log.WithFields(log.Fields{
		"out": out,
	}).Info("gedit")
	return 0
}
