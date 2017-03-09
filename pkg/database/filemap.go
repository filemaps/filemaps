// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package database

import (
	"fmt"
	"time"
)

// FileMap is a database struct for FileMap.
type FileMap struct {
	ID     int       `json:"id"`
	Title  string    `json:"title"`
	Base   string    `json:"base"`
	File   string    `json:"file"`
	Opened time.Time `json:"opened"`
}

// CreateTableFileMaps creates filemaps table if it does not exist.
func (db *Database) CreateTableFileMaps() error {
	// sqlite creates rowid column automatically
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS filemaps (
		title text,
		base text,
		file text,
		opened integer
	)
	`
	_, err := db.h.Exec(sqlStmt)
	return err
}

// AddFileMap inserts new FileMap to database and updates
// FileMap.ID with new ID.
func (db *Database) AddFileMap(fm *FileMap) error {
	stmt, err := db.h.Prepare("INSERT INTO filemaps (title, base, file, opened) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	rslt, err := stmt.Exec(fm.Title, fm.Base, fm.File, fm.Opened.Unix())
	if err != nil {
		return err
	}

	id, err := rslt.LastInsertId()
	fm.ID = int(id)

	return err
}

// GetFileMaps returns FileMap rows from database, giving max limit rows,
// last opened first. If limit is < 1, returns all rows.
func (db *Database) GetFileMaps(limit int) ([]FileMap, error) {
	var fileMaps []FileMap
	stmt := "SELECT rowid, title, base, file, opened FROM filemaps ORDER BY opened DESC"
	if limit > 0 {
		stmt += fmt.Sprintf(" LIMIT %d", limit)
	}
	rows, err := db.h.Query(stmt)
	if err != nil {
		return fileMaps, err
	}
	defer rows.Close()

	for rows.Next() {
		fm := FileMap{}
		var id int64
		var opened int64
		err = rows.Scan(&id, &fm.Title, &fm.Base, &fm.File, &opened)
		fm.ID = int(id)
		fm.Opened = time.Unix(opened, 0)
		if err != nil {
			return fileMaps, err
		}
		fileMaps = append(fileMaps, fm)
	}
	err = rows.Err()
	return fileMaps, err
}

// GetFileMap returns a FileMap row by given ID.
func (db *Database) GetFileMap(ID int) (FileMap, error) {
	fm := FileMap{}
	stmt, err := db.h.Prepare("SELECT rowid, title, base, file, opened FROM filemaps WHERE rowid = ?")
	if err != nil {
		return fm, err
	}
	defer stmt.Close()

	var i64 int64
	var opened int64
	err = stmt.QueryRow(ID).Scan(&i64, &fm.Title, &fm.Base, &fm.File, &opened)
	fm.ID = int(i64)
	fm.Opened = time.Unix(opened, 0)
	return fm, err
}

// UpdateFileMap updates given FileMap row in database.
func (db *Database) UpdateFileMap(fm FileMap) error {
	stmt, err := db.h.Prepare("UPDATE filemaps SET title = ?, base = ?, file = ?, opened = ? WHERE rowid = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(fm.Title, fm.Base, fm.File, fm.Opened.Unix(), fm.ID)
	return err
}

// DeleteFileMap deletes FileMap row by given ID.
func (db *Database) DeleteFileMap(ID int) error {
	stmt := fmt.Sprintf("DELETE FROM filemaps WHERE rowid = %d", ID)
	_, err := db.h.Exec(stmt)
	return err
}
