// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package httpd

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"

	"github.com/filemaps/filemaps/pkg/config"
	"github.com/filemaps/filemaps/pkg/fileapp"
)

func routeConfig(r *httprouter.Router) {
	url := APIURL + "/config"
	r.GET(url, ReadConfig)
	r.PUT(url, WriteConfig)
}

func ReadConfig(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := make(map[string]interface{})
	resp["config"] = config.GetConfiguration()
	resp["fileApps"] = fileapp.GetInfos()
	WriteJSON(w, resp)
}

// ConfigCtrl is controller for configuration.
func WriteConfig(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cfg := config.GetConfiguration()
	d := json.NewDecoder(r.Body)
	err := d.Decode(cfg)
	r.Body.Close()
	if err != nil {
		WriteJSONError(w, 400, "bad request")
		return
	}

	err = cfg.Write()
	if err != nil {
		WriteJSONError(w, 500, "could not save config")
		return
	}

	fmt.Println("{}")
}
