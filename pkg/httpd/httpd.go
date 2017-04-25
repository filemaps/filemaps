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
	"net"
	"net/http"
	"strings"

	"github.com/filemaps/filemaps-backend/pkg/config"
	"github.com/filemaps/filemaps-backend/pkg/model"
)

// RunHTTP starts HTTP server
func RunHTTP(addr string, webUIPath string) {
	router := httprouter.New()
	route(router, webUIPath)
	handler := authMiddleware(router)
	log.WithFields(log.Fields{
		"transport": "HTTP",
		"addr":      addr,
	}).Info("Starting server")
	log.Fatal(http.ListenAndServe(addr, handler))
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

// authMiddleware authenticates the request.
// Request must come from trusted address or X-API-Key header must
// contain a valid API key.
func authMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		model.GetAPIKeyManager()
		if addrIsTrusted(r.RemoteAddr) {
			handler.ServeHTTP(w, r)
			return
		}

		if model.GetAPIKeyManager().IsValidAPIKey(r.Header.Get("X-API-Key")) {
			handler.ServeHTTP(w, r)
			return
		}

		log.WithFields(log.Fields{
			"requestURI": r.RequestURI,
			"remoteAddr": r.RemoteAddr,
		}).Error("Access denied")
		w.WriteHeader(403)
	})
}

// addrIsTrusted returns true if given address is trusted.
// addr is request.RemoteAddr which has format IP:port
func addrIsTrusted(addr string) bool {
	// strip port
	addr = addr[:strings.LastIndex(addr, ":")]
	remoteIP := net.ParseIP(addr)
	if remoteIP == nil {
		log.WithFields(log.Fields{
			"ip": addr,
		}).Error("Could not parse remote IP")
		return false
	}

	// trust loopback addresses
	if remoteIP.IsLoopback() {
		return true
	}

	// check trusted addresses from config
	cfg := config.GetConfiguration()
	addrs := strings.Split(cfg.TrustedAddresses, ",")
	for _, a := range addrs {
		trustedIP := net.ParseIP(a)
		if trustedIP != nil && trustedIP.Equal(remoteIP) {
			return true
		}
	}

	return false
}
