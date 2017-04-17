// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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
