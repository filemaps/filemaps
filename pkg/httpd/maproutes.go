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
	"strconv"
	"time"

	"github.com/filemaps/filemaps-backend/pkg/database"
	"github.com/filemaps/filemaps-backend/pkg/model"
)

// GetMaps is controller for getting maps.
// Returns all maps, sorted by Opened field.
func GetMaps(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mm := model.GetMapManager()
	resp := make(map[string]interface{})
	resp["maps"] = mm.GetMaps()
	WriteJSON(w, resp)
}

// CreateMap creates new Map.
func CreateMap(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type JSONRequest struct {
		Title string `json:"title"`
		Base  string `json:"base"`
		File  string `json:"file"`
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
		"title": jr.Title,
		"base":  jr.Base,
		"file":  jr.File,
	}).Info("Create Map")

	fm := database.FileMap{
		Title:  jr.Title,
		Base:   jr.Base,
		File:   jr.File,
		Opened: time.Now(),
	}
	mm := model.GetMapManager()
	mm.AddMap(&fm)
	writeMap(w, fm.ID)
}

// GetMap is controller for getting a map.
func GetMap(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("mapid"))
	if err != nil {
		WriteJSONError(w, 400, "map id must be integer")
		return
	}
	writeMap(w, id)
}

func writeMap(w http.ResponseWriter, id int) {
	mm := model.GetMapManager()
	m := mm.GetMap(id)
	if m != nil {
		WriteJSON(w, m)
	} else {
		WriteJSONError(w, 404, "map not found")
	}
}
