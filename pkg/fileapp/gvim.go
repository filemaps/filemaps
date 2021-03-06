// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package fileapp

import (
	log "github.com/Sirupsen/logrus"
	"os/exec"
)

const (
	gvimCmd = "gvim"
)

func init() {
	// register file app if command exists
	_, err := exec.Command("which", gvimCmd).Output()
	if err == nil {
		register(NewGVim())
	}
}

type GVim struct {
}

func NewGVim() *GVim {
	return &GVim{}
}

func (a *GVim) getInfo() FileAppInfo {
	return FileAppInfo{
		ID:   "gvim",
		Name: "GVim",
	}
}

func (a *GVim) open(path string) int {
	log.WithFields(log.Fields{
		"path": path,
	}).Info("GVim: open")

	out, err := exec.Command(gvimCmd, "-p", "--remote-tab-silent", path).Output()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("GVim open error")
		return -1
	}
	log.WithFields(log.Fields{
		"out": out,
	}).Info("GVim")
	return 0
}
