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
)

func route(r *httprouter.Router) {
	r.GET("/", Index)
	r.GET("/hello/:name", Hello)
	r.ServeFiles("/gl/*filepath", http.Dir("filemaps-webui/build"))
	routeMaps(r)
}

// Index is controller for root URL
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("root request")
	fmt.Fprint(w, "Welcome!\n")
}

// Hello is controller for hello
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Info("hello request")
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
