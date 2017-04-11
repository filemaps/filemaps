// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package httpd

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os/user"

	"github.com/filemaps/filemaps-backend/pkg/filemaps"
)

const (
	// APIURL is prefix for REST API URL.
	APIURL = "/api"
)

func route(r *httprouter.Router) {
	r.GET("/", Index)
	r.GET(APIURL+"/info", Info)
	staticPath := "filemaps-webui"
	r.ServeFiles("/gl/*filepath", http.Dir(staticPath))
	routeMaps(r)
	routeBrowse(r)
}

// Index is controller for root URL
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("root request")
	fmt.Fprint(w, "Welcome!\n")
}

// Info is controller for information
func Info(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := make(map[string]interface{})
	resp["version"] = filemaps.Version
	usr, _ := user.Current()
	if usr != nil {
		resp["homeDir"] = usr.HomeDir
	}

	WriteJSON(w, resp)
}
