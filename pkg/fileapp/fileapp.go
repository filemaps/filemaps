// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

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
