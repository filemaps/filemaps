// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build !windows

package main

import (
	"os/exec"
	"runtime"
	"syscall"
)

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		// OSX
		return exec.Command("open", url).Run()

	default:
		cmd := exec.Command("xdg-open", url)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}
		return cmd.Run()
	}
}
