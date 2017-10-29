// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package config

import (
	"testing"
)

func TestReadWrite(t *testing.T) {
	cfg, err := readFile(getTestPath())
	if err != nil {
		t.Error("Error in GetOrCreate", err)
	} else {
		err = cfg.Write()
		if err != nil {
			t.Error("Error in write", err)
		}
	}
}

func getTestPath() string {
	return "/tmp/filemaps.config"
}
