// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	colorable "github.com/mattn/go-colorable"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/filemaps/filemaps/pkg/config"
	"github.com/filemaps/filemaps/pkg/httpd"
	"github.com/filemaps/filemaps/pkg/model"
)

var (
	noBrowser bool
	port      int
	webUIPath string
)

func init() {
	rand.Seed(time.Now().UnixNano())

	// show only clock time in log output
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05.000",
	})

	// colors for windows
	if runtime.GOOS == "windows" {
		log.SetOutput(colorable.NewColorableStdout())
	}

	flag.BoolVar(&noBrowser, "no-browser", false, "Do not open browser")
	flag.BoolVar(&httpd.CORSEnabled, "cors", false, "Enable CORS")
	flag.IntVar(&port, "port", 8338, "Port to listen to")
	flag.StringVar(&webUIPath, "webui", "", "Path of Web UI files")
}

func main() {
	log.Info("File Maps starting")

	flag.Parse()

	if err := config.EnsureDir(); err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"path": config.GetDir(),
		}).Fatal("Config directory could not be created")
	}

	// create singleton instance from MapManager
	_, err := model.CreateMapManager()
	if err != nil {
		log.Fatal(err)
	}

	// create singleton instance from APIKeyManager
	_, err = model.CreateAPIKeyManager()
	if err != nil {
		log.Fatal(err)
	}

	config.CreateConfiguration()

	addr := ":" + strconv.Itoa(port)

	if noBrowser == false {
		openBrowser("http://localhost" + addr + httpd.UIURL)
	}
	httpd.RunHTTP(addr, webUIPath)
}
