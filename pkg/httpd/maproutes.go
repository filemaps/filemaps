// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package httpd

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"

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

// GetMap is controller for getting a map.
func GetMap(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("mapid"))
	if err != nil {
		WriteJSONError(w, 400, "map id must be integer")
		return
	}

	mm := model.GetMapManager()
	m := mm.GetMap(id)
	WriteJSON(w, m)
}
