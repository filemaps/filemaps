// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

// +build windows

package main

import (
	"os/exec"
)

func openBrowser(url string) error {
	return exec.Command("cmd.exe", "/C", "start "+url).Run()
}
