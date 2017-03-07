// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package database

import (
	"fmt"
)

const (
	// CurrentVersion defines current db migration level
	CurrentVersion = 1
)

// CreateTableMigrations creates migrations table if it does not exist.
func (db *Database) CreateTableMigrations() error {
	stmt := `
	CREATE TABLE migrations (
		version integer
	)
	`
	_, err := db.h.Exec(stmt)
	if err == nil {
		stmt = fmt.Sprintf("INSERT INTO migrations (version) VALUES (%d)", CurrentVersion)
		_, err = db.h.Exec(stmt)
		return err
	}
	return nil
}

// RunMigrations runs needed database migrations.
func (db *Database) RunMigrations() error {
	var version int
	err := db.h.QueryRow("SELECT version FROM migrations").Scan(&version)
	return err
}
