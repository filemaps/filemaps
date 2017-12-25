// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package model

// Style defines graphical style rules.
type Style struct {
	SClass string            `json:"sClass"`
	Rules  map[string]string `json:"rules"`
}
