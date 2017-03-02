// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.Info("Building and installing File Maps")
	target := "github.com/filemaps/filemaps-backend/cmd/filemaps"
	run("go", "install", target)
}

func run(cmd string, args ...string) {
	cmdh := exec.Command(cmd, args...)
	cmdh.Stdout = os.Stdout
	cmdh.Stderr = os.Stderr
	err := cmdh.Run()
	if err != nil {
		log.WithFields(log.Fields{
			"cmd":  cmd,
			"args": args,
		}).Error(err)
	}
}
