// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package model

import (
	"math/rand"
	"time"
)

const (
	keyBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// APIKey is a struct for API key.
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
