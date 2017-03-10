// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package httpd

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"

	"github.com/filemaps/filemaps-backend/pkg/database"
	"github.com/filemaps/filemaps-backend/pkg/model"
)

const (
	// APIURL is prefix for REST API URL.
	APIURL = "/api"
)

func routeMaps(r *httprouter.Router) {
	mapsURL := APIURL + "/maps"
	r.GET(mapsURL, GetMaps)
	r.POST(mapsURL, CreateMap)

	mapURL := mapsURL + "/:mapid"
	r.GET(mapURL, ReadMap)
	r.PUT(mapURL, UpdateMap)
	r.DELETE(mapURL, DeleteMap)
}

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

// ReadMap is controller for getting a map.
func ReadMap(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("mapid"))
	if err != nil {
		WriteJSONError(w, 400, "map id must be integer")
		return
	}
	writeMap(w, id)
}

// UpdateMap updates existing Map.
func UpdateMap(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("mapid"))
	if err != nil {
		WriteJSONError(w, 400, "map id must be integer")
		return
	}

	type JSONRequest struct {
		Title string `json:"title"`
		Base  string `json:"base"`
		File  string `json:"file"`
	}
	var jr JSONRequest
	d := json.NewDecoder(r.Body)
	err = d.Decode(&jr)
	r.Body.Close()
	if err != nil {
		WriteJSONError(w, 400, "bad request")
		return
	}

	log.WithFields(log.Fields{
		"title": jr.Title,
		"base":  jr.Base,
		"file":  jr.File,
	}).Info("Update Map")

	mm := model.GetMapManager()
	pm := mm.Maps[id]
	if pm == nil {
		WriteJSONError(w, 404, "map not found")
	}
	pm.SetTitle(jr.Title)
	pm.SetBase(jr.Base)
	pm.SetFile(jr.File)
	pm.Write()

	writeMap(w, id)
}

// DeleteMap is controller for deleting a map.
func DeleteMap(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("mapid"))
	if err != nil {
		WriteJSONError(w, 400, "map id must be integer")
		return
	}

	mm := model.GetMapManager()
	if err := mm.DeleteMap(id); err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"id":  id,
		}).Error("Could not remove map")
		WriteJSONError(w, 500, "could not remove map")
		return
	}
	fmt.Fprint(w, "{}")
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
