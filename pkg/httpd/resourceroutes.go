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

	"github.com/filemaps/filemaps-backend/pkg/model"
)

func routeResources(r *httprouter.Router, mapURL string) {
	resourcesURL := mapURL + "/resources"
	r.POST(resourcesURL, CreateResource)

	resourceURL := resourcesURL + "/:rid"
	r.GET(resourceURL, ReadResource)
	r.PUT(resourceURL, UpdateResource)
	r.DELETE(resourceURL, DeleteResource)
}

// CreateResource creates new Resource.
func CreateResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pm := findProxyMap(ps.ByName("mapid"))
	if pm == nil {
		WriteJSONError(w, 404, "map not found")
		return
	}

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
	}).Info("Create Resource")

	rsrc := model.Resource{
		Type: model.ResourceFile,
		Path: jr.Path,
		Pos: model.ResourcePos{
			X: 0,
			Y: 0,
			Z: 0,
		},
	}

	rID := pm.AddResource(&rsrc)
	pm.Write()

	writeResource(w, pm.Map, rID)
}

// ReadResource is controller for getting a resource.
func ReadResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pm := findProxyMap(ps.ByName("mapid"))
	if pm == nil {
		WriteJSONError(w, 404, "map not found")
		return
	}

	id, err := strconv.Atoi(ps.ByName("rid"))
	if err != nil {
		WriteJSONError(w, 404, "resource not found")
		return
	}

	pm.Read()
	writeResource(w, pm.Map, id)
}

// UpdateResource updates existing resource.
func UpdateResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pm := findProxyMap(ps.ByName("mapid"))
	if pm == nil {
		WriteJSONError(w, 404, "map not found")
		return
	}

	id, err := strconv.Atoi(ps.ByName("rid"))
	if err != nil {
		WriteJSONError(w, 404, "resource not found")
		return
	}

	type JSONRequest struct {
		Pos model.ResourcePos `json:"pos"`
	}
	var jr JSONRequest
	d := json.NewDecoder(r.Body)
	err = d.Decode(&jr)
	r.Body.Close()
	if err != nil {
		WriteJSONError(w, 400, "bad request")
		return
	}

	pm.Read()
	rsrc := pm.Resources[id]
	if rsrc == nil {
		WriteJSONError(w, 404, "resource not found")
		return
	}

	rsrc.Pos = jr.Pos
	pm.Write()
	writeResource(w, pm.Map, id)
}

// DeleteResource is controller for deleting a resource.
func DeleteResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pm := findProxyMap(ps.ByName("mapid"))
	if pm == nil {
		WriteJSONError(w, 404, "map not found")
		return
	}

	id, err := strconv.Atoi(ps.ByName("rid"))
	if err != nil {
		WriteJSONError(w, 404, "resource not found")
		return
	}

	pm.DeleteResource(id)
	pm.Write()

	fmt.Fprint(w, "{}")
}

// ResourceResponse is struct used for JSON response.
type ResourceResponse struct {
	ID int `json:"id"`
	*model.Resource
}

func writeResource(w http.ResponseWriter, m *model.Map, id int) {
	rsrc := m.Resources[id]
	if rsrc != nil {
		rr := ResourceResponse{
			ID:       id,
			Resource: rsrc,
		}
		WriteJSON(w, rr)
	} else {
		WriteJSONError(w, 404, "resource not found")
	}
}
