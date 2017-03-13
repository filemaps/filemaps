// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"time"

	"github.com/filemaps/filemaps-backend/pkg/config"
	"github.com/filemaps/filemaps-backend/pkg/database"
	"github.com/filemaps/filemaps-backend/pkg/httpd"
	"github.com/filemaps/filemaps-backend/pkg/model"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	// show only clock time in log output
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05.000",
	})
}

func main() {
	log.Info("File Maps starting")

	// initialize database
	err := database.InitDatabase()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Database initialization error")
	}

	// create singleton instance from MapManager
	_, err = model.CreateMapManager()
	if err != nil {
		log.Fatal(err)
	}

	// create singleton instance from APIKeyManager
	_, err = model.CreateAPIKeyManager()
	if err != nil {
		log.Fatal(err)
	}

	cfg, _ := config.Read()
	config.Write(cfg)

	addr := ":8338"
	httpd.RunHTTP(addr)
}
