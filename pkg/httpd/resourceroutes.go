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
	"path/filepath"
	"strconv"

	"github.com/filemaps/filemaps/pkg/model"
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

	r.GET(resourceURL+"/open", OpenResource)
}

// CreateResource creates new Resources.
func CreateResources(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pm := findProxyMap(ps.ByName("mapid"))
	if pm == nil {
		WriteJSONError(w, 404, "map not found")
		return
	}

	type Item struct {
		Path string            `json:"path"`
		Pos  model.ResourcePos `json:"pos"`
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

	var ids []model.ResourceID
	for _, item := range jr.Items {
		// convert absolute path to relative
		path, err := filepath.Rel(pm.Base, item.Path)
		if err != nil {
			log.WithFields(log.Fields{
				"basepath": pm.Base,
				"targpath": item.Path,
			}).Error("Could not make relative path")
			path = item.Path
		}
		rsrc := model.Resource{
			Type: model.ResourceFile,
			Path: path,
			Pos:  item.Pos,
		}

		rID := pm.AddResource(&rsrc)
		ids = append(ids, rID)
	}
	pm.Write()

	writeResources(w, pm, ids)
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
	writeResource(w, pm, model.ResourceID(id))
}

// UpdateResource updates existing resource.
func UpdateResources(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pm := findProxyMap(ps.ByName("mapid"))
	if pm == nil {
		WriteJSONError(w, 404, "map not found")
		return
	}

	type ResourceData struct {
		ID  model.ResourceID  `json:"id"`
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
	var ids []model.ResourceID
	for _, resData := range jr.Resources {
		rsrc := pm.GetResource(resData.ID)
		if rsrc == nil {
			WriteJSONError(w, 404, fmt.Sprintf("resource %d not found", resData.ID))
			return
		}
		rsrc.Pos = resData.Pos
		ids = append(ids, rsrc.ResourceID)
	}
	pm.Write()

	writeResources(w, pm, ids)
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
		pm.DeleteResource(model.ResourceID(id))
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

	pm.DeleteResource(model.ResourceID(id))
	pm.Write()

	fmt.Fprint(w, "{}")
}

// OpenResource is controller for opening a resource.
func OpenResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	rsrc := pm.GetResource(model.ResourceID(id))
	if rsrc != nil {
		pm.OpenResource(rsrc)
	} else {
		WriteJSONError(w, 404, "resource not found")
	}

	fmt.Fprint(w, "{}")
}

// ResourceResponse is struct used for JSON response.
type ResourcesResponse struct {
	Resources []*model.Resource `json:"resources"`
}

func writeResource(w http.ResponseWriter, pm *model.ProxyMap, id model.ResourceID) {
	rsrc := pm.GetResource(id)
	if rsrc != nil {
		WriteJSON(w, rsrc)
	} else {
		WriteJSONError(w, 404, "resource not found")
	}
}

func writeResources(w http.ResponseWriter, pm *model.ProxyMap, ids []model.ResourceID) {
	resp := ResourcesResponse{}
	for _, id := range ids {
		rsrc := pm.GetResource(id)
		if rsrc != nil {
			resp.Resources = append(resp.Resources, rsrc)
		} else {
			WriteJSONError(w, 404, "resource not found")
		}
	}
	WriteJSON(w, resp)
}
