// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package database

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"

	"github.com/filemaps/filemaps-backend/pkg/config"
)

// Database struct
type Database struct {
	h *sql.DB // database handle
}

// NewDB returns new Database struct.
func NewDB() Database {
	return Database{}
}

// Open opens database connection.
// Remember to call Close().
func (db *Database) Open() error {
	path := filepath.Join(config.GetDir(), "filemaps.db")
	return db.openFile(path)
}

func (db *Database) openFile(path string) error {
	var err error
	log.WithFields(log.Fields{
		"type": "sqlite3",
		"path": path,
	}).Info("Database")
	db.h, err = sql.Open("sqlite3", path)
	return err
}

// Init makes sure database is up-to-date.
func (db *Database) Init() error {
	if err := db.CreateTableFileMaps(); err != nil {
		return err
	}
	if err := db.CreateTableMigrations(); err != nil {
		return err
	}
	if err := db.RunMigrations(); err != nil {
		return err
	}
	return nil
}

// Closes active database connection.
func (db *Database) Close() {
	if db.h != nil {
		db.h.Close()
	}
}
