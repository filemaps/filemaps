// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package fileapp

import (
	log "github.com/Sirupsen/logrus"
	"os/exec"
	"strings"

	"github.com/filemaps/filemaps/pkg/config"
)

const (
	// argSeparator separates arguments from custom command string
	argSeparator = ","
)

func init() {
	register(NewCustom())
}

type Custom struct {
}

func NewCustom() *Custom {
	return &Custom{}
}

func (a *Custom) getInfo() FileAppInfo {
	return FileAppInfo{
		ID:   "custom1",
		Name: "Custom",
	}
}

func (a *Custom) open(path string) int {
	log.WithFields(log.Fields{
		"path": path,
	}).Info("Custom: open")

	cfg := config.GetConfiguration()
	cmd := strings.Split(cfg.TextEditorCustom1Cmd, argSeparator)
	cmd = append(cmd, path)

	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Custom open error")
		return -1
	}
	log.WithFields(log.Fields{
		"out": out,
	}).Info("Custom")
	return 0
}
