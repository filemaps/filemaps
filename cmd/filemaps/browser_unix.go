// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

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
