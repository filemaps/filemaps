// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/filemaps/filemaps-backend/pkg/config"
	"github.com/filemaps/filemaps-backend/pkg/httpd"
)

func init() {
	// show only clock time in log output
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05.000",
	})
}

func main() {
	log.Info("File Maps starting")
	cfg, _ := config.GetOrCreate()
	cfg.Version += 1
	config.Write(cfg)
	addr := ":8080"
	httpd.RunHTTP(addr)
}
