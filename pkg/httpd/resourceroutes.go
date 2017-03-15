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
	r.POST(resourcesURL, CreateResources)
	r.PUT(resourcesURL, UpdateResources)

	resourceURL := resourcesURL + "/:rid"
	r.GET(resourceURL, ReadResource)
	// DELETE with JSON request body is problematic,
	// using POST for multi-delete
	r.POST(resourceURL, DeleteResources)
	r.DELETE(resourceURL, DeleteResource)
}

// CreateResource creates new Resources.
func CreateResources(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pm := findProxyMap(ps.ByName("mapid"))
	if pm == nil {
		WriteJSONError(w, 404, "map not found")
		return
	}

	type Item struct {
		Path string `json:"path"`
	}
	type JSONRequest struct {
		Items []Item `json:"items"`
	}
	var jr JSONRequest
	err := json.NewDecoder(r.Body).Decode(&jr)
	r.Body.Close()
	if err != nil {
		WriteJSONError(w, 400, "bad request")
		return
	}

	log.WithFields(log.Fields{
		"items": jr.Items,
	}).Info("Create Resources")

	var ids []int
	for _, item := range jr.Items {
		rsrc := model.Resource{
			Type: model.ResourceFile,
			Path: item.Path,
			Pos: model.ResourcePos{
				X: 0,
				Y: 0,
				Z: 0,
			},
		}

		rID := pm.AddResource(&rsrc)
		ids = append(ids, rID)
	}
	pm.Write()

	writeResources(w, pm.Map, ids)
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
func UpdateResources(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pm := findProxyMap(ps.ByName("mapid"))
	if pm == nil {
		WriteJSONError(w, 404, "map not found")
		return
	}

	type ResourceData struct {
		ID  int               `json:"id"`
		Pos model.ResourcePos `json:"pos"`
	}
	type JSONRequest struct {
		Resources []ResourceData `json:"resources"`
	}

	var jr JSONRequest
	err := json.NewDecoder(r.Body).Decode(&jr)
	r.Body.Close()
	if err != nil {
		WriteJSONError(w, 400, "bad request")
		return
	}

	pm.Read()
	var ids []int
	for _, resData := range jr.Resources {
		rsrc := pm.Resources[resData.ID]
		if rsrc == nil {
			WriteJSONError(w, 404, fmt.Sprintf("resource %d not found", resData.ID))
			return
		}
		rsrc.Pos = resData.Pos
		ids = append(ids, resData.ID)
	}
	pm.Write()
	writeResources(w, pm.Map, ids)
}

// DeleteResources is controller for deleting multiple resources.
func DeleteResources(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pm := findProxyMap(ps.ByName("mapid"))
	if pm == nil {
		WriteJSONError(w, 404, "map not found")
		return
	}

	if ps.ByName("rid") != "delete" {
		WriteJSONError(w, 400, "bad request")
		return
	}

	type JSONRequest struct {
		IDs []int `json:"ids"`
	}
	var jr JSONRequest
	err := json.NewDecoder(r.Body).Decode(&jr)
	r.Body.Close()
	if err != nil {
		WriteJSONError(w, 400, "bad request")
		return
	}

	for _, id := range jr.IDs {
		pm.DeleteResource(id)
	}

	pm.Write()
	fmt.Fprint(w, "{}")
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

// ResourceResponse is struct used for JSON response.
type ResourcesResponse struct {
	Resources []ResourceResponse `json:"resources"`
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

func writeResources(w http.ResponseWriter, m *model.Map, ids []int) {
	resp := ResourcesResponse{}
	for _, id := range ids {
		rsrc := m.Resources[id]
		if rsrc != nil {
			rr := ResourceResponse{
				ID:       id,
				Resource: rsrc,
			}
			resp.Resources = append(resp.Resources, rr)
		} else {
			WriteJSONError(w, 404, "resource not found")
		}
	}
	WriteJSON(w, resp)
}
