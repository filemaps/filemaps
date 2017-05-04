// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fileapp

import (
	log "github.com/Sirupsen/logrus"

	"github.com/filemaps/filemaps/pkg/config"
)

type FileAppInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FileApp interface {
	getInfo() FileAppInfo
	open(path string) int
}

var (
	apps []FileApp
)

func GetInfos() []FileAppInfo {
	var infos []FileAppInfo
	for i := 0; i < len(apps); i++ {
		info := apps[i].getInfo()
		if info.ID != "custom1" {
			infos = append(infos, info)
		}
	}
	return infos
}

func Open(path string) int {
	app := routeApp(path)
	if app == nil {
		return -1
	}
	return app.open(path)
}

func register(app FileApp) {
	log.WithFields(log.Fields{
		"app": app.getInfo().Name,
	}).Info("register file app")

	apps = append(apps, app)
}

func routeApp(path string) FileApp {
	cfg := config.GetConfiguration()
	for _, app := range apps {
		if app.getInfo().ID == cfg.TextEditor {
			return app
		}
	}

	// text editor not configured or not available
	for _, app := range apps {
		if app.getInfo().ID != "custom1" {
			return app
		}
	}
	return apps[0]
}
