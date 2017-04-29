// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"os"
	"testing"
)

func TestAPIKeyManager(t *testing.T) {
	CreateAPIKeyManager()
	m := GetAPIKeyManager()
	path := getTestPath()
	err := m.readFile(path)
	if err != nil {
		t.Error("Error in readFile", err)
	}

	ak := m.CreateAPIKey()
	m.writeFile(path)

	err = m.readFile(path)
	if err != nil {
		t.Error("Error in readFile", err)
	}

	if m.IsValidAPIKey(ak) == false {
		t.Error("Expected created API key to be valid")
	}

	m.DeleteAPIKey(ak)
	if m.IsValidAPIKey(ak) == true {
		t.Error("Expected API key to be deleted")
	}

	// remove json file after test
	defer os.Remove(path)
}

func getTestPath() string {
	return "testdata/apikeys.json"
}
