// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package httpd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// RunHTTP starts HTTP server
func RunHTTP(addr string) {
	router := httprouter.New()
	route(router)
	log.WithFields(log.Fields{
		"transport": "HTTP",
		"addr":      addr,
	}).Info("Starting server")
	log.Fatal(http.ListenAndServe(addr, router))
}
