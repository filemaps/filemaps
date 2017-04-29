// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"math/rand"
	"time"
)

const (
	keyBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// APIKey is a database struct for APIKey.
type APIKey struct {
	APIKey  string    `json:"apikey"`
	Expires time.Time `json:"expires"`
}

// NewAPIKey returns new APIKey with random API key for one year.
func NewAPIKey() APIKey {
	return APIKey{
		APIKey:  randString(32),
		Expires: time.Now().AddDate(1, 0, 0),
	}
}

// randString generates a random string with given length.
func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = keyBytes[rand.Intn(len(keyBytes))]
	}
	return string(b)
}
