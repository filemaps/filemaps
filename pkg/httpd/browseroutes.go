// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package httpd

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"net/http"

	"github.com/filemaps/filemaps/pkg/browser"
)

func routeBrowse(r *httprouter.Router) {
	browseURL := APIURL + "/browse"
	r.POST(browseURL, Browse)
}

// Browse is controller for file browse.
func Browse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type JSONRequest struct {
		Path string `json:"path"`
	}
	var jr JSONRequest
	d := json.NewDecoder(r.Body)
	err := d.Decode(&jr)
	r.Body.Close()
	if err != nil {
		WriteJSONError(w, 400, "bad request")
		return
	}

	log.WithFields(log.Fields{
		"path": jr.Path,
	}).Info("Browse")

	if jr.Path == "" {
		WriteJSONError(w, 400, "empty path")
		return
	}

	dir, err := browser.ScanDir(jr.Path)
	if err != nil {
		log.WithFields(log.Fields{
			"path": jr.Path,
			"err":  err,
		}).Error("Could not scan dir")
		WriteJSONError(w, 500, "browse failed")
		return
	}
	WriteJSON(w, dir)
}
