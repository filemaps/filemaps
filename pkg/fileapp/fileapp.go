// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fileapp

import (
	log "github.com/Sirupsen/logrus"
)

type FileApp interface {
	getName() string
	open(path string) int
}

var (
	apps []FileApp
)

func Open(path string) int {
	app := routeApp(path)
	if app == nil {
		return -1
	}
	return app.open(path)
}

func register(app FileApp) {
	log.WithFields(log.Fields{
		"app": app.getName(),
	}).Info("register file app")
	apps = append(apps, app)
}

func routeApp(path string) FileApp {
	return apps[0]
}
