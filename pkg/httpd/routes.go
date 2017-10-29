// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package httpd

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os/user"

	"github.com/filemaps/filemaps/pkg/filemaps"
)

const (
	// APIURL is prefix for REST API URL.
	APIURL = "/api"
	// UIURL is prefix for Web UI URL.
	UIURL = "/ui"
)

func route(r *httprouter.Router, webUIPath string) {
	r.GET("/", Index)
	r.GET(APIURL+"/info", Info)

	routeMaps(r)
	routeBrowse(r)
	routeConfig(r)
	routeWebUI(r, webUIPath)
}

// Index is controller for root URL
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, UIURL+"/", http.StatusTemporaryRedirect)
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
