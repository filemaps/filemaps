// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fileapp

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
)

const (
	geditBin = "/usr/bin/gedit"
)

func init() {
	// register file app if binary exists
	if _, err := os.Stat(geditBin); err == nil {
		register(NewGedit())
	}
}

type Gedit struct {
}

func NewGedit() *Gedit {
	return &Gedit{}
}

func (a *Gedit) getName() string {
	return "gedit"
}

func (a *Gedit) open(path string) int {
	log.WithFields(log.Fields{
		"path": path,
	}).Info("Gedit: open")

	out, err := exec.Command(geditBin, path).Output()
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
