// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package httpd

import (
	"encoding/json"
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

// WriteJSON writes JSON response
func WriteJSON(w http.ResponseWriter, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
	return nil
}

// WriteJSONError writes error JSON response
func WriteJSONError(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	WriteJSON(w, map[string]string{
		"error": err,
	})
}
